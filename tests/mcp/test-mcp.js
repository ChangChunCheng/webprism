#!/usr/bin/env node
/**
 * WEBPRISM MCP Server 基礎連線測試
 * 測試項目：
 * 1. MCP server 啟動
 * 2. 列出所有工具
 * 3. 驗證 6 個工具是否正確註冊
 */

const { spawn } = require('child_process');
const readline = require('readline');

console.log('=================================');
console.log('WEBPRISM MCP 基礎連線測試');
console.log('=================================\n');

// 啟動 MCP server（載入 .env 環境變數）
console.log('📡 正在啟動 MCP server...');

// 讀取 .env 檔案（從專案根目錄）
const fs = require('fs');
const path = require('path');
const projectRoot = path.join(__dirname, '../..');
const mcpBinary = path.join(projectRoot, 'bin/webprism-mcp');
const envFile = fs.readFileSync(path.join(projectRoot, '.env'), 'utf8');
const envVars = {};

envFile.split('\n').forEach(line => {
  line = line.trim();
  if (line && !line.startsWith('#')) {
    const [key, ...valueParts] = line.split('=');
    if (key && valueParts.length > 0) {
      envVars[key] = valueParts.join('=');
    }
  }
});

console.log('✓ 已載入環境變數');
console.log('  - WEBPRISM_DATABASE_PASSWORD:', envVars.WEBPRISM_DATABASE_PASSWORD ? '***' : '未設定');
console.log('  - WEBPRISM_SECURITY_ENCRYPTION_KEY:', envVars.WEBPRISM_SECURITY_ENCRYPTION_KEY ? '***' : '未設定');

const mcpServer = spawn(mcpBinary, [], {
  stdio: ['pipe', 'pipe', 'pipe'],
  env: { ...process.env, ...envVars },
  cwd: projectRoot
});

let responseBuffer = '';
let testsPassed = 0;
let testsFailed = 0;

// 處理 stderr（日誌輸出）
mcpServer.stderr.on('data', (data) => {
  const log = data.toString();
  // 只顯示重要的日誌
  if (log.includes('MCP server')) {
    console.log('📋', log.trim());
  }
});

// 處理 stdout（JSON-RPC 回應）
const rl = readline.createInterface({
  input: mcpServer.stdout,
  crlfDelay: Infinity
});

let requestId = 0;

// 發送 JSON-RPC 請求的輔助函數
function sendRequest(method, params = {}) {
  requestId++;
  const request = {
    jsonrpc: '2.0',
    id: requestId,
    method: method,
    params: params
  };

  console.log(`\n📤 發送請求: ${method}`);
  mcpServer.stdin.write(JSON.stringify(request) + '\n');
  return requestId;
}

// 監聽回應
rl.on('line', (line) => {
  try {
    const response = JSON.parse(line);

    if (response.error) {
      console.log('❌ 錯誤回應:', response.error);
      testsFailed++;
    } else if (response.result !== undefined) {
      handleResponse(response);
    }
  } catch (e) {
    // 忽略非 JSON 輸出
  }
});

function handleResponse(response) {
  const method = response.id === 1 ? 'initialize' :
                 response.id === 2 ? 'tools/list' : 'unknown';

  console.log(`✅ 收到回應: ${method}`);

  if (response.id === 1) {
    // initialize 回應
    if (response.result && response.result.capabilities) {
      console.log('✓ MCP server 初始化成功');
      console.log('  - Server name:', response.result.serverInfo?.name || 'N/A');
      console.log('  - Server version:', response.result.serverInfo?.version || 'N/A');
      testsPassed++;

      // 發送 tools/list 請求
      setTimeout(() => sendRequest('tools/list'), 100);
    }
  } else if (response.id === 2) {
    // tools/list 回應
    if (response.result && response.result.tools) {
      const tools = response.result.tools;
      console.log(`\n✓ 成功列出 ${tools.length} 個工具:`);

      const expectedTools = [
        'upload_api_spec',
        'list_api_specs',
        'get_api_spec',
        'set_auth',
        'call_api',
        'health_check'
      ];

      tools.forEach((tool, idx) => {
        const hasExpected = expectedTools.includes(tool.name);
        const icon = hasExpected ? '✓' : '⚠';
        console.log(`  ${icon} ${idx + 1}. ${tool.name}`);
        console.log(`     描述: ${tool.description}`);
        if (hasExpected) testsPassed++;
      });

      // 檢查是否所有預期工具都存在
      const toolNames = tools.map(t => t.name);
      const missingTools = expectedTools.filter(t => !toolNames.includes(t));

      if (missingTools.length > 0) {
        console.log('\n❌ 缺少以下工具:', missingTools.join(', '));
        testsFailed += missingTools.length;
      }

      if (tools.length === 6 && missingTools.length === 0) {
        console.log('\n✅ 所有 6 個工具都正確註冊！');
      } else {
        console.log(`\n⚠️  預期 6 個工具，實際 ${tools.length} 個`);
      }

      // 測試完成
      setTimeout(() => {
        printSummary();
        process.exit(testsFailed > 0 ? 1 : 0);
      }, 500);
    }
  }
}

function printSummary() {
  console.log('\n=================================');
  console.log('測試摘要');
  console.log('=================================');
  console.log(`✅ 通過: ${testsPassed}`);
  console.log(`❌ 失敗: ${testsFailed}`);
  console.log(`📊 總計: ${testsPassed + testsFailed}`);

  if (testsFailed === 0) {
    console.log('\n🎉 所有測試通過！MCP 介面可用！');
  } else {
    console.log('\n⚠️  部分測試失敗，請檢查');
  }
  console.log('=================================\n');
}

// 錯誤處理
mcpServer.on('error', (error) => {
  console.error('❌ 啟動 MCP server 失敗:', error.message);
  process.exit(1);
});

mcpServer.on('exit', (code) => {
  if (code !== 0 && code !== null) {
    console.error(`❌ MCP server 異常退出，代碼: ${code}`);
  }
});

// 超時處理
const timeout = setTimeout(() => {
  console.error('❌ 測試超時（30秒）');
  mcpServer.kill();
  process.exit(1);
}, 30000);

// 開始測試：發送 initialize 請求
setTimeout(() => {
  sendRequest('initialize', {
    protocolVersion: '2024-11-05',
    capabilities: {},
    clientInfo: {
      name: 'webprism-test',
      version: '1.0.0'
    }
  });
}, 1000);

// 處理 Ctrl+C
process.on('SIGINT', () => {
  console.log('\n\n⚠️  測試被中斷');
  mcpServer.kill();
  clearTimeout(timeout);
  process.exit(130);
});

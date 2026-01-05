#!/usr/bin/env node
/**
 * WEBPRISM MCP Web 測試伺服器
 * 提供互動式 Web UI 來測試 MCP 介面
 */

const http = require('http');
const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');

const PORT = 3000;
const projectRoot = path.join(__dirname, '../..');

// MCP server 實例
let mcpServer = null;
let mcpReady = false;

// 載入環境變數
function loadEnv() {
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

  return envVars;
}

// 啟動 MCP server
function startMCPServer() {
  const envVars = loadEnv();
  const mcpBinary = path.join(projectRoot, 'bin/webprism-mcp');

  console.log('🚀 啟動 MCP server...');

  mcpServer = spawn(mcpBinary, [], {
    stdio: ['pipe', 'pipe', 'pipe'],
    env: { ...process.env, ...envVars },
    cwd: projectRoot
  });

  mcpServer.stderr.on('data', (data) => {
    const log = data.toString();
    if (log.includes('MCP server starting')) {
      mcpReady = true;
      console.log('✅ MCP server 已就緒');
    }
  });

  mcpServer.on('error', (error) => {
    console.error('❌ MCP server 錯誤:', error);
  });

  mcpServer.on('exit', (code) => {
    console.log('⚠️  MCP server 已退出:', code);
    mcpReady = false;
  });

  // 初始化
  setTimeout(() => {
    sendMCPRequest({
      jsonrpc: '2.0',
      id: 0,
      method: 'initialize',
      params: {
        protocolVersion: '2024-11-05',
        capabilities: {},
        clientInfo: { name: 'webprism-web-test', version: '1.0.0' }
      }
    });
  }, 1000);
}

// 發送 MCP 請求
function sendMCPRequest(request) {
  return new Promise((resolve, reject) => {
    if (!mcpServer) {
      reject(new Error('MCP server 未啟動'));
      return;
    }

    const requestStr = JSON.stringify(request) + '\n';
    mcpServer.stdin.write(requestStr);

    // 監聽回應
    const timeout = setTimeout(() => {
      reject(new Error('請求超時'));
    }, 10000);

    const handler = (data) => {
      const lines = data.toString().split('\n').filter(l => l.trim());
      for (const line of lines) {
        try {
          const response = JSON.parse(line);
          if (response.id === request.id) {
            clearTimeout(timeout);
            mcpServer.stdout.removeListener('data', handler);
            resolve(response);
            return;
          }
        } catch (e) {
          // 忽略非 JSON 輸出
        }
      }
    };

    mcpServer.stdout.on('data', handler);
  });
}

// HTTP 伺服器
const server = http.createServer(async (req, res) => {
  // CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (req.method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }

  // 路由
  if (req.url === '/' || req.url === '/index.html') {
    // 提供 HTML 頁面
    res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
    res.end(getHTMLPage());
  } else if (req.url === '/api/status') {
    // 狀態檢查
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      ready: mcpReady,
      serverRunning: mcpServer !== null
    }));
  } else if (req.url === '/api/tools/list' && req.method === 'GET') {
    // 列出工具
    try {
      const response = await sendMCPRequest({
        jsonrpc: '2.0',
        id: Date.now(),
        method: 'tools/list',
        params: {}
      });
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(response));
    } catch (error) {
      res.writeHead(500, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: error.message }));
    }
  } else if (req.url === '/api/tools/call' && req.method === 'POST') {
    // 呼叫工具
    let body = '';
    req.on('data', chunk => { body += chunk; });
    req.on('end', async () => {
      try {
        const { toolName, arguments: toolArgs } = JSON.parse(body);
        const response = await sendMCPRequest({
          jsonrpc: '2.0',
          id: Date.now(),
          method: 'tools/call',
          params: {
            name: toolName,
            arguments: toolArgs || {}
          }
        });
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(response));
      } catch (error) {
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: error.message }));
      }
    });
  } else {
    res.writeHead(404);
    res.end('Not Found');
  }
});

// HTML 頁面
function getHTMLPage() {
  return `<!DOCTYPE html>
<html lang="zh-TW">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>WEBPRISM MCP 互動測試</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      min-height: 100vh;
      padding: 20px;
    }
    .container {
      max-width: 1200px;
      margin: 0 auto;
    }
    .header {
      background: white;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 6px rgba(0,0,0,0.1);
      margin-bottom: 20px;
    }
    .header h1 {
      color: #667eea;
      margin-bottom: 10px;
      font-size: 32px;
    }
    .status {
      display: inline-block;
      padding: 5px 15px;
      border-radius: 20px;
      font-size: 14px;
      font-weight: 600;
    }
    .status.ready { background: #10b981; color: white; }
    .status.not-ready { background: #ef4444; color: white; }

    .grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 20px;
      margin-bottom: 20px;
    }

    .panel {
      background: white;
      padding: 25px;
      border-radius: 10px;
      box-shadow: 0 4px 6px rgba(0,0,0,0.1);
    }
    .panel h2 {
      color: #333;
      margin-bottom: 20px;
      font-size: 20px;
      border-bottom: 2px solid #667eea;
      padding-bottom: 10px;
    }

    .tool-list {
      max-height: 400px;
      overflow-y: auto;
    }
    .tool-item {
      padding: 15px;
      margin-bottom: 10px;
      background: #f9fafb;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.2s;
      border-left: 4px solid transparent;
    }
    .tool-item:hover {
      background: #f3f4f6;
      border-left-color: #667eea;
      transform: translateX(5px);
    }
    .tool-item.selected {
      background: #eef2ff;
      border-left-color: #667eea;
    }
    .tool-name {
      font-weight: 600;
      color: #667eea;
      margin-bottom: 5px;
    }
    .tool-desc {
      font-size: 14px;
      color: #6b7280;
    }

    .form-group {
      margin-bottom: 20px;
    }
    .form-group label {
      display: block;
      margin-bottom: 8px;
      font-weight: 600;
      color: #374151;
    }
    .form-group textarea {
      width: 100%;
      padding: 12px;
      border: 2px solid #e5e7eb;
      border-radius: 8px;
      font-family: 'Monaco', 'Courier New', monospace;
      font-size: 14px;
      resize: vertical;
      min-height: 150px;
    }
    .form-group textarea:focus {
      outline: none;
      border-color: #667eea;
    }

    .btn {
      padding: 12px 30px;
      background: #667eea;
      color: white;
      border: none;
      border-radius: 8px;
      font-size: 16px;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s;
    }
    .btn:hover {
      background: #5568d3;
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
    }
    .btn:disabled {
      background: #9ca3af;
      cursor: not-allowed;
      transform: none;
    }

    .result {
      background: #1f2937;
      color: #f9fafb;
      padding: 20px;
      border-radius: 8px;
      font-family: 'Monaco', 'Courier New', monospace;
      font-size: 14px;
      max-height: 500px;
      overflow-y: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
    .result.error {
      background: #fef2f2;
      color: #991b1b;
      border: 2px solid #fecaca;
    }

    .examples {
      background: white;
      padding: 25px;
      border-radius: 10px;
      box-shadow: 0 4px 6px rgba(0,0,0,0.1);
    }
    .example-btn {
      display: inline-block;
      margin: 5px;
      padding: 8px 16px;
      background: #f3f4f6;
      color: #374151;
      border-radius: 6px;
      cursor: pointer;
      font-size: 14px;
      transition: all 0.2s;
    }
    .example-btn:hover {
      background: #667eea;
      color: white;
    }

    .loading {
      text-align: center;
      padding: 20px;
      color: #6b7280;
    }

    @media (max-width: 768px) {
      .grid { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🌐 WEBPRISM MCP 互動測試</h1>
      <p style="color: #6b7280; margin-top: 10px;">
        狀態: <span id="status" class="status not-ready">檢查中...</span>
      </p>
    </div>

    <div class="grid">
      <div class="panel">
        <h2>📋 可用工具</h2>
        <div id="toolList" class="tool-list">
          <div class="loading">載入中...</div>
        </div>
      </div>

      <div class="panel">
        <h2>🎯 測試工具</h2>
        <div class="form-group">
          <label>選擇的工具: <span id="selectedTool" style="color: #667eea;">未選擇</span></label>
        </div>
        <div class="form-group">
          <label>參數 (JSON 格式)</label>
          <textarea id="toolArgs" placeholder='例如: {"limit": 10} 或 {}'>&#123;&#125;</textarea>
        </div>
        <button class="btn" onclick="callTool()" id="callBtn" disabled>執行工具</button>
      </div>
    </div>

    <div class="panel">
      <h2>📊 執行結果</h2>
      <div id="result" class="result">等待執行...</div>
    </div>

    <div class="examples">
      <h2>💡 快速測試範例</h2>
      <div style="margin-top: 15px;">
        <span class="example-btn" onclick="quickTest('list_specs')">列出所有規格</span>
        <span class="example-btn" onclick="quickTest('upload_spec')">上傳測試規格</span>
        <span class="example-btn" onclick="quickTest('health_check')">健康檢查</span>
      </div>
    </div>
  </div>

  <script>
    let selectedToolName = null;
    let tools = [];

    // 檢查狀態
    async function checkStatus() {
      try {
        const res = await fetch('/api/status');
        const data = await res.json();
        const statusEl = document.getElementById('status');
        if (data.ready) {
          statusEl.textContent = '✅ MCP Server 就緒';
          statusEl.className = 'status ready';
          loadTools();
        } else {
          statusEl.textContent = '⏳ MCP Server 啟動中...';
          statusEl.className = 'status not-ready';
          setTimeout(checkStatus, 2000);
        }
      } catch (error) {
        console.error('狀態檢查失敗:', error);
        setTimeout(checkStatus, 5000);
      }
    }

    // 載入工具列表
    async function loadTools() {
      try {
        const res = await fetch('/api/tools/list');
        const data = await res.json();

        if (data.result && data.result.tools) {
          tools = data.result.tools;
          renderTools(tools);
        }
      } catch (error) {
        document.getElementById('toolList').innerHTML =
          '<div class="loading" style="color: #ef4444;">載入失敗: ' + error.message + '</div>';
      }
    }

    // 渲染工具列表
    function renderTools(toolList) {
      const html = toolList.map(tool => \`
        <div class="tool-item" onclick="selectTool('\${tool.name}')">
          <div class="tool-name">\${tool.name}</div>
          <div class="tool-desc">\${tool.description}</div>
        </div>
      \`).join('');
      document.getElementById('toolList').innerHTML = html;
    }

    // 選擇工具
    function selectTool(name) {
      selectedToolName = name;
      document.getElementById('selectedTool').textContent = name;
      document.getElementById('callBtn').disabled = false;

      // 更新選中狀態
      document.querySelectorAll('.tool-item').forEach(el => {
        el.classList.remove('selected');
      });
      event.target.closest('.tool-item').classList.add('selected');

      // 根據工具類型設定預設參數
      const examples = {
        'list_api_specs': '{"limit": 10}',
        'upload_api_spec': '{"name": "test-api", "version": "1.0.0", "spec_data": {}}',
        'get_api_spec': '{"spec_id": "your-spec-id"}',
        'set_auth': '{"spec_id": "your-spec-id", "auth_type": "api_key", "credentials": {"key": "your-key"}}',
        'call_api': '{"spec_id": "your-spec-id", "operation_id": "getOperation"}',
        'health_check': '{"spec_id": "your-spec-id"}'
      };

      if (examples[name]) {
        document.getElementById('toolArgs').value = examples[name];
      }
    }

    // 呼叫工具
    async function callTool() {
      if (!selectedToolName) return;

      const argsText = document.getElementById('toolArgs').value;
      const resultEl = document.getElementById('result');
      const btn = document.getElementById('callBtn');

      try {
        const args = JSON.parse(argsText);
        btn.disabled = true;
        btn.textContent = '執行中...';
        resultEl.textContent = '⏳ 執行中...';
        resultEl.className = 'result';

        const res = await fetch('/api/tools/call', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            toolName: selectedToolName,
            arguments: args
          })
        });

        const data = await res.json();
        resultEl.textContent = JSON.stringify(data, null, 2);
        resultEl.className = data.error ? 'result error' : 'result';
      } catch (error) {
        resultEl.textContent = '❌ 錯誤: ' + error.message;
        resultEl.className = 'result error';
      } finally {
        btn.disabled = false;
        btn.textContent = '執行工具';
      }
    }

    // 快速測試
    function quickTest(type) {
      if (type === 'list_specs') {
        selectToolByName('list_api_specs');
      } else if (type === 'upload_spec') {
        selectToolByName('upload_api_spec');
      } else if (type === 'health_check') {
        // 需要先有 spec_id
        alert('請先上傳一個 API 規格，然後使用 get_api_spec 取得 spec_id');
      }
    }

    function selectToolByName(name) {
      const toolItem = Array.from(document.querySelectorAll('.tool-item'))
        .find(el => el.querySelector('.tool-name').textContent === name);
      if (toolItem) toolItem.click();
    }

    // 初始化
    checkStatus();
  </script>
</body>
</html>`;
}

// 啟動伺服器
server.listen(PORT, () => {
  console.log('╔════════════════════════════════════════════╗');
  console.log('║   WEBPRISM MCP 互動測試伺服器              ║');
  console.log('╚════════════════════════════════════════════╝');
  console.log('');
  console.log(`🌐 Web UI: http://localhost:${PORT}`);
  console.log('');
  console.log('📋 功能:');
  console.log('  • 視覺化工具列表');
  console.log('  • 互動式工具測試');
  console.log('  • 即時結果顯示');
  console.log('  • JSON 格式化');
  console.log('');
  console.log('⌨️  按 Ctrl+C 停止伺服器');
  console.log('');

  startMCPServer();
});

// 清理
process.on('SIGINT', () => {
  console.log('\\n\\n🛑 正在關閉伺服器...');
  if (mcpServer) {
    mcpServer.kill();
  }
  server.close(() => {
    console.log('✅ 伺服器已關閉');
    process.exit(0);
  });
});

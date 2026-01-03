package version

import (
	"fmt"
	"runtime"
	"time"
)

// 這些變數會在編譯時透過 -ldflags 注入
var (
	// Version 是應用程式版本號（例如：v1.0.0）
	Version = "dev"

	// GitCommit 是 Git commit hash
	GitCommit = "unknown"

	// GitBranch 是 Git 分支名稱
	GitBranch = "unknown"

	// BuildTime 是構建時間（RFC3339 格式）
	BuildTime = "unknown"

	// GoVersion 是 Go 編譯器版本
	GoVersion = runtime.Version()

	// Platform 是目標平台（OS/Arch）
	Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// Info 包含所有版本資訊
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	GitBranch string `json:"git_branch"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get 返回當前版本資訊
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		GitBranch: GitBranch,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
		Platform:  Platform,
	}
}

// String 返回格式化的版本字串
func (i Info) String() string {
	return fmt.Sprintf(
		"WEBPRISM %s\nGit: %s (%s)\nBuilt: %s\nGo: %s\nPlatform: %s",
		i.Version,
		i.GitCommit,
		i.GitBranch,
		i.BuildTime,
		i.GoVersion,
		i.Platform,
	)
}

// Short 返回簡短的版本字串
func (i Info) Short() string {
	commitHash := i.GitCommit
	if len(commitHash) > 7 {
		commitHash = commitHash[:7]
	}

	fullVer := i.FullVersion()
	return fmt.Sprintf("WEBPRISM %s (commit: %s)", fullVer, commitHash)
}

// FullVersion 返回完整的版本字串（根據狀況決定格式）
// 狀況一：開發中 -> branch+commit
// 狀況二：測試版本 -> alpha/beta+tag+commit
// 狀況三：正式版 -> release+tag+commit
func (i Info) FullVersion() string {
	commitHash := i.GitCommit
	if len(commitHash) > 7 {
		commitHash = commitHash[:7]
	}

	// 檢查是否有 tag（不只是 commit hash）
	hasTag := i.Version != "dev" && i.Version != "unknown" &&
		len(i.Version) > 7 && !isCommitHash(i.Version)

	if hasTag {
		// 有 tag 的情況
		cleanVer := i.Version

		// 移除 -dirty 後綴來判斷版本類型
		if len(cleanVer) > 6 && cleanVer[len(cleanVer)-6:] == "-dirty" {
			cleanVer = cleanVer[:len(cleanVer)-6]
		}

		// 判斷是測試版本還是正式版本
		if containsAlphaBeta(cleanVer) {
			// 測試版本：alpha/beta + tag + commit
			return fmt.Sprintf("alpha/beta+%s+%s", i.Version, commitHash)
		}

		// 正式版本：release + tag + commit
		return fmt.Sprintf("release+%s+%s", i.Version, commitHash)
	}

	// 開發中：branch + commit
	if i.GitBranch != "unknown" && i.GitBranch != "" {
		return fmt.Sprintf("%s+%s", i.GitBranch, commitHash)
	}

	// 降級：只有 commit hash
	return commitHash
}

// isCommitHash 檢查字串是否看起來像 commit hash
func isCommitHash(s string) bool {
	if len(s) < 7 {
		return false
	}
	// 檢查前7個字符是否都是十六進制
	for i := 0; i < 7 && i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// containsAlphaBeta 檢查版本號是否包含 alpha、beta 或 rc
func containsAlphaBeta(version string) bool {
	// 檢查是否包含 -alpha、-beta 或 -rc
	return len(version) > 0 && (
		contains(version, "-alpha") ||
		contains(version, "-beta") ||
		contains(version, "-rc"))
}

// contains 簡單的字串包含檢查
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOfSubstring(s, substr) >= 0
}

// indexOfSubstring 返回子字串在字串中的位置
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// IsDevelopment 判斷是否為開發版本
func (i Info) IsDevelopment() bool {
	// 有 tag 且不包含 -dirty 才是非開發版本
	hasCleanTag := i.Version != "dev" && i.Version != "unknown" &&
		!isCommitHash(i.Version) &&
		!(len(i.Version) > 6 && i.Version[len(i.Version)-6:] == "-dirty")

	return !hasCleanTag
}

// GetReleaseType 返回版本類型
func (i Info) GetReleaseType() string {
	if i.IsDevelopment() {
		return "DEVELOPMENT"
	}

	cleanVer := i.Version
	if len(cleanVer) > 6 && cleanVer[len(cleanVer)-6:] == "-dirty" {
		cleanVer = cleanVer[:len(cleanVer)-6]
	}

	if containsAlphaBeta(cleanVer) {
		return "ALPHA/BETA"
	}

	return "RELEASE"
}

// ParseBuildTime 解析構建時間為 time.Time
func (i Info) ParseBuildTime() (time.Time, error) {
	if i.BuildTime == "unknown" {
		return time.Time{}, fmt.Errorf("build time is unknown")
	}
	return time.Parse(time.RFC3339, i.BuildTime)
}

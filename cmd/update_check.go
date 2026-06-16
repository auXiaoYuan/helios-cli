package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	npmLatestURL = "https://registry.npmjs.org/helios-cli/latest"
	cacheTTL     = 24 * time.Hour
	httpTimeout  = 1500 * time.Millisecond
	reportWait   = 250 * time.Millisecond
)

type updateResult struct {
	latest string
	err    error
}

// startUpdateCheck 异步去 npm registry 查最新版本号。
// 返回一个 channel，调用方在命令结束时尝试读取结果。
// 设计原则：
//   - 任何网络错误都静默忽略，绝不影响 CLI 主命令
//   - 24h 缓存避免每次启动都打 HTTP
//   - 开发态（version="dev"）跳过检查
//   - 用户可通过 HELIOS_CLI_NO_UPDATE_CHECK=1 关闭
func startUpdateCheck(currentVersion string) <-chan updateResult {
	ch := make(chan updateResult, 1)

	if currentVersion == "" || currentVersion == "dev" {
		close(ch)
		return ch
	}
	if os.Getenv("HELIOS_CLI_NO_UPDATE_CHECK") == "1" {
		close(ch)
		return ch
	}

	go func() {
		defer close(ch)

		if cached, ok := readCachedLatest(); ok {
			ch <- updateResult{latest: cached}
			return
		}

		latest, err := fetchLatestVersion()
		if err != nil {
			ch <- updateResult{err: err}
			return
		}
		_ = writeCachedLatest(latest)
		ch <- updateResult{latest: latest}
	}()

	return ch
}

// reportUpdate 在主命令执行完后短暂等待 update 检查结果，
// 拿到且版本不一致就在 stderr 打一行温和提示；其他情况静默。
func reportUpdate(ch <-chan updateResult, currentVersion string) {
	if ch == nil {
		return
	}
	select {
	case res, ok := <-ch:
		if !ok || res.err != nil || res.latest == "" {
			return
		}
		if res.latest == currentVersion {
			return
		}
		fmt.Fprintf(os.Stderr,
			"\n[helios-cli] new version available: %s (current %s)\n"+
				"            run: npm install -g helios-cli@latest\n",
			res.latest, currentVersion,
		)
	case <-time.After(reportWait):
		// 在限定时间内没拿到结果，放弃，避免拖慢命令
	}
}

func fetchLatestVersion() (string, error) {
	client := &http.Client{Timeout: httpTimeout}
	req, err := http.NewRequest(http.MethodGet, npmLatestURL, nil)
	if err != nil {
		return "", err
	}
	// npm 推荐的轻量 metadata content-type，响应更小
	req.Header.Set("Accept", "application/vnd.npm.install-v1+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registry returned %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if err != nil {
		return "", err
	}
	var data struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	return strings.TrimSpace(data.Version), nil
}

func cacheFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(home, ".cache", "helios-cli")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ""
	}
	return filepath.Join(dir, "update-check")
}

func readCachedLatest() (string, bool) {
	path := cacheFilePath()
	if path == "" {
		return "", false
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", false
	}
	if time.Since(info.ModTime()) > cacheTTL {
		return "", false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(data)), true
}

func writeCachedLatest(v string) error {
	path := cacheFilePath()
	if path == "" {
		return nil
	}
	return os.WriteFile(path, []byte(v), 0o644)
}

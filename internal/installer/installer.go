package installer

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"WinDevReady/internal/config"
	"WinDevReady/internal/logger"
	"WinDevReady/internal/store"
)

// Installer 安装引擎
type Installer struct {
	log   *logger.Logger
	store *store.Records
}

// New 创建安装引擎
func New(log *logger.Logger, store *store.Records) *Installer {
	return &Installer{log: log, store: store}
}

// InstallTools 批量安装工具，返回成功/失败数量
func (i *Installer) InstallTools(tools []config.ToolDef) (int, int) {
	var (
		mu      sync.Mutex
		success int
		failed  int
	)

	for _, tool := range tools {
		// 检查是否已安装
		if installed, _ := i.isInstalled(tool); installed {
			i.log.Info(tool.ID, fmt.Sprintf("%s 已安装，跳过", tool.Name))
			mu.Lock()
			success++
			mu.Unlock()
			continue
		}

		// 执行安装
		if err := i.installOne(tool); err != nil {
			i.log.Error(tool.ID, fmt.Sprintf("安装失败: %s", err))
			mu.Lock()
			failed++
			mu.Unlock()
			continue
		}

		// 验证安装结果
		if version, ok := i.verifyInstall(tool); ok {
			// 写入安装记录
			i.store.Upsert(store.Record{
				ToolID:        tool.ID,
				Name:          tool.Name,
				Version:       version,
				InstallMethod: string(tool.Method),
			})
			_ = i.store.Save()
			i.log.Success(tool.ID, fmt.Sprintf("安装成功 [%s] %s", tool.Name, version))
			mu.Lock()
			success++
			mu.Unlock()
		} else {
			i.log.Warn(tool.ID, fmt.Sprintf("%s 已安装但无法获取版本号", tool.Name))
			mu.Lock()
			success++
			mu.Unlock()
		}
	}

	return success, failed
}

// installOne 执行单个工具安装
func (i *Installer) installOne(tool config.ToolDef) error {
	i.log.Info(tool.ID, fmt.Sprintf("正在安装 %s ...", tool.Name))

	switch tool.Method {
	case config.MethodWinget:
		return i.installWinget(tool)
	case config.MethodNpm:
		return i.installNpm(tool)
	case config.MethodPip:
		return i.installPip(tool)
	case config.MethodChoco:
		return i.installChoco(tool)
	case config.MethodDownload:
		return i.installDownload(tool)
	default:
		return fmt.Errorf("不支持的安装方式: %s", tool.Method)
	}
}

// installWinget 通过 winget 安装
func (i *Installer) installWinget(tool config.ToolDef) error {
	args := []string{"install", "--id", tool.WingetID, "--accept-package-agreements", "--accept-source-agreements"}
	cmd := exec.Command("winget", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

// installNpm 通过 npm 全局安装
func (i *Installer) installNpm(tool config.ToolDef) error {
	args := []string{"install", "-g", tool.Package}
	cmd := exec.Command("npm", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

// installPip 通过 pip 安装
func (i *Installer) installPip(tool config.ToolDef) error {
	args := []string{"install", tool.Package}
	cmd := exec.Command("pip", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

// installChoco 通过 chocolatey 安装
func (i *Installer) installChoco(tool config.ToolDef) error {
	args := []string{"install", tool.Package, "-y"}
	cmd := exec.Command("choco", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

// installDownload 直接下载安装包并静默安装
func (i *Installer) installDownload(tool config.ToolDef) error {
	if len(tool.DownloadURLs) == 0 {
		return fmt.Errorf("未配置下载地址")
	}

	// 创建下载器
	dl, err := NewDownloader()
	if err != nil {
		return err
	}
	defer dl.Cleanup()

	// 从 URL 中提取文件名
	filename := extractFilename(tool.DownloadURLs[0])
	i.log.Info(tool.ID, fmt.Sprintf("正在下载: %s", filename))

	// 带回退的下载
	savePath, err := dl.DownloadWithFallback(tool.DownloadURLs, filename, func(downloaded, total int64) {
		if total > 0 {
			percent := float64(downloaded) / float64(total) * 100
			i.log.Info(tool.ID, fmt.Sprintf("下载进度: %.1f%% (%d/%d MB)", percent, downloaded/1024/1024, total/1024/1024))
		}
	})
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	i.log.Info(tool.ID, "下载完成，开始静默安装...")

	// 执行静默安装
	return i.runSilentInstall(tool, savePath)
}

// runSilentInstall 执行静默安装
func (i *Installer) runSilentInstall(tool config.ToolDef, installerPath string) error {
	args := []string{installerPath}
	// 添加静默安装参数
	if tool.SilentArgs != "" {
		for _, arg := range parseArgs(tool.SilentArgs) {
			args = append(args, arg)
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	return i.runCmdWithLog(tool.ID, cmd)
}

// extractFilename 从 URL 中提取文件名
func extractFilename(url string) string {
	// 找到最后一个 / 的位置
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			name := url[i+1:]
			if name != "" {
				return name
			}
		}
	}
	return "installer.exe"
}

// parseArgs 解析带空格的参数字符串
func parseArgs(s string) []string {
	var args []string
	current := ""
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch == '"' || ch == '\'':
			if inQuote && ch == quoteChar {
				inQuote = false
			} else if !inQuote {
				inQuote = true
				quoteChar = ch
			} else {
				current += string(ch)
			}
		case ch == ' ' && !inQuote:
			if current != "" {
				args = append(args, current)
				current = ""
			}
		default:
			current += string(ch)
		}
	}
	if current != "" {
		args = append(args, current)
	}
	return args
}

// isInstalled 检查工具是否已安装
func (i *Installer) isInstalled(tool config.ToolDef) (bool, string) {
	if tool.VerifyCmd == "" {
		return false, ""
	}
	parts := strings.Fields(tool.VerifyCmd)
	if len(parts) == 0 {
		return false, ""
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, ""
	}
	return true, strings.TrimSpace(string(out))
}

// verifyInstall 安装后验证，返回版本号
func (i *Installer) verifyInstall(tool config.ToolDef) (string, bool) {
	_, version := i.isInstalled(tool)
	if version == "" {
		return "", false
	}
	return version, true
}

// runCmdWithLog 运行命令并实时输出日志
func (i *Installer) runCmdWithLog(toolID string, cmd *exec.Cmd) error {
	// 使用 CombinedOutput 捕获输出
	// 注意：真实场景需要用 cmd.StdoutPipe + scanner 实现流式输出
	// 这里先用简化的同步方式
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 将错误输出也记录到日志
		if len(output) > 0 {
			for _, line := range strings.Split(string(output), "\n") {
				line = strings.TrimSpace(line)
				if line != "" {
					i.log.Info(toolID, line)
				}
			}
		}
		return fmt.Errorf("命令执行失败: %w, 输出: %s", err, string(output))
	}

	// 成功输出
	if len(output) > 0 {
		for _, line := range strings.Split(string(output), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				i.log.Info(toolID, line)
			}
		}
	}
	return nil
}

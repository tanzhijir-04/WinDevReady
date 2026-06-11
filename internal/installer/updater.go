package installer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"WinDevReady/internal/config"
	"WinDevReady/internal/store"
)

// UpdateInfo 单个工具的更新信息
type UpdateInfo struct {
	ToolID         string `json:"tool_id"`
	ToolName       string `json:"tool_name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	HasUpdate      bool   `json:"has_update"`
}

// CheckUpdates 检查所有已安装工具的更新
func (i *Installer) CheckUpdates() []UpdateInfo {
	records := i.store.GetAll()
	var updates []UpdateInfo

	for _, rec := range records {
		info := UpdateInfo{
			ToolID:         rec.ToolID,
			ToolName:       rec.Name,
			CurrentVersion: rec.Version,
		}

		// 根据安装方式查询最新版本
		tool, found := config.GetToolByID(rec.ToolID)
		if !found {
			continue
		}

		latest, err := i.fetchLatestVersion(tool)
		if err != nil {
			i.log.Warn(rec.ToolID, fmt.Sprintf("查询最新版本失败: %s", err))
			continue
		}

		info.LatestVersion = latest
		info.HasUpdate = info.CurrentVersion != latest && latest != ""
		updates = append(updates, info)
	}

	return updates
}

// UpdateTool 升级单个工具
func (i *Installer) UpdateTool(toolID string) error {
	tool, found := config.GetToolByID(toolID)
	if !found {
		return fmt.Errorf("工具定义不存在: %s", toolID)
	}

	i.log.Info(toolID, fmt.Sprintf("正在升级 %s ...", tool.Name))

	// winget 方式可以直接升级
	if tool.Method == config.MethodWinget {
		args := []string{"upgrade", "--id", tool.WingetID, "--accept-package-agreements", "--accept-source-agreements"}
		cmd := createCommand(args...)
		if err := i.runCmdWithLog(toolID, cmd); err != nil {
			return err
		}
	} else {
		// npm/pip 等方式：先卸载再安装
		if err := i.UninstallTool(toolID); err != nil {
			i.log.Warn(toolID, fmt.Sprintf("卸载旧版本警告: %s", err))
		}
		if err := i.installOne(tool); err != nil {
			return err
		}
	}

	// 更新安装记录
	if version, ok := i.verifyInstall(tool); ok {
		i.store.Upsert(store.Record{
			ToolID:        toolID,
			Name:          tool.Name,
			Version:       version,
			InstallMethod: string(tool.Method),
		})
		_ = i.store.Save()
		i.log.Success(toolID, fmt.Sprintf("升级完成 [%s] %s", tool.Name, version))
	}

	return nil
}

// fetchLatestVersion 查询工具最新版本
func (i *Installer) fetchLatestVersion(tool config.ToolDef) (string, error) {
	switch tool.Method {
	case config.MethodNpm:
		return i.fetchNpmLatest(tool.Package)
	case config.MethodWinget:
		return i.fetchWingetLatest(tool.WingetID)
	default:
		return "", fmt.Errorf("暂不支持查询 %s 方式的最新版本", tool.Method)
	}
}

// fetchNpmLatest 从 npm registry 查询最新版本
func (i *Installer) fetchNpmLatest(pkg string) (string, error) {
	url := fmt.Sprintf("https://registry.npmmirror.com/%s/latest", pkg)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		// 回退到官方 registry
		url = fmt.Sprintf("https://registry.npmjs.org/%s/latest", pkg)
		resp, err = client.Get(url)
		if err != nil {
			return "", err
		}
	}
	defer resp.Body.Close()

	var result struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Version, nil
}

// fetchWingetLatest 通过 winget 查询最新版本
func (i *Installer) fetchWingetLatest(wingetID string) (string, error) {
	// winget 没有直接查版本的命令，用 show 获取
	args := []string{"show", "--id", wingetID}
	cmd := createCommand(args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// 解析输出中的版本号
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "版本:") || strings.HasPrefix(line, "Version:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", nil
}

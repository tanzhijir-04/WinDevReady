package installer

import (
	"fmt"
	"os/exec"
	"strings"

	"WinDevReady/internal/config"
	"WinDevReady/internal/store" // 导入安装记录模块
)

// UninstallTool 卸载单个工具
func (i *Installer) UninstallTool(toolID string) error {
	// 只允许卸载安装记录中的工具
	record, ok := i.store.Get(toolID)
	if !ok {
		return fmt.Errorf("工具 %s 不在安装记录中，跳过卸载", toolID)
	}

	// 获取工具定义
	tool, found := config.GetToolByID(toolID)
	if !found {
		return fmt.Errorf("工具定义不存在: %s", toolID)
	}

	i.log.Info(toolID, fmt.Sprintf("正在卸载 %s ...", record.Name))

	var err error
	switch tool.Method {
	case config.MethodWinget:
		err = i.uninstallWinget(tool)
	case config.MethodNpm:
		err = i.uninstallNpm(tool)
	case config.MethodPip:
		err = i.uninstallPip(tool)
	case config.MethodChoco:
		err = i.uninstallChoco(tool)
	default:
		return fmt.Errorf("不支持的卸载方式: %s", tool.Method)
	}

	if err != nil {
		i.log.Error(toolID, fmt.Sprintf("卸载失败: %s", err))
		return err
	}

	// 清理安装记录
	i.store.Remove(toolID)
	_ = i.store.Save()
	i.log.Success(toolID, fmt.Sprintf("%s 卸载完成", record.Name))
	return nil
}

// UninstallAll 卸载所有已记录的工具
func (i *Installer) UninstallAll() (int, int) {
	records := i.store.GetAll()
	var success, failed int
	for _, rec := range records {
		if err := i.UninstallTool(rec.ToolID); err != nil {
			failed++
		} else {
			success++
		}
	}
	return success, failed
}

func (i *Installer) uninstallWinget(tool config.ToolDef) error {
	args := []string{"uninstall", "--id", tool.WingetID, "--accept-source-agreements"}
	cmd := exec.Command("winget", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

func (i *Installer) uninstallNpm(tool config.ToolDef) error {
	args := []string{"uninstall", "-g", tool.Package}
	cmd := exec.Command("npm", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

func (i *Installer) uninstallPip(tool config.ToolDef) error {
	args := []string{"uninstall", "-y", tool.Package}
	cmd := exec.Command("pip", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

func (i *Installer) uninstallChoco(tool config.ToolDef) error {
	args := []string{"uninstall", tool.Package, "-y"}
	cmd := exec.Command("choco", args...)
	return i.runCmdWithLog(tool.ID, cmd)
}

// CleanUninstallReport 卸载后输出清理报告
func (i *Installer) CleanUninstallReport() string {
	records := i.store.GetAll()
	var sb strings.Builder
	sb.WriteString("=== 卸载清理报告 ===\n")
	if len(records) == 0 {
		sb.WriteString("当前无已记录的安装工具\n")
	} else {
		sb.WriteString(fmt.Sprintf("已清理 %d 个工具的安装记录\n", len(records)))
		for _, rec := range records {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", rec.Name, rec.ToolID))
		}
	}
	sb.WriteString(fmt.Sprintf("\n安装记录文件: %s", i.store.FilePath()))
	return sb.String()
}

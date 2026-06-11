package verify

import (
	"fmt"
	"os/exec"
	"strings"

	"WinDevReady/internal/config"
	"WinDevReady/internal/logger"
)

// ToolStatus 单个工具的验证结果
type ToolStatus struct {
	ToolID    string `json:"tool_id"`
	ToolName  string `json:"tool_name"`
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
	Status    string `json:"status"` // "✅ 正常" / "❌ 未安装" / "⚠️ 异常"
}

// Report 完整环境报告
type Report struct {
	Items  []ToolStatus `json:"items"`
	Summary string       `json:"summary"`
}

// Verifier 环境验证器
type Verifier struct {
	log *logger.Logger
}

// New 创建验证器
func New(log *logger.Logger) *Verifier {
	return &Verifier{log: log}
}

// VerifyAll 逐项检测所有工具，生成报告卡
func (v *Verifier) VerifyAll() Report {
	report := Report{}
	passCount := 0
	total := len(config.Tools)

	for _, tool := range config.Tools {
		v.log.Info(tool.ID, fmt.Sprintf("正在验证 %s ...", tool.Name))
		status := v.verifyOne(tool)
		report.Items = append(report.Items, status)

		if status.Installed {
			passCount++
		}
	}

	report.Summary = fmt.Sprintf("环境检测完成：%d/%d 项通过", passCount, total)
	return report
}

// VerifyTools 只验证指定工具列表
func (v *Verifier) VerifyTools(tools []config.ToolDef) Report {
	report := Report{}
	passCount := 0

	for _, tool := range tools {
		status := v.verifyOne(tool)
		report.Items = append(report.Items, status)
		if status.Installed {
			passCount++
		}
	}

	report.Summary = fmt.Sprintf("环境检测完成：%d/%d 项通过", passCount, len(tools))
	return report
}

// verifyOne 验证单个工具
func (v *Verifier) verifyOne(tool config.ToolDef) ToolStatus {
	status := ToolStatus{
		ToolID:   tool.ID,
		ToolName: tool.Name,
	}

	if tool.VerifyCmd == "" {
		status.Status = "⚠️ 无验证命令"
		return status
	}

	parts := strings.Fields(tool.VerifyCmd)
	if len(parts) == 0 {
		status.Status = "⚠️ 无验证命令"
		return status
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		status.Installed = false
		status.Status = "❌ 未安装"
		return status
	}

	status.Installed = true
	status.Version = strings.TrimSpace(string(output))
	status.Status = "✅ 正常"
	return status
}

// FormatReport 将报告格式化为可显示的字符串
func FormatReport(report Report) string {
	var sb strings.Builder
	sb.WriteString("╔══════════════════════════════════════════════╗\n")
	sb.WriteString("║           环 境 验 证 报 告 卡              ║\n")
	sb.WriteString("╠══════════════════════════════════════════════╣\n")

	// 按分组输出
	currentGroup := ""
	for _, item := range report.Items {
		// 找到工具对应的分组
		tool, _ := config.GetToolByID(item.ToolID)
		if string(tool.Group) != currentGroup {
			currentGroup = string(tool.Group)
			sb.WriteString(fmt.Sprintf("\n【%s】\n", currentGroup))
		}
		sb.WriteString(fmt.Sprintf("  %s  %-20s  %s\n", item.Status, item.ToolName, item.Version))
	}

	sb.WriteString("╠══════════════════════════════════════════════╣\n")
	sb.WriteString(fmt.Sprintf("  %s\n", report.Summary))
	sb.WriteString("╚══════════════════════════════════════════════╝\n")
	return sb.String()
}

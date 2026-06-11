package installer

import (
	"os/exec"
)

// createCommand 跨平台创建命令（Windows 用 cmd /c 包装）
func createCommand(args ...string) *exec.Cmd {
	if len(args) == 0 {
		return exec.Command("cmd")
	}
	// 使用 cmd /c 执行，确保 Windows 兼容
	return exec.Command("cmd", append([]string{"/c"}, args...)...)
}

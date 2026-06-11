package config

// ============================================================
// 工具配置数据 —— 社区贡献新工具只需在此添加一条记录
// ============================================================

// InstallMethod 安装方式
type InstallMethod string

const (
	MethodWinget     InstallMethod = "winget"      // 通过 winget 安装
	MethodNpm        InstallMethod = "npm"         // 通过 npm install -g 安装
	MethodPip        InstallMethod = "pip"         // 通过 pip install 安装
	MethodDownload   InstallMethod = "download"    // 直接下载安装包执行
	MethodChoco      InstallMethod = "choco"       // 通过 chocolatey 安装
)

// ToolGroup 工具分组
type ToolGroup string

const (
	GroupRuntime   ToolGroup = "基础运行时"
	GroupAICLI     ToolGroup = "AI CLI 工具"
	GroupEditor    ToolGroup = "编辑器"
	GroupTerminal  ToolGroup = "终端增强"
)

// ToolDef 单个工具的配置定义
type ToolDef struct {
	ID            string        `json:"id"`            // 唯一标识，如 "nodejs"
	Name          string        `json:"name"`          // 显示名称，如 "Node.js LTS"
	Group         ToolGroup     `json:"group"`         // 所属分组
	Method        InstallMethod `json:"method"`        // 安装方式
	Package       string        `json:"package"`       // 包名 / npm 包名
	VerifyCmd     string        `json:"verify_cmd"`    // 验证命令，如 "node -v"
	MinVersion    string        `json:"min_version"`   // 最低版本（可选）
	Description   string        `json:"description"`   // 一行说明
	WingetID      string        `json:"winget_id"`     // winget 包 ID（winget 方式时必填）
	PostInstall   string        `json:"post_install"`  // 安装后执行的命令（可选）
	DownloadURLs  []string      `json:"download_urls"` // 下载地址列表（download 方式时必填，按优先级排列）
	SilentArgs    string        `json:"silent_args"`   // 静默安装参数（download 方式时必填）
}

// Tools 全部工具定义（数据驱动，按分组排列）
var Tools = []ToolDef{
	// ── 基础运行时 ──────────────────────────────────────
	{
		ID:          "nodejs",
		Name:        "Node.js LTS",
		Group:       GroupRuntime,
		Method:      MethodWinget,
		Package:     "OpenJS.NodeJS.LTS",
		VerifyCmd:   "node -v",
		Description: "JavaScript 运行时，AI CLI 工具依赖",
		WingetID:    "OpenJS.NodeJS.LTS",
	},
	{
		ID:       "git",
		Name:     "Git",
		Group:    GroupRuntime,
		Method:   MethodDownload,
		Package:  "Git.Git",
		VerifyCmd: "git --version",
		Description: "版本控制工具，Copilot/Claude Code 依赖",
		WingetID: "Git.Git",
		// 优先使用 GitHub 镜像，回退到官方和国内镜像
		DownloadURLs: []string{
			"https://ghfast.top/https://github.com/git-for-windows/git/releases/latest/download/Git-2.49.0-64-bit.exe",
			"https://mirror.ghproxy.com/https://github.com/git-for-windows/git/releases/latest/download/Git-2.49.0-64-bit.exe",
			"https://github.com/git-for-windows/git/releases/latest/download/Git-2.49.0-64-bit.exe",
			"https://registry.npmmirror.com/-/binary/git-for-windows/2.49.0.windows.1/Git-2.49.0-64-bit.exe",
		},
		SilentArgs: "/VERYSILENT /NORESTART /NOCANCEL /SP- /CLOSEAPPLICATIONS /RESTARTAPPLICATIONS /COMPONENTS=icons,ext,ext\\shell\\here,gitlfs,assoc,assoc_sh",
	},
	{
		ID:          "python",
		Name:        "Python 3",
		Group:       GroupRuntime,
		Method:      MethodWinget,
		Package:     "Python.Python.3.12",
		VerifyCmd:   "python --version",
		Description: "Python 运行时",
		WingetID:    "Python.Python.3.12",
	},

	// ── AI CLI 工具 ──────────────────────────────────────
	{
		ID:          "claude-code",
		Name:        "Claude Code",
		Group:       GroupAICLI,
		Method:      MethodNpm,
		Package:     "@anthropic-ai/claude-code",
		VerifyCmd:   "claude --version",
		Description: "Anthropic 官方 AI 编程 CLI",
	},
	{
		ID:          "codex-cli",
		Name:        "Codex CLI",
		Group:       GroupAICLI,
		Method:      MethodNpm,
		Package:     "@openai/codex",
		VerifyCmd:   "codex --version",
		Description: "OpenAI 官方 AI 编程 CLI",
	},
	{
		ID:          "echobird",
		Name:        "Echobird",
		Group:       GroupAICLI,
		Method:      MethodNpm,
		Package:     "echobird",
		VerifyCmd:   "echobird --version",
		Description: "AI 编程助手 CLI",
	},

	// ── 编辑器 ──────────────────────────────────────────
	{
		ID:          "vscode",
		Name:        "VS Code",
		Group:       GroupEditor,
		Method:      MethodWinget,
		Package:     "Microsoft.VisualStudioCode",
		VerifyCmd:   "code --version",
		Description: "轻量级代码编辑器",
		WingetID:    "Microsoft.VisualStudioCode",
	},
	{
		ID:          "cursor",
		Name:        "Cursor",
		Group:       GroupEditor,
		Method:      MethodWinget,
		Package:     "Cursor.Cursor",
		VerifyCmd:   "cursor --version",
		Description: "AI 原生代码编辑器",
		WingetID:    "Cursor.Cursor",
	},

	// ── 终端增强 ──────────────────────────────────────────
	{
		ID:          "windows-terminal",
		Name:        "Windows Terminal",
		Group:       GroupTerminal,
		Method:      MethodWinget,
		Package:     "Microsoft.WindowsTerminal",
		VerifyCmd:   "wt --version",
		Description: "现代终端模拟器",
		WingetID:    "Microsoft.WindowsTerminal",
	},
	{
		ID:          "oh-my-posh",
		Name:        "Oh My Posh",
		Group:       GroupTerminal,
		Method:      MethodWinget,
		Package:     "JanDeDobbeleer.OhMyPosh",
		VerifyCmd:   "oh-my-posh --version",
		Description: "终端提示符美化工具",
		WingetID:    "JanDeDobbeleer.OhMyPosh",
	},
}

// GetToolsByGroup 按分组获取工具列表
func GetToolsByGroup(group ToolGroup) []ToolDef {
	var result []ToolDef
	for _, t := range Tools {
		if t.Group == group {
			result = append(result, t)
		}
	}
	return result
}

// GetToolByID 根据 ID 获取工具定义
func GetToolByID(id string) (ToolDef, bool) {
	for _, t := range Tools {
		if t.ID == id {
			return t, true
		}
	}
	return ToolDef{}, false
}

// AllGroups 返回所有分组（有序）
var AllGroups = []ToolGroup{
	GroupRuntime,
	GroupAICLI,
	GroupEditor,
	GroupTerminal,
}

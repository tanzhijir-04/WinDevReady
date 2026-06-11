<p align="center">
  <img src="assets/Icon.png" width="120" alt="WinDevReady Icon">
</p>

<h1 align="center">WinDevReady</h1>

<p align="center">
  <strong>Windows 一键 AI 开发环境配置工具</strong><br>
  <sub>帮你的朋友快速配好 AI 开发环境，操作者是开发者本人，受益者是对方。</sub>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Windows-10%2F11-0078D4?logo=windows&logoColor=white" alt="Windows">
  <img src="https://img.shields.io/badge/License-MIT-brightgreen" alt="License">
  <img src="https://img.shields.io/badge/Fyne-v2-5FAA5F?logo=github&logoColor=white" alt="Fyne">
</p>

---

## 功能概览

| 模块 | 说明 |
|------|------|
| **网络检测** | 自动检测 npm/GitHub 连通性，不通则切换国内镜像；检测代理并自动配置 |
| **一键安装** | 勾选工具，点击安装，已装自动跳过，日志实时输出 |
| **版本更新** | 对比最新版本，支持逐个或批量升级 |
| **卸载清理** | 只清理本工具安装的记录，不触碰用户原有环境 |
| **环境验证** | 逐项检测工具可用性，输出环境报告卡 |

## 可安装工具

| 分组 | 工具 |
|------|------|
| 基础运行时 | Node.js LTS · Git · Python 3 |
| AI CLI 工具 | Claude Code · Codex CLI · Echobird |
| 编辑器 | VS Code · Cursor |
| 终端增强 | Windows Terminal · Oh My Posh |

> 社区贡献新工具只需在 `internal/config/tools.go` 添加一条记录。

## 快速开始

### 下载

从 [GitHub Actions](https://github.com/tanzhijir-04/WinDevReady/actions) 下载最新编译产物，解压后双击运行。

### 编译

```bash
git clone https://github.com/tanzhijir-04/WinDevReady.git
cd WinDevReady
go mod tidy
go build -o WinDevReady.exe -ldflags="-s -w" .
```

## 添加新工具

在 `internal/config/tools.go` 中添加：

```go
{
    ID:          "my-tool",
    Name:        "My Tool",
    Group:       GroupAICLI,
    Method:      MethodNpm,
    Package:     "my-tool-pkg",
    VerifyCmd:   "my-tool --version",
    Description: "工具说明",
},
```

支持的安装方式：`winget` · `npm` · `pip` · `download` · `choco`

## 项目结构

```
WinDevReady/
├── main.go                    # 入口
├── assets/                    # 图标资源
├── internal/
│   ├── config/tools.go        # 工具配置（数据驱动）
│   ├── installer/             # 安装/卸载/升级引擎
│   ├── network/detector.go    # 网络检测
│   ├── logger/logger.go       # 流式日志
│   ├── store/records.go       # 安装记录
│   ├── verify/verifier.go     # 环境验证
│   └── ui/                    # 界面
│       ├── app.go             # 主窗口
│       ├── theme.go           # 自定义深色主题
│       └── ..._page.go        # 各功能页面
└── .github/workflows/         # CI 自动编译
```

## 技术栈

- **语言**：Go 1.23
- **GUI**：[Fyne v2](https://fyne.io/) · 自定义深色主题
- **产物**：单文件 `.exe`，无运行时依赖
- **CI**：GitHub Actions · Windows Runner

## 赞助

如果这个工具帮到了你，欢迎请我喝杯咖啡 ☕

👉 [爱发电赞助](https://ifdian.net/a/tanz666)

## License

[MIT](LICENSE)

# WinDevReady

> 帮朋友一键配好 Windows 上的 AI 开发环境。

你有没有过这种经历：朋友说「我想试试更高级的 AI 工具」，然后你花了两小时帮他装 Node.js、Git、Python、Claude Code……装完还各种报错。

**WinDevReady 就是干这个的**——打开工具，勾选要装的东西，点一下，剩下的全自动。

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/Windows-10%2F11-0078D4?logo=windows&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-brightgreen)

---

## 它能干什么

### 1. 自动搞定网络问题

检测你的网络能不能访问 npm 和 GitHub。如果不通，自动切换到国内镜像（npmmirror、清华源）。如果开了代理，自动给 npm、git、pip 配好代理地址。

**你不需要手动配置任何东西。**

### 2. 一键安装开发工具

打开工具，你会看到四个分组，勾选要装的，点「一键安装」：

| 分组 | 有哪些工具 | 怎么装 |
|------|-----------|--------|
| **基础运行时** | Node.js、Git、Python 3 | winget / 直接下载 |
| **AI CLI 工具** | Claude Code、Codex CLI、Echobird | npm 全局安装 |
| **编辑器** | VS Code、Cursor | winget |
| **终端增强** | Windows Terminal、Oh My Posh | winget |

- 每个工具装之前会先检测有没有装过，装过的自动跳过
- Git 支持从 GitHub 镜像直接下载安装包静默安装（不用手动点下一步）
- 安装过程实时显示日志

### 3. 一键升级

工具装完之后，能检测每个工具的最新版本，支持逐个升级或全部升级。

### 4. 卸载清理

只卸载本工具帮你装的东西，不会碰你自己原来装的软件。卸载后自动更新记录。

### 5. 环境验证

装完之后点「环境验证」，会逐个检测所有工具能不能正常调用，输出一张报告卡：

```
╔══════════════════════════════════════════════╗
║           环 境 验 证 报 告 卡              ║
╠══════════════════════════════════════════════╣
【基础运行时】
  ✅  Node.js LTS          v22.15.0
  ✅  Git                  v2.49.0
  ✅  Python 3             v3.12.3
【AI CLI 工具】
  ✅  Claude Code          v1.0.30
  ❌  Codex CLI            未安装
【编辑器】
  ✅  VS Code              1.99.2
  ✅  Cursor               0.49.6
╠══════════════════════════════════════════════╣
  环境检测完成：7/8 项通过
╚══════════════════════════════════════════════╝
```

---

## 怎么用

### 方式一：下载 exe（推荐）

1. 去 [GitHub Actions](https://github.com/tanzhijir-04/WinDevReady/actions) 页面
2. 点最新的构建记录 → 底部 Artifacts 区域下载 `WinDevReady-windows-amd64.zip`
3. 解压，双击 `WinDevReady.exe` 运行

### 方式二：自己编译

需要先装好 [Go](https://go.dev/dl/)（1.23+）和 Git：

```bash
git clone https://github.com/tanzhijir-04/WinDevReady.git
cd WinDevReady
go mod tidy
go build -o WinDevReady.exe -ldflags="-s -w" .
```

编译好的 `WinDevReady.exe` 就在当前目录。

---

## 怎么添加新工具

只需要在 `internal/config/tools.go` 文件里加一条记录，比如加一个 `npm` 方式的工具：

```go
{
    ID:          "my-tool",
    Name:        "我的工具",
    Group:       GroupAICLI,        // 选分组：GroupRuntime / GroupAICLI / GroupEditor / GroupTerminal
    Method:      MethodNpm,         // 选安装方式
    Package:     "my-tool-pkg",     // 包名
    VerifyCmd:   "my-tool -v",      // 验证命令（用来判断装没装成功）
    Description: "一句话说明这个工具是干嘛的",
},
```

支持的安装方式：

| 方式 | 说明 | 需要填写的字段 |
|------|------|---------------|
| `winget` | Windows 包管理器 | `WingetID` |
| `npm` | npm 全局安装 | `Package` |
| `pip` | pip 安装 | `Package` |
| `download` | 直接下载安装包 | `DownloadURLs` + `SilentArgs` |
| `choco` | Chocolatey | `Package` |

---

## 项目结构

```
WinDevReady/
├── main.go                      # 程序入口，启动 GUI
├── internal/
│   ├── config/tools.go          # 所有工具的配置（加新工具改这里就行）
│   ├── logger/logger.go         # 日志系统（实时输出到界面）
│   ├── store/records.go         # 安装记录（存在 AppData 里）
│   ├── network/detector.go      # 网络检测 + 镜像/代理自动切换
│   ├── installer/
│   │   ├── installer.go         # 安装引擎（核心）
│   │   ├── downloader.go        # HTTP 下载器（多源回退）
│   │   ├── uninstaller.go       # 卸载清理
│   │   └── updater.go           # 版本对比和升级
│   ├── verify/verifier.go       # 环境验证报告卡
│   └── ui/                      # 界面
│       ├── app.go               # 主窗口和侧边栏
│       ├── install_page.go      # 安装页面
│       ├── update_page.go       # 更新页面
│       ├── uninstall_page.go    # 卸载页面
│       └── report_page.go       # 报告卡页面
└── .github/workflows/build.yml  # GitHub Actions 自动编译
```

---

## 技术细节

- **语言**：Go 1.23
- **GUI**：[Fyne v2](https://fyne.io/) — 用 Go 写的跨平台 UI 框架
- **编译产物**：单个 `.exe` 文件，不依赖任何运行时，拷到对方电脑就能跑
- **CI/CD**：推代码到 GitHub 后自动用 Windows runner 编译

---

## 贡献

欢迎提 Issue 和 PR。最简单的贡献方式就是**加一个新工具**——照着上面的格式在 `config/tools.go` 里加一条就行。

## License

[MIT](LICENSE)

package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"WinDevReady/internal/installer"
	"WinDevReady/internal/logger"
	"WinDevReady/internal/network"
	"WinDevReady/internal/store"
	"WinDevReady/internal/verify"
)

// ============================================================
// 应用主窗口 —— 自定义主题 + 侧边栏 + 底部栏
// ============================================================

// App 应用主结构
type App struct {
	fyneApp    fyne.App
	mainWindow fyne.Window

	// 核心模块
	log     *logger.Logger
	store   *store.Records
	network *network.Detector
	inst    *installer.Installer
	verify  *verify.Verifier

	// 页面
	installPage   *InstallPage
	updatePage    *UpdatePage
	uninstallPage *UninstallPage
	reportPage    *ReportPage

	currentPage string
}

// NewApp 创建应用实例
func NewApp() *App {
	a := &App{}

	// 初始化 Fyne 应用 + 自定义主题
	a.fyneApp = app.NewWithID("com.windevready.app")
	a.fyneApp.Settings().SetTheme(NewPrimerTheme())

	// 创建主窗口
	a.mainWindow = a.fyneApp.NewWindow("WinDevReady")
	a.mainWindow.Resize(fyne.NewSize(1024, 680))
	a.mainWindow.SetFixedSize(false)

	// 加载自定义图标
	if iconRes, err := fyne.LoadResourceFromPath("assets/Icon.png"); err == nil {
		a.mainWindow.SetIcon(iconRes)
	}

	// 初始化核心模块
	a.log = logger.New()
	a.store, _ = store.NewRecords()
	a.network = network.New(a.log)
	a.inst = installer.New(a.log, a.store)
	a.verify = verify.New(a.log)

	return a
}

// Run 启动应用
func (a *App) Run() {
	content := a.buildLayout()
	a.mainWindow.SetContent(content)
	a.mainWindow.ShowAndRun()
}

// buildLayout 构建主界面布局：侧边栏 | 主内容 + 底部栏
func (a *App) buildLayout() fyne.CanvasObject {
	// 侧边栏
	sidebar := a.buildSidebar()

	// 主内容区
	mainContent := container.NewStack()

	// 初始化页面
	a.installPage = NewInstallPage(a.inst, a.log, a.network, mainContent)
	a.updatePage = NewUpdatePage(a.inst, a.log, a.store, mainContent)
	a.uninstallPage = NewUninstallPage(a.inst, a.log, a.store, mainContent)
	a.reportPage = NewReportPage(a.verify, a.log, mainContent)

	// 底部栏
	footer := a.buildFooter()

	// 右侧 = 主内容 + 底部栏
	rightSide := container.NewBorder(nil, footer, nil, nil, mainContent)

	// 整体布局：侧边栏 + 右侧
	root := container.NewBorder(nil, nil, sidebar, nil, rightSide)

	// 默认显示安装页
	a.installPage.Show()

	return root
}

// buildSidebar 构建侧边栏
func (a *App) buildSidebar() *fyne.Container {
	// 侧边栏按钮
	btnDefs := []struct {
		id    string
		label string
		icon  fyne.Resource
	}{
		{"install", "安装", theme.DownloadIcon()},
		{"update", "更新", theme.ViewRefreshIcon()},
		{"uninstall", "卸载", theme.DeleteIcon()},
		{"report", "报告", theme.ConfirmIcon()},
	}

	// Logo 区域
	logo := canvas.NewText("⚡ WinDevReady", PrimerColors.Primary)
	logo.TextSize = 18
	logo.TextStyle = fyne.TextStyle{Bold: true}
	logoBox := container.NewPadded(logo)

	// 版本标签
	version := canvas.NewText("v1.0.0", PrimerColors.TextMuted)
	version.TextSize = 11
	versionBox := container.NewPadded(version)

	// 分隔线
	sep := canvas.NewRectangle(PrimerColors.Border)
	sep.SetMinSize(fyne.NewSize(0, 1))

	var items []fyne.CanvasObject
	items = append(items, logoBox, sep, versionBox)

	for _, def := range btnDefs {
		def := def
		btn := widgetNewSidebarBtn(def.label, def.icon, func() {
			a.switchPage(def.id)
		})
		items = append(items, btn)
	}

	// 底部填充
	items = append(items, layout.NewSpacer())

	// 侧边栏背景
 sidebarBg := canvas.NewRectangle(PrimerColors.Sidebar)
 sidebarContent := container.NewVBox(items...)

	return container.NewStack(sidebarBg, sidebarContent)
}

// buildFooter 构建底部栏（可点击仓库链接 + 赞助按钮）
func (a *App) buildFooter() fyne.CanvasObject {
	// 分隔线
	sep := canvas.NewRectangle(PrimerColors.Border)
	sep.SetMinSize(fyne.NewSize(0, 1))

	// 仓库地址 —— 可点击跳转浏览器
	repoBtn := widget.NewButtonWithIcon(
		"github.com/tanzhijir-04/WinDevReady",
		theme.ComputerIcon(),
		func() {
			openURL("https://github.com/tanzhijir-04/WinDevReady")
		},
	)
	repoBtn.Importance = widget.LowImportance

	// 赞助按钮 —— 爱发电
	sponsorBtn := widget.NewButtonWithIcon(
		"❤️ 爱发电赞助",
		theme.ContentAddIcon(),
		func() {
			openURL("https://ifdian.net/a/tanz666")
		},
	)
	sponsorBtn.Importance = widget.WarningImportance

	// 版本号
	verText := canvas.NewText("v1.0.0", PrimerColors.TextMuted)
	verText.TextSize = 11

	footerContent := container.NewHBox(repoBtn, layout.NewSpacer(), sponsorBtn, layout.NewSpacer(), verText)
	footerPadded := container.NewPadded(footerContent)

	return container.NewBorder(sep, nil, nil, nil, footerPadded)
}

// openURL 跨平台打开浏览器
func openURL(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	fyne.CurrentApp().OpenURL(u)
}

// switchPage 切换页面
func (a *App) switchPage(pageID string) {
	a.currentPage = pageID
	switch pageID {
	case "install":
		a.installPage.Show()
	case "update":
		a.updatePage.Show()
	case "uninstall":
		a.uninstallPage.Show()
	case "report":
		a.reportPage.Show()
	}
}

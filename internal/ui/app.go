package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"WinDevReady/internal/installer"
	"WinDevReady/internal/logger"
	"WinDevReady/internal/network"
	"WinDevReady/internal/store"
	"WinDevReady/internal/verify"
)

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

	// 侧边栏按钮
	sidebarBtns []*fyne.Container
	currentPage string
}

// NewApp 创建应用实例
func NewApp() *App {
	a := &App{}

	// 初始化 Fyne 应用
	a.fyneApp = app.New()
	a.mainWindow = a.fyneApp.NewWindow("WinDevReady - Windows AI 开发环境配置工具")
	a.mainWindow.Resize(fyne.NewSize(960, 640))

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

// buildLayout 构建主界面布局
func (a *App) buildLayout() fyne.CanvasObject {
	// 侧边栏
	sidebar := a.buildSidebar()

	// 主内容区（默认显示安装页）
	mainContent := container.NewStack()

	// 页面容器
	pages := container.NewBorder(nil, nil, sidebar, nil, mainContent)

	// 设置默认页面
	a.installPage = NewInstallPage(a.inst, a.log, a.network, mainContent)
	a.updatePage = NewUpdatePage(a.inst, a.log, a.store, mainContent)
	a.uninstallPage = NewUninstallPage(a.inst, a.log, a.store, mainContent)
	a.reportPage = NewReportPage(a.verify, a.log, mainContent)

	// 默认显示安装页
	a.installPage.Show()

	return pages
}

// buildSidebar 构建侧边栏导航
func (a *App) buildSidebar() *fyne.Container {
	// 页面按钮定义
	btnDefs := []struct {
		id    string
		label string
		icon  fyne.Resource
	}{
		{"install", "安装工具", theme.DownloadIcon()},
		{"update", "版本更新", theme.ViewRefreshIcon()},
		{"uninstall", "卸载清理", theme.DeleteIcon()},
		{"report", "环境报告", theme.ConfirmIcon()},
	}

	var buttons []*fyne.Container
	for _, def := range btnDefs {
		def := def // 捕获循环变量
		btn := createSidebarButton(def.label, def.icon, func() {
			a.switchPage(def.id)
		})
		buttons = append(buttons, btn)
	}

	// Logo + 标题
	title := newTitle("WinDevReady")

	sidebarItems := []fyne.CanvasObject{title}
	for _, btn := range buttons {
		sidebarItems = append(sidebarItems, btn)
		sidebarItems = append(sidebarItems, layout.NewSpacer())
	}

	return container.NewVBox(sidebarItems...)
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

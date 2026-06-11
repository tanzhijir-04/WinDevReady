package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"WinDevReady/internal/config"
	"WinDevReady/internal/installer"
	"WinDevReady/internal/logger"
	"WinDevReady/internal/network"
)

// InstallPage 安装页面
type InstallPage struct {
	parent    *fyne.Container
	content   fyne.CanvasObject
	inst      *installer.Installer
	log       *logger.Logger
	net       *network.Detector
	checkboxes map[string]*widget.Check // toolID -> checkbox
	logEntry   *widget.TextEntry
}

// NewInstallPage 创建安装页面
func NewInstallPage(inst *installer.Installer, log *logger.Logger, net *network.Detector, parent *fyne.Container) *InstallPage {
	p := &InstallPage{
		parent:    parent,
		inst:      inst,
		log:       log,
		net:       net,
		checkboxes: make(map[string]*widget.Check),
	}
	p.build()
	return p
}

// build 构建安装页面内容
func (p *InstallPage) build() {
	// 工具勾选区
	toolSection := p.buildToolSection()

	// 操作按钮
	actionBar := p.buildActionBar()

	// 日志输出区
	p.logEntry = newLogArea()
	logCard := newGroupCard("安装日志", p.logEntry)

	// 组装布局
	p.content = container.NewBorder(
		toolSection,  // 上：工具勾选
		actionBar,    // 下：操作按钮
		nil, nil,
		logCard,      // 中：日志输出
	)
}

// buildToolSection 构建工具勾选区
func (p *InstallPage) buildToolSection() fyne.CanvasObject {
	var groupCards []fyne.CanvasObject

	for _, group := range config.AllGroups {
		tools := config.GetToolsByGroup(group)
		var checkboxContainer []fyne.CanvasObject

		for _, tool := range tools {
			tool := tool
			cb := widget.NewCheck(tool.Name, nil)
			cb.OnChanged = func(checked bool) {
				// 预留：可在此添加实时反馈
			}
			p.checkboxes[tool.ID] = cb
			checkboxContainer = append(checkboxContainer, cb)
		}

		card := newGroupCard(string(group), container.NewVBox(checkboxContainer...))
		groupCards = append(groupCards, card)
	}

	return container.NewHBox(groupCards...)
}

// buildActionBar 构建操作按钮区
func (p *InstallPage) buildActionBar() fyne.CanvasObject {
	// 全选/全不选
	selectAllBtn := widget.NewButton("全选", func() {
		for _, cb := range p.checkboxes {
			cb.SetChecked(true)
		})
	})
	deselectAllBtn := widget.NewButton("全不选", func() {
		for _, cb := range p.checkboxes {
			cb.SetChecked(false)
		})
	})

	// 一键安装
	installBtn := widget.NewButton("一键安装", func() {
		go p.startInstall()
	})
	installBtn.Importance = widget.HighImportance

	return container.NewHBox(
		selectAllBtn,
		deselectAllBtn,
		layout.NewSpacer(),
		installBtn,
	)
}

// startInstall 开始安装（在 goroutine 中执行）
func (p *InstallPage) startInstall() {
	// 订阅日志
	logCh := p.log.Subscribe()
	defer p.log.Unsubscribe(logCh)

	// 异步更新日志 UI
	go func() {
		for entry := range logCh {
			p.logEntry.SetText(p.logEntry.Text + logger.FormatEntry(entry) + "\n")
			// 滚动到底部
			p.logEntry.CursorRow = len(p.logEntry.Text)
		}
	}()

	// 网络检测
	p.log.Info("", "========== 开始网络检测 ==========")
	status := p.net.CheckAll()

	// 应用网络配置
	if status.RegistryURL != "" {
		_ = p.net.ApplyNPMRegistry(status.RegistryURL)
	}
	if status.ProxyActive && status.ProxyPort != "" {
		_ = p.net.ApplyProxy(status.ProxyPort)
	}

	// 收集用户勾选的工具
	var selectedTools []config.ToolDef
	for id, cb := range p.checkboxes {
		if cb.Checked {
			if tool, ok := config.GetToolByID(id); ok {
				selectedTools = append(selectedTools, tool)
			}
		}
	}

	if len(selectedTools) == 0 {
		p.log.Warn("", "未选择任何工具")
		return
	}

	// 执行安装
	p.log.Info("", "========== 开始安装 ==========")
	success, failed := p.inst.InstallTools(selectedTools)
	p.log.Success("", "========== 安装完成 ==========")
	p.log.Success("", "成功: "+string(rune(success+'0'))+" 失败: "+string(rune(failed+'0')))
}

// Show 显示安装页面
func (p *InstallPage) Show() {
	p.parent.Objects = []fyne.CanvasObject{p.content}
	p.parent.Refresh()
}

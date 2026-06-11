package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"WinDevReady/internal/installer"
	"WinDevReady/internal/logger"
	"WinDevReady/internal/store"
)

// UninstallPage 卸载页面
type UninstallPage struct {
	parent   *fyne.Container
	content  fyne.CanvasObject
	inst     *installer.Installer
	log      *logger.Logger
	store    *store.Records
	logEntry *widget.Entry
	list     *widget.List
	records  []store.Record
}

// NewUninstallPage 创建卸载页面
func NewUninstallPage(inst *installer.Installer, log *logger.Logger, store *store.Records, parent *fyne.Container) *UninstallPage {
	p := &UninstallPage{
		parent: parent,
		inst:   inst,
		log:    log,
		store:  store,
	}
	p.build()
	return p
}

// build 构建卸载页面
func (p *UninstallPage) build() {
	// 刷新按钮
	refreshBtn := widget.NewButton("刷新列表", func() {
		p.refreshList()
	})

	// 全部卸载
	uninstallAllBtn := widget.NewButton("全部卸载", func() {
		go p.uninstallAll()
	})
	uninstallAllBtn.Importance = widget.DangerImportance

	actionBar := container.NewHBox(refreshBtn, layout.NewSpacer(), uninstallAllBtn)

	// 已安装工具列表
	p.list = widget.NewList(
		func() int { return len(p.records) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("工具名称"),
				widget.NewLabel("版本"),
				widget.NewLabel("安装时间"),
				widget.NewButton("卸载", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(p.records) {
				return
			}
			rec := p.records[id]
			box := obj.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(rec.Name)
			box.Objects[1].(*widget.Label).SetText(rec.Version)
			box.Objects[2].(*widget.Label).SetText(rec.InstalledAt.Format("2006-01-02 15:04"))
			btn := box.Objects[3].(*widget.Button)
			btn.OnTapped = func() {
				go p.uninstallOne(rec.ToolID)
			}
		},
	)

	// 日志区
	p.logEntry = newLogArea()
	logCard := newGroupCard("卸载日志", p.logEntry)

	listCard := newGroupCard("已安装工具", p.list)

	p.content = container.NewBorder(
		actionBar, nil, nil, nil,
		container.NewBorder(listCard, nil, nil, nil, logCard),
	)

	// 初始加载
	p.refreshList()
}

// refreshList 刷新已安装工具列表
func (p *UninstallPage) refreshList() {
	p.records = p.store.GetAll()
	p.list.Refresh()
}

// uninstallOne 卸载单个工具
func (p *UninstallPage) uninstallOne(toolID string) {
	logCh := p.log.Subscribe()
	defer p.log.Unsubscribe(logCh)

	go func() {
		for entry := range logCh {
			p.logEntry.SetText(p.logEntry.Text + logger.FormatEntry(entry) + "\n")
		}
	}()

	_ = p.inst.UninstallTool(toolID)
	p.refreshList()
}

// uninstallAll 卸载全部
func (p *UninstallPage) uninstallAll() {
	logCh := p.log.Subscribe()
	defer p.log.Unsubscribe(logCh)

	go func() {
		for entry := range logCh {
			p.logEntry.SetText(p.logEntry.Text + logger.FormatEntry(entry) + "\n")
		}
	}()

	p.log.Info("", "========== 开始卸载清理 ==========")
	success, failed := p.inst.UninstallAll()
	p.log.Success("", "========== 卸载完成 ==========")
	p.log.Success("", "成功卸载: "+string(rune(success+'0'))+" 个工具，失败: "+string(rune(failed+'0')))
	p.refreshList()
}

// Show 显示卸载页面
func (p *UninstallPage) Show() {
	p.parent.Objects = []fyne.CanvasObject{p.content}
	p.parent.Refresh()
}

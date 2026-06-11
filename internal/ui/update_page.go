package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"WinDevReady/internal/installer"
	"WinDevReady/internal/logger"
	"WinDevReady/internal/store"
)

// UpdatePage 更新页面
type UpdatePage struct {
	parent   *fyne.Container
	content  fyne.CanvasObject
	inst     *installer.Installer
	log      *logger.Logger
	store    *store.Records
	logEntry *widget.Entry
	list     *widget.List
	updates  []installer.UpdateInfo
}

// NewUpdatePage 创建更新页面
func NewUpdatePage(inst *installer.Installer, log *logger.Logger, store *store.Records, parent *fyne.Container) *UpdatePage {
	p := &UpdatePage{
		parent: parent,
		inst:   inst,
		log:    log,
		store:  store,
	}
	p.build()
	return p
}

// build 构建更新页面
func (p *UpdatePage) build() {
	// 操作按钮
	checkBtn := widget.NewButton("检查更新", func() {
		go p.checkUpdates()
	})
	checkBtn.Importance = widget.HighImportance

	updateAllBtn := widget.NewButton("全部升级", func() {
		go p.updateAll()
	})

	actionBar := container.NewHBox(checkBtn, updateAllBtn, layout.NewSpacer())

	// 更新列表
	p.list = widget.NewList(
		func() int { return len(p.updates) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("工具名称"),
				widget.NewLabel("当前版本"),
				widget.NewLabel("最新版本"),
				widget.NewButton("升级", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(p.updates) {
				return
			}
			info := p.updates[id]
			box := obj.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(info.ToolName)
			box.Objects[1].(*widget.Label).SetText(info.CurrentVersion)
			box.Objects[2].(*widget.Label).SetText(info.LatestVersion)
			btn := box.Objects[3].(*widget.Button)
			if info.HasUpdate {
				btn.Enable()
				btn.OnTapped = func() {
					go p.updateOne(info.ToolID)
				}
			} else {
				btn.SetText("已是最新")
				btn.Disable()
			}
		},
	)

	// 日志区
	p.logEntry = newLogArea()
	logCard := newGroupCard("更新日志", p.logEntry)

	listCard := newGroupCard("可用更新", p.list)

	p.content = container.NewBorder(
		actionBar, nil, nil, nil,
		container.NewBorder(listCard, nil, nil, nil, logCard),
	)
}

// checkUpdates 检查更新
func (p *UpdatePage) checkUpdates() {
	p.log.Info("", "正在检查更新...")
	p.updates = p.inst.CheckUpdates()
	p.list.Refresh()
	p.log.Success("", fmt.Sprintf("检查完成，发现 %d 个可更新工具", len(p.updates)))
}

// updateAll 升级全部
func (p *UpdatePage) updateAll() {
	logCh := p.log.Subscribe()
	defer p.log.Unsubscribe(logCh)

	go func() {
		for entry := range logCh {
			p.logEntry.SetText(p.logEntry.Text + logger.FormatEntry(entry) + "\n")
		}
	}()

	p.log.Info("", "========== 开始批量升级 ==========")
	for _, info := range p.updates {
		if info.HasUpdate {
			_ = p.inst.UpdateTool(info.ToolID)
		}
	}
	p.log.Success("", "========== 批量升级完成 ==========")

	// 刷新列表
	p.updates = p.inst.CheckUpdates()
	p.list.Refresh()
}

// updateOne 升级单个工具
func (p *UpdatePage) updateOne(toolID string) {
	logCh := p.log.Subscribe()
	defer p.log.Unsubscribe(logCh)

	go func() {
		for entry := range logCh {
			p.logEntry.SetText(p.logEntry.Text + logger.FormatEntry(entry) + "\n")
		}
	}()

	_ = p.inst.UpdateTool(toolID)

	// 刷新列表
	p.updates = p.inst.CheckUpdates()
	p.list.Refresh()
}

// Show 显示更新页面
func (p *UpdatePage) Show() {
	p.parent.Objects = []fyne.CanvasObject{p.content}
	p.parent.Refresh()
}

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"WinDevReady/internal/logger"
	"WinDevReady/internal/verify"
)

// ReportPage 环境报告页面
type ReportPage struct {
	parent    *fyne.Container
	content   fyne.CanvasObject
	verifier  *verify.Verifier
	log       *logger.Logger
	logEntry  *widget.Entry
	reportLabel *widget.Label
}

// NewReportPage 创建报告页面
func NewReportPage(v *verify.Verifier, log *logger.Logger, parent *fyne.Container) *ReportPage {
	p := &ReportPage{
		parent:   parent,
		verifier: v,
		log:      log,
	}
	p.build()
	return p
}

// build 构建报告页面
func (p *ReportPage) build() {
	// 操作按钮
	verifyBtn := widget.NewButton("开始验证", func() {
		go p.runVerification()
	})
	verifyBtn.Importance = widget.HighImportance

	actionBar := container.NewHBox(verifyBtn, layout.NewSpacer())

	// 报告卡显示区
	p.reportLabel = widget.NewLabel("点击「开始验证」检测环境状态")
	p.reportLabel.Wrapping = fyne.TextWrapBreak

	reportCard := newGroupCard("环境报告卡", container.NewScroll(p.reportLabel))

	// 日志区
	p.logEntry = newLogArea()
	logCard := newGroupCard("验证日志", p.logEntry)

	p.content = container.NewBorder(
		actionBar, nil, nil, nil,
		container.NewBorder(reportCard, nil, nil, nil, logCard),
	)
}

// runVerification 执行环境验证
func (p *ReportPage) runVerification() {
	logCh := p.log.Subscribe()
	defer p.log.Unsubscribe(logCh)

	go func() {
		for entry := range logCh {
			p.logEntry.SetText(p.logEntry.Text + logger.FormatEntry(entry) + "\n")
		}
	}()

	p.log.Info("", "========== 开始环境验证 ==========")
	report := p.verifier.VerifyAll()
	p.log.Success("", "========== 验证完成 ==========")

	// 显示报告卡
	p.reportLabel.SetText(verify.FormatReport(report))
}

// Show 显示报告页面
func (p *ReportPage) Show() {
	p.parent.Objects = []fyne.CanvasObject{p.content}
	p.parent.Refresh()
}

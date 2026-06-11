package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// widgetNewSidebarBtn 创建侧边栏按钮（深色风格）
func widgetNewSidebarBtn(label string, icon fyne.Resource, onTap func()) fyne.CanvasObject {
	btn := widget.NewButtonWithIcon(label, icon, onTap)
	btn.Importance = widget.LowImportance
	btn.Alignment = widget.ButtonAlignLeading

	// 按钮背景
	bg := canvas.NewRectangle(PrimerColors.Sidebar)
	bg.SetMinSize(fyne.NewSize(180, 40))

	return container.NewStack(bg, container.NewPadded(btn))
}

// newTitle 创建应用标题组件
func newTitle(text string) fyne.CanvasObject {
	title := canvas.NewText(text, color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 20
	return container.NewCenter(title)
}

// newLogArea 创建日志输出区域
func newLogArea() *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapBreak
	entry.Disable()
	return entry
}

// newGroupCard 创建分组卡片
func newGroupCard(title string, content fyne.CanvasObject) *widget.Card {
	return widget.NewCard(title, "", content)
}

// newCheckboxList 创建可勾选的工具列表
type ToolCheckbox struct {
	ID       string
	Name     string
	CheckBox *widget.Check
}

// createToolCheckboxList 为指定分组创建勾选列表
func createToolCheckboxList(items []ToolCheckbox) *fyne.Container {
	var checkboxes []fyne.CanvasObject
	for i := range items {
		checkboxes = append(checkboxes, items[i].CheckBox)
		if i < len(items)-1 {
			checkboxes = append(checkboxes, layout.NewSpacer())
		}
	}
	return container.NewVBox(checkboxes...)
}

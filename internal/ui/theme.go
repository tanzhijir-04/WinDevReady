package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ============================================================
// 自定义主题 —— 深色科技风格 + 自定义字体加载
// ============================================================

// PrimerTheme 自定义主题
type PrimerTheme struct {
	normalFont fyne.Resource
	monoFont   fyne.Resource
}

// NewPrimerTheme 创建自定义主题，fontDir 为字体目录路径
func NewPrimerTheme() *PrimerTheme {
	return &PrimerTheme{
		normalFont: theme.DefaultTheme().Font(fyne.TextStyle{}),
		monoFont:   theme.DefaultTheme().Font(fyne.TextStyle{Monospace: true}),
	}
}

// Font 实现 fyne.Theme 接口 —— 根据样式返回对应字体
func (t *PrimerTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return t.monoFont
	}
	return t.normalFont
}

// Color 实现 fyne.Theme 接口 —— 深色配色方案
func (t *PrimerTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 24, G: 28, B: 38, A: 255} // 深蓝灰底色
	case theme.ColorNameButton:
		return color.NRGBA{R: 45, G: 52, B: 68, A: 255} // 按钮底色
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 35, G: 40, B: 52, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 220, G: 225, B: 235, A: 255} // 浅色文字
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 120, G: 130, B: 150, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 55, G: 65, B: 85, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 32, G: 36, B: 48, A: 255}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 60, G: 70, B: 90, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 88, G: 166, B: 255, A: 255} // 主题蓝
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 80, G: 200, B: 120, A: 255} // 成功绿
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 180, B: 50, A: 255} // 警告橙
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 80, B: 80, A: 255}  // 错误红
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 20, G: 24, B: 34, A: 255}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 20, G: 24, B: 34, A: 255}
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 30, G: 34, B: 46, A: 255}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 50, G: 58, B: 75, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Icon 实现 fyne.Theme 接口 —— 使用默认图标集
func (t *PrimerTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size 实现 fyne.Theme 接口 —— 微调尺寸
func (t *PrimerTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNamePadding:
		return 12
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameInputBorder:
		return 1
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// PrimerColors 导出配色常量供 UI 使用
var PrimerColors = struct {
	BG         color.NRGBA // 背景
	Sidebar    color.NRGBA // 侧边栏
	Card       color.NRGBA // 卡片
	Primary    color.NRGBA // 主题蓝
	Success    color.NRGBA // 成功绿
	Warning    color.NRGBA // 警告橙
	Error      color.NRGBA // 错误红
	Text       color.NRGBA // 主文字
	TextMuted  color.NRGBA // 次要文字
	Border     color.NRGBA // 边框
}{
	BG:        color.NRGBA{R: 24, G: 28, B: 38, A: 255},
	Sidebar:   color.NRGBA{R: 18, G: 22, B: 30, A: 255},
	Card:      color.NRGBA{R: 30, G: 34, B: 46, A: 255},
	Primary:   color.NRGBA{R: 88, G: 166, B: 255, A: 255},
	Success:   color.NRGBA{R: 80, G: 200, B: 120, A: 255},
	Warning:   color.NRGBA{R: 255, G: 180, B: 50, A: 255},
	Error:     color.NRGBA{R: 255, G: 80, B: 80, A: 255},
	Text:      color.NRGBA{R: 220, G: 225, B: 235, A: 255},
	TextMuted: color.NRGBA{R: 120, G: 130, B: 150, A: 255},
	Border:    color.NRGBA{R: 50, G: 58, B: 75, A: 255},
}

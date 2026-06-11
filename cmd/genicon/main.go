// genicon —— 生成 WinDevReady 应用图标（256x256 PNG）
// 用法: go run ./cmd/genicon
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

const (
	size    = 256
	padding = 20
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// 背景：深蓝灰圆角矩形
	bgColor := color.NRGBA{R: 24, G: 28, B: 38, A: 255}
	drawRoundedRect(img, 0, 0, size, size, 40, bgColor)

	// 中心圆形高光
	highlight := color.NRGBA{R: 88, G: 166, B: 255, A: 40}
	drawCircle(img, size/2, size/2, 80, highlight)

	// "W" 字母 —— 用像素画
	drawW(img, 60, 70, 136, 120, color.NRGBA{R: 88, G: 166, B: 255, A: 255})

	// 底部横线装饰
	lineColor := color.NRGBA{R: 80, G: 200, B: 120, A: 255}
	for x := 70; x < 186; x++ {
		for dy := 0; dy < 4; dy++ {
			img.Set(x, 210+dy, lineColor)
		}
	}

	// 保存
	f, err := os.Create("assets/icon.png")
	if err != nil {
		fmt.Println("创建文件失败:", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Println("编码失败:", err)
		os.Exit(1)
	}
	fmt.Println("✅ 图标已生成: assets/icon.png (256x256)")
}

// drawRoundedRect 绘制圆角矩形
func drawRoundedRect(img *image.RGBA, x0, y0, x1, y1, r int, c color.Color) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if isInsideRoundedRect(x, y, x0, y0, x1, y1, r) {
				img.Set(x, y, c)
			}
		}
	}
}

func isInsideRoundedRect(x, y, x0, y0, x1, y1, r int) bool {
	if x < x0 || x >= x1 || y < y0 || y >= y1 {
		return false
	}
	// 四个角的圆角检测
	corners := [][2]int{{x0 + r, y0 + r}, {x1 - r - 1, y0 + r}, {x0 + r, y1 - r - 1}, {x1 - r - 1, y1 - r - 1}}
	for _, c := range corners {
		if (x < x0+r || x >= x1-r) && (y < y0+r || y >= y1-r) {
			dx := float64(x - c[0])
			dy := float64(y - c[1])
			if math.Sqrt(dx*dx+dy*dy) > float64(r) {
				return false
			}
		}
	}
	return true
}

// drawCircle 绘制半透明圆形
func drawCircle(img *image.RGBA, cx, cy, radius int, c color.Color) {
	r := c.(color.NRGBA)
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= float64(radius) {
				// 距离越远越透明
				alpha := uint8(float64(r.A) * (1.0 - dist/float64(radius)))
				img.Set(x, y, color.NRGBA{R: r.R, G: r.G, B: r.B, A: alpha})
			}
		}
	}
}

// drawW 用线段绘制 "W" 字母
func drawW(img *image.RGBA, x, y, w, h int, c color.Color) {
	// W 的五个关键点
	points := [5][2]float64{
		{float64(x), float64(y)},           // 左上
		{float64(x + w/4), float64(y + h)}, // 左下中
		{float64(x + w/2), float64(y + h*2/3)}, // 中间
		{float64(x + w*3/4), float64(y + h)},   // 右下中
		{float64(x + w), float64(y)},            // 右上
	}

	thickness := 8.0
	// 绘制四条线段
	for i := 0; i < 4; i++ {
		drawThickLine(img, points[i], points[i+1], thickness, c)
	}
}

// drawThickLine 画粗线
func drawThickLine(img *image.RGBA, from, to [2]float64, thickness float64, c color.Color) {
	dx := to[0] - from[0]
	dy := to[1] - from[1]
	length := math.Sqrt(dx*dx + dy*dy)
	if length == 0 {
		return
	}
	// 法线方向
	nx := -dy / length * thickness / 2
	ny := dx / length * thickness / 2

	steps := int(length * 2) // 两倍采样保证不漏
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		cx := from[0] + dx*t
		cy := from[1] + dy*t
		// 在法线方向扩展
		for _, sign := range []float64{-1, 0, 1} {
			px := int(cx + nx*sign)
			py := int(cy + ny*sign)
			img.Set(px, py, c)
		}
	}
}

// ensure draw is imported (used by drawCircle via draw.Draw)
var _ = draw.Draw

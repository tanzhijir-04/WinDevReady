# 生成 WinDevReady 图标 (256x256 PNG) - 简化版
Add-Type -AssemblyName System.Drawing

$size = 256
$bmp = New-Object System.Drawing.Bitmap($size, $size)
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias

# 背景
$g.Clear([System.Drawing.Color]::FromArgb(24, 28, 38))

# 圆形高光
$glow = New-Object System.Drawing.Drawing2D.LinearGradientBrush(
    (New-Object System.Drawing.Point(128, 30)),
    (New-Object System.Drawing.Point(128, 230)),
    [System.Drawing.Color]::FromArgb(80, 88, 166, 255),
    [System.Drawing.Color]::FromArgb(10, 88, 166, 255)
)
$g.FillEllipse($glow, 30, 30, 196, 196)

# "W" 字母 - 用贝塞尔曲线
$wBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb(88, 166, 255))
$wPen = New-Object System.Drawing.Pen($wBrush, 10)
$wPen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
$wPen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round

$pts = [System.Drawing.PointF[]]@(
    (New-Object System.Drawing.PointF(60.0, 75.0)),
    (New-Object System.Drawing.PointF(100.0, 170.0)),
    (New-Object System.Drawing.PointF(128.0, 125.0)),
    (New-Object System.Drawing.PointF(156.0, 170.0)),
    (New-Object System.Drawing.PointF(196.0, 75.0))
)
$g.DrawLines($wPen, $pts)

# 底部绿线
$green = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb(80, 200, 120))
$g.FillRectangle($green, 60, 208, 136, 4)

# 保存
$bmp.Save("assets/icon.png", [System.Drawing.Imaging.ImageFormat]::Png)
$g.Dispose()
$bmp.Dispose()
Write-Host "Done: assets/icon.png"

package tank

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	mplusNormalFont *text.GoTextFaceSource
)

func init() {

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont = s
}

// 绘制坦克周围的 按键特效
func KeyPressDrawAroundTank(t *Tank, screen *ebiten.Image) {

	op := &text.DrawOptions{}

	op.ColorScale.ScaleWithColor(color.RGBA{210, 105, 30, 255})
	keyWord := ""
	x, y := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		x, y = -5.0, -25.0
		keyWord = "W"
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		x, y = -5.0, 25.0
		keyWord = "S"
	}

	op.GeoM.Translate(x, y)
	angleRad := t.Angle * math.Pi / 180.0 // 角度转弧度
	op.GeoM.Rotate(angleRad)
	// x,y 经过旋转 angleRad 角度后的位置坐标 x1,y1
	x1, y1 := x*math.Cos(angleRad)-y*math.Sin(angleRad), x*math.Sin(angleRad)+y*math.Cos(angleRad)
	//op.LineSpacing = 100
	op.GeoM.Translate(x1+t.X, y1+t.Y)
	text.Draw(screen, keyWord, &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   20}, op)

	// 重置
	op.GeoM.Reset()
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		x, y = -30.0, -5.0
		keyWord = "A"
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		x, y = 20.0, -5.0
		keyWord = "D"
	}

	op.GeoM.Translate(x, y)
	op.GeoM.Rotate(angleRad)
	x1, y1 = x*math.Cos(angleRad)-y*math.Sin(angleRad), x*math.Sin(angleRad)+y*math.Cos(angleRad)
	op.LineSpacing = 100
	op.GeoM.Translate(x1+t.X, y1+t.Y)
	text.Draw(screen, keyWord, &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   20}, op)
}

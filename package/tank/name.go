package tank

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type killedName struct {
	name       string
	updateTime time.Time
}

var (
	killedNames = make([]killedName, 0, 3)
)

// 更新名字列表
func UpdateNameList(name string) {

	kn := killedName{
		name:       name,
		updateTime: time.Now(),
	}
	killedNames = append(killedNames, kn)
}

// 绘制名字列表
func DrawNameList(screen *ebiten.Image) {

	x, y := 0., 25.
	for i, killeName := range killedNames {
		if time.Since(killeName.updateTime) > 5*time.Second {
			killedNames = append(killedNames[:i], killedNames[i+1:]...) // 去掉
			continue
		}
		drawName(screen, x, y, killeName.name)
		y += 50.
	}
}

var (
	mplusSmallFontFace font.Face
)

func init() {
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	mplusSmallFontFace, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 1,
		DPI:  72,
	})
}
func drawName(screen *ebiten.Image, x, y float64, txt string) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.RGBA{255, 97, 3, 255})
	op.GeoM.Translate(x, y)

	text.Draw(screen, txt, &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   50}, op)

	op.GeoM.Reset()
	op.ColorScale.Reset()

	bounds, _ := font.BoundString(mplusSmallFontFace, txt)
	width := float64(bounds.Max.X - bounds.Min.X)

	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(x+width, y)
	text.Draw(screen, " is killed", &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   50}, op)
}

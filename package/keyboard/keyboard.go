package keyboard

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
func Draw(screen *ebiten.Image) {

	vector.StrokeRect(screen, 500, 500, 25, 25, 1, color.Black, true)
	vector.DrawFilledRect(screen, 500+1, 500+1, 25-2, 25-2, color.Black, true)

	op := &text.DrawOptions{}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		op.ColorScale.ScaleWithColor(color.RGBA{0, 255, 0, 255})
	} else {
		op.ColorScale.ScaleWithColor(color.White)
	}

	op.GeoM.Translate(500+1, 500)
	op.LineSpacing = 100
	text.Draw(screen, "W", &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   23}, op)
}

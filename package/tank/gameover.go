package tank

import (
	"image/color"
	"os"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/gofish2020/tankgame/package/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFontFace font.Face
)

func init() {
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	mplusNormalFontFace, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 1,
		DPI:  100,
	})
}
func GameOverDraw(screen *ebiten.Image) {

	if utils.GameProgress == "over" {
		screen.Fill(color.Black)

		bounds, _ := font.BoundString(mplusNormalFontFace, "Game Over")
		width := float64(bounds.Max.X - bounds.Min.X)

		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(color.White)
		op.GeoM.Translate(monitor.ScreenWidth/2-width/2.0, monitor.ScreenHeight/2-90)

		text.Draw(screen, "Game Over", &text.GoTextFace{
			Source: mplusNormalFont,
			Size:   90}, op)

		op.GeoM.Reset()
		bounds, _ = font.BoundString(mplusNormalFontFace, "Press [Enter] to try again")
		width = float64(bounds.Max.X - bounds.Min.X)
		op.GeoM.Translate(monitor.ScreenWidth/2-width/2.0, monitor.ScreenHeight/2)
		text.Draw(screen, "Press [Enter] to try again", &text.GoTextFace{
			Source: mplusNormalFont,
			Size:   90}, op)

		op.GeoM.Reset()
		bounds, _ = font.BoundString(mplusNormalFontFace, "Press [ESC] to exit game")
		width = float64(bounds.Max.X - bounds.Min.X)
		op.GeoM.Translate(monitor.ScreenWidth/2-width/2.0, monitor.ScreenHeight/2+90)
		text.Draw(screen, "Press [ESC] to exit game", &text.GoTextFace{
			Source: mplusNormalFont,
			Size:   90}, op)
	}

}

func GameOverUpdate() {

	if utils.GameProgress == "over" {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			utils.GameProgress = "prepare"
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
	}
}

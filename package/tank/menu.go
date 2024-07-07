package tank

import (
	"fmt"
	"image/color"
	"os"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// 按钮坐标
type Coordinates struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

var (
	playButtonImg, exitButtonImg, logoImg *ebiten.Image

	playButton Coordinates
	exitButton Coordinates
)

func init() {
	// 按钮图片
	playButtonImage, _, _ := ebitenutil.NewImageFromFile("resource/play_button.png")
	playButtonImg = playButtonImage

	exitButtonImage, _, _ := ebitenutil.NewImageFromFile("resource/exit_game_button.png")
	exitButtonImg = exitButtonImage

	logoImage, _, _ := ebitenutil.NewImageFromFile("resource/logo.png")
	logoImg = logoImage

	// 按钮坐标
	playButton = Coordinates{
		X:      monitor.ScreenWidth - 250.0,
		Y:      0,
		Width:  250,
		Height: 74,
	}

	exitButton = Coordinates{
		X:      0,
		Y:      monitor.ScreenHeight - 74.0,
		Width:  250.0,
		Height: 74.0,
	}
}

// **************************** 更新主菜单坐标 *********************
func MenuUpdate(tk []*Tank) {

	playButton.Y += 3
	if playButton.Y >= monitor.ScreenHeight {
		playButton.Y = 0
	}

	exitButton.X += 5
	if exitButton.X >= monitor.ScreenWidth {
		exitButton.X = 0
	}

	for _, t := range tk {
		for _, projectile := range t.Projectiles {

			if checkForCollisions(projectile.X, projectile.Y, playButton.X, playButton.Y, playButton.X+playButton.Width, playButton.Y+playButton.Height) {
				fmt.Println("goal!!!!!")
				projectile.IsExplode = true
			}

			if checkForCollisions(projectile.X, projectile.Y, exitButton.X, exitButton.Y, exitButton.X+exitButton.Width, exitButton.Y+exitButton.Height) {
				os.Exit(0)
			}
		}
	}
}

func checkForCollisions(x, y float64, x1, y1, x2, y2 float64) bool {

	return x1 < x && x < x2 && y1 < y && y < y2
}

// ******************* 绘制主菜单 ******************
func MenuDraw(screen *ebiten.Image) {
	drawButton(screen)
	drawTip(screen)
	drawLogo(screen)
	drawKeyborad(screen)
}

func drawOneKey(x, y float32, w float32, keyWord string, screen *ebiten.Image) {

	defaultClr := color.RGBA{255, 215, 0, 255}
	pressClr := color.RGBA{255, 128, 0, 255}

	vector.StrokeRect(screen, x, y, w, 25, 1, color.Black, true)
	vector.DrawFilledRect(screen, x+1, y+1, w-2, 25-2, defaultClr, true)

	flag := false
	switch keyWord {
	case "W":
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			flag = true
		}
	case "S":
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			flag = true
		}
	case "A":
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			flag = true
		}
	case "D":
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			flag = true
		}
	case "J":
		if ebiten.IsKeyPressed(ebiten.KeyJ) {
			flag = true
		}
	case "K":
		if ebiten.IsKeyPressed(ebiten.KeyK) {
			flag = true
		}
	case "Space":
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			flag = true
		}
	}
	if flag {
		vector.DrawFilledRect(screen, x+1, y+1, w-2, 25-2, pressClr, true)
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.Black)
	op.GeoM.Translate(float64(x+2.0), float64(y-2.0))

	text.Draw(screen, keyWord, &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   23}, op)

}

func drawKeyborad(screen *ebiten.Image) {

	drawOneKey(100.0, 400.0, 25.0, "W", screen)
	drawOneKey(100.0, 425.0, 25.0, "S", screen)
	drawOneKey(75.0, 425.0, 25.0, "A", screen)
	drawOneKey(125.0, 425.0, 25.0, "D", screen)

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.Black)
	op.GeoM.Translate(float64(75.0), float64(450.0))
	text.Draw(screen, "Move", &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   23}, op)

	drawOneKey(300.0, 425.0, 25.0, "J", screen)
	drawOneKey(325.0, 425.0, 25.0, "K", screen)
	op.GeoM.Reset()
	op.ColorScale.ScaleWithColor(color.Black)
	op.GeoM.Translate(float64(300.0), float64(450.0))
	text.Draw(screen, "Aim", &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   23}, op)

	drawOneKey(175.0, 425.0, 100, "Space", screen)
	op.GeoM.Reset()
	op.ColorScale.ScaleWithColor(color.Black)
	op.GeoM.Translate(float64(200), float64(450.0))
	text.Draw(screen, "Shoot", &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   23}, op)
}

func drawLogo(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(.25, .25)
	op.GeoM.Translate(100, 150)
	screen.DrawImage(logoImg, op)
}
func drawTip(screen *ebiten.Image) {

	op := &text.DrawOptions{}

	op.ColorScale.ScaleWithColor(color.RGBA{128, 138, 135, 255})
	op.GeoM.Translate(100, 50)

	text.Draw(screen, "github.com/gofish2020/tankgame", &text.GoTextFace{
		Source: mplusNormalFont,
		Size:   50}, op)
}
func drawButton(screen *ebiten.Image) {

	// play button
	buttonOp := &ebiten.DrawImageOptions{}

	buttonOp.GeoM.Translate(playButton.X, playButton.Y)
	screen.DrawImage(playButtonImg, buttonOp)

	// exit button

	buttonOp.GeoM.Reset()

	buttonOp.GeoM.Translate(exitButton.X, exitButton.Y)
	screen.DrawImage(exitButtonImg, buttonOp)
}

//////////////////////////////////////////

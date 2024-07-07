package monitor

import "github.com/hajimehoshi/ebiten/v2"

var (
	ScreenWidth  float64
	ScreenHeight float64
)

func init() {
	w, h := ebiten.Monitor().Size()
	ScreenWidth, ScreenHeight = float64(w), float64(h)
}

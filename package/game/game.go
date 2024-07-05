package game

import (
	"math"

	"github.com/gofish2020/tankgame/package/keyboard"
	"github.com/gofish2020/tankgame/package/tank"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ScreenWidth  int
	ScreenHeight int

	hullImage = "resource/green_tank_hull.png"
)

func init() {

	ScreenWidth, ScreenHeight = ebiten.Monitor().Size()

}

type Game struct {
	tks []*tank.Tank
}

func NewGame() *Game {

	game := Game{}
	game.tks = append(game.tks, tank.NewTank(float64(ScreenWidth/2.0), float64(ScreenHeight/2.0), tank.TankTypePlayer))

	game.tks = append(game.tks, tank.NewTank(float64(ScreenWidth/2.0), float64(ScreenHeight/2.0), tank.TankTypeNPC))
	return &game
}

type Points struct {
	x float64
	y float64
	t *tank.Tank
}

func (g *Game) Update() error {

	var playerPoints Points

	var npcPoints []Points
	// 更新每个坦克数据
	for _, tk := range g.tks {
		tk.Update()
		tk.LimitRange(30, 30, float64(ScreenWidth)-30, float64(ScreenHeight)-30)

		if tk.TkType == tank.TankTypePlayer {
			playerPoints.x = tk.X
			playerPoints.y = tk.Y
			playerPoints.t = tk
		} else {
			npcPoints = append(npcPoints, Points{x: tk.X, y: tk.Y, t: tk})
		}

	}

	minDistance := math.MaxInt
	// 计算距离 player 最近的坦克
	playerPoints.t.Enemy = nil
	for _, npcPoint := range npcPoints {
		x := npcPoint.x - playerPoints.x
		y := npcPoint.y - playerPoints.y
		distance := int(math.Sqrt(x*x + y*y))

		if minDistance > distance { // 找距离最近的坦克
			minDistance = distance
			if minDistance <= 300 {
				playerPoints.t.Enemy = npcPoint.t
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 清屏
	screen.Clear()

	// 绘制每个坦克
	for _, tk := range g.tks {
		tk.Draw(screen)
	}

	// 绘制键盘

	keyboard.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

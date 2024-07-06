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

	game.tks = append(game.tks, tank.NewTank(float64(ScreenWidth/2.0+100), float64(ScreenHeight/2.0+100), tank.TankTypeNPC))
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

	// 更新处于npc攻击范围内的坦克
	for _, npcPoint := range npcPoints {

		// 默认（无敌人）坦克
		npcPoint.t.Enemy = nil

		x := playerPoints.x - npcPoint.x
		y := playerPoints.y - npcPoint.y
		distance := math.Sqrt(x*x + y*y)

		if npcPoint.t.Turrent.RangeDistance >= distance { // 在攻击范围内

			// 在视野内
			angle := math.Atan2(y, x) * 180 / math.Pi
			if angle < 0 {
				angle += 360.0
			}

			startAngle, endAngle := npcPoint.t.Turrent.Angle-npcPoint.t.Turrent.RangeAngle, npcPoint.t.Turrent.Angle+npcPoint.t.Turrent.RangeAngle
			if endAngle > 360 {
				endAngle -= 360
			}
			if startAngle < 0 {
				startAngle += 360
			}
			// 正常情况下 startAngle <= endAngle
			if startAngle <= endAngle && startAngle <= angle && angle <= endAngle {
				npcPoint.t.Enemy = playerPoints.t
			} else {
				// 如果处于 0 or 360的分割位置，startAngle > endAngle
				if angle <= endAngle || angle >= startAngle {
					npcPoint.t.Enemy = playerPoints.t
				}
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
	keyboard.Draw(g.tks[0], screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

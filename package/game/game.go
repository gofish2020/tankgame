package game

import (
	"math"

	"github.com/gofish2020/tankgame/package/tank"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ScreenWidth  int
	ScreenHeight int
)

func init() {

	ScreenWidth, ScreenHeight = ebiten.Monitor().Size()

}

type Game struct {
	tks   []*tank.Tank
	incre int16
}

func NewGame() *Game {

	game := Game{}
	game.tks = append(game.tks, tank.NewTank(float64(ScreenWidth/2.0), float64(ScreenHeight/2.0), tank.TankTypePlayer))
	//game.AddEnemy(3)
	return &game
}

// 新增敌人
func (g *Game) AddEnemy(count int) {

	for range count {
		x, y := tank.MinXCoordinates, tank.MinYCoordinates
		switch g.incre % 3 { // 按照轮询的方式，选择放置位置
		case 0:

		case 1:
			x = float64(ScreenWidth) / 2.0
		case 2:
			x = float64(ScreenWidth) - tank.MinXCoordinates
		}
		g.tks = append(g.tks, tank.NewTank(x, y, tank.TankTypeNPC))
		g.incre++
	}

	// game.tks = append(game.tks, tank.NewTank(float64(ScreenWidth/2.0+100), float64(ScreenHeight/2.0+100), tank.TankTypeNPC))
}
func (g *Game) Update() error {

	var playerPosition tank.TankPosition

	var npcPositions []tank.TankPosition
	// 更新每个坦克数据
	for _, tk := range g.tks {
		tk.Update()
		// 限制坦克运动范围
		tk.LimitRange(tank.MinXCoordinates, tank.MinYCoordinates, float64(ScreenWidth)-30, float64(ScreenHeight)-30)

		// 记录下坦克当前的位置
		if tk.TkType == tank.TankTypePlayer {

			playerPosition.X = tk.X
			playerPosition.Y = tk.Y
			playerPosition.TK = tk
		} else {
			npcPositions = append(npcPositions, tank.TankPosition{X: tk.X, Y: tk.Y, TK: tk})
		}

	}

	// 更新npc攻击范围内的坦克(为了做自动攻击)
	for _, npcPosition := range npcPositions {

		// 默认（无敌人）坦克
		npcPosition.TK.Enemy = nil

		x := playerPosition.X - npcPosition.X
		y := playerPosition.Y - npcPosition.Y
		distance := math.Sqrt(x*x + y*y)

		if npcPosition.TK.Turrent.RangeDistance >= distance { // 在攻击范围内

			// 在视野内
			angle := math.Atan2(y, x) * 180 / math.Pi
			if angle < 0 {
				angle += 360.0
			}

			startAngle, endAngle := npcPosition.TK.Turrent.Angle-npcPosition.TK.Turrent.RangeAngle, npcPosition.TK.Turrent.Angle+npcPosition.TK.Turrent.RangeAngle
			if endAngle > 360 {
				endAngle -= 360
			}
			if startAngle < 0 {
				startAngle += 360
			}
			// 正常情况下 startAngle <= endAngle
			if startAngle <= endAngle && startAngle <= angle && angle <= endAngle {
				npcPosition.TK.Enemy = playerPosition.TK
			} else {
				// 如果处于 0 or 360的分割位置，startAngle > endAngle
				if angle <= endAngle || angle >= startAngle {
					npcPosition.TK.Enemy = playerPosition.TK
				}
			}
		}
	}
	return nil
}

var (
	menuType = "init"
)

func (g *Game) Draw(screen *ebiten.Image) {
	// 清屏
	screen.Clear()
	//screen.Fill(color.RGBA{240, 222, 180, 255}) // Desert background

	if menuType == "init" {
		tank.MenuDraw(ScreenWidth, ScreenHeight, screen)
	}

	// 绘制每个坦克
	for _, tk := range g.tks {
		tk.Draw(screen)
		// 绘制按键
		if tk.TkType == tank.TankTypePlayer {
			tank.KeyPressDrawAroundTank(tk, screen)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

package game

import (
	"image/color"
	"math"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/gofish2020/tankgame/package/tank"
	"github.com/gofish2020/tankgame/package/utils"
	"github.com/gofish2020/tankgame/package/utils/sound"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	tks  []*tank.Tank
	incr int16

	barriers []*tank.Barrier
}

const defaultFreq = 50

func NewGame() *Game {

	game := Game{}

	game.tks = append(game.tks, tank.NewTank(float64(monitor.ScreenWidth/2.0), float64(monitor.ScreenHeight-30), tank.TankTypePlayer))
	return &game
}

func (g *Game) initData() {
	utils.GameLevel++
	if utils.GameLevel > 4 {
		utils.GameProgress = "pass" // 通关
		return
	}
	g.barriers = tank.NewMap()
	g.tks = nil
	g.tks = append(g.tks, tank.NewTank(float64(monitor.ScreenWidth/2.0), float64(monitor.ScreenHeight-30), tank.TankTypePlayer))
	g.AddEnemy(2 * utils.GameLevel)
	utils.GameProgress = "play"
}

func (g *Game) Restart() {
	if utils.GameProgress == "prepare" || utils.GameProgress == "next" {

		if utils.GameProgress == "prepare" {
			utils.KilledCount = 0
			utils.GameLevel = 0
		}

		g.initData()
	}
}

// 新增敌人
func (g *Game) AddEnemy(count int) {

	for range count {
		x, y := tank.MinXCoordinates, tank.MinYCoordinates
		switch g.incr % 3 { // 按照轮询的方式，选择放置位置
		case 0:
		case 1:
			x = float64(monitor.ScreenWidth) / 2.0
		case 2:
			x = float64(monitor.ScreenWidth) - tank.MinXCoordinates
		}

		g.tks = append(g.tks, tank.NewTank(x, y, tank.TankTypeNPC))
		g.incr++
	}
}

func (g *Game) Update() error {

	enemyCount := 0
	g.Restart()

	// 播放 bgm
	sound.PlayBGM()

	// 分离 player 和 npc 坦克
	var playerPosition tank.TankPosition
	var npcPositions []tank.TankPosition

	// 检测存活的坦克
	liveTanks := []*tank.Tank{}

	for _, tk := range g.tks {
		// 更新坦克
		tk.Update()
		// 检测子弹碰撞
		tk.CheckCollisions(g.tks, g.barriers)
		// 限制坦克运动范围
		tk.LimitTankRange(tank.MinXCoordinates, tank.MinYCoordinates, float64(monitor.ScreenWidth)-30, float64(monitor.ScreenHeight)-30)

		// 记录下坦克当前的位置
		if tk.TkType == tank.TankTypePlayer {
			playerPosition.X = tk.X
			playerPosition.Y = tk.Y
			playerPosition.TK = tk
			if tk.HealthPoints == 0 {
				tank.UpdateNameList(tk.Name)
				utils.GameProgress = "over"
				sound.PlaySound("yiwai")
				break
			}
		} else {
			// 记录npc的位置
			if tk.HealthPoints == 0 {
				tank.UpdateNameList(tk.Name)
				utils.KilledCount++
				tk.DeathSound()
			} else {
				enemyCount++
				npcPositions = append(npcPositions, tank.TankPosition{X: tk.X, Y: tk.Y, TK: tk})
			}
		}

		if tk.HealthPoints != 0 {
			liveTanks = append(liveTanks, tk)
		}
	}

	g.tks = liveTanks

	// 初始界面
	if utils.GameProgress == "init" || utils.GameProgress == "pass" {
		tank.MenuUpdate(g.tks) //  按钮移动 + 炮弹和按钮碰撞
	} else if utils.GameProgress == "play" { // 游戏界面

		// 更新npc攻击范围内的坦克(为了做自动攻击)
		for _, npcPosition := range npcPositions {

			npcPosition.TK.Enemy = nil // 默认无敌人

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
				if startAngle <= endAngle {
					if startAngle <= angle && angle <= endAngle {
						npcPosition.TK.Enemy = playerPosition.TK
					}
				} else {
					// 如果处于 0 or 360的分割位置，startAngle > endAngle
					if angle <= endAngle || angle >= startAngle {
						npcPosition.TK.Enemy = playerPosition.TK
					}
				}
			}

			// 说明视野内没有敌人，自动旋转炮塔
			if npcPosition.TK.Enemy == nil {
				npcPosition.TK.AddTurrentAngle(2.0)

				// 转向player的方向
				angle := math.Atan2(y, x) * 180 / math.Pi
				if angle < 0 {
					angle += 360.0
				}

				//	npcPosition.TK.Angle 表示 坦克 和 x 轴的夹角
				// angle 表示两个坦克连线 和 x轴的夹角
				if npcPosition.TK.Angle > angle {
					// 目的让t.Turrent.Angle 往夹角小的方向移动，让炮台尽可能快的对准敌人
					if npcPosition.TK.Angle-angle > 180 {
						npcPosition.TK.AddTankAngle(1.0)
					} else {
						npcPosition.TK.AddTankAngle(-1)
					}
				} else if npcPosition.TK.Angle < angle {

					if angle-npcPosition.TK.Angle > 180 {
						npcPosition.TK.AddTankAngle(-1)
					} else {
						npcPosition.TK.AddTankAngle(1)
					}
				}

				// 移动坦克
				npcPosition.TK.X += npcPosition.TK.ForwardSpeed * math.Cos(npcPosition.TK.Angle*math.Pi/180)
				npcPosition.TK.Y += npcPosition.TK.ForwardSpeed * math.Sin(npcPosition.TK.Angle*math.Pi/180)
			}

		}
	}

	if utils.GameProgress == "play" && enemyCount == 0 { // 全部消灭
		utils.GameProgress = "next" // 下一关
	}

	tank.GameOverUpdate()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{240, 222, 180, 215})

	if utils.GameProgress == "init" || utils.GameProgress == "pass" {
		tank.MenuDraw(screen)
	}

	x, y := 0.0, 0.0
	// 绘制每个坦克
	for _, tk := range g.tks {
		tk.Draw(screen)
		// 绘制按键
		if tk.TkType == tank.TankTypePlayer {
			tank.KeyPressDrawAroundTank(tk, screen)
			x, y = tk.X, tk.Y // 以player的视角
		}
	}

	// 绘制战争迷雾
	if utils.GameProgress == "play" {
		tank.DrawWarFog(screen, x, y, g.barriers)
	}

	// 绘制死亡名单
	tank.DrawNameList(screen)

	tank.GameOverDraw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(monitor.ScreenWidth), int(monitor.ScreenHeight)
}

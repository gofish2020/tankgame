package game

import (
	"image/color"

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

func NewGame() *Game {
	game := Game{}
	game.tks = append(game.tks, tank.NewTank(float64(monitor.ScreenWidth/2.0), float64(monitor.ScreenHeight-30), tank.TankTypePlayer))
	return &game
}

func (g *Game) Restart() {
	if utils.GameProgress == "prepare" || utils.GameProgress == "next" {

		if utils.GameProgress == "prepare" {
			utils.KilledCount = 0
			utils.GameLevel = 0
		}

		utils.GameLevel++
		if utils.GameLevel > 4 {
			utils.GameProgress = "pass" // 通关
			return
		}
		// 新地图
		g.barriers = tank.NewMap()
		g.tks = nil
		g.tks = append(g.tks, tank.NewTank(float64(monitor.ScreenWidth/2.0), float64(monitor.ScreenHeight-30), tank.TankTypePlayer))
		g.AddEnemy(2 * utils.GameLevel)
		utils.GameProgress = "play"
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

// 更新数据
func (g *Game) Update() error {

	enemyCount := 0
	// 游戏重启
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
		// 检测碰撞
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

	// 更新 g.tks,剩余的坦克
	g.tks = liveTanks

	// 初始界面
	if utils.GameProgress == "init" || utils.GameProgress == "pass" {
		tank.MenuUpdate(g.tks) //  按钮移动 + 炮弹和按钮碰撞
	} else if utils.GameProgress == "play" {
		// 移动 npc 坦克，并检测攻击范围内敌人
		tank.MoveAndFineEnemyTank(playerPosition, npcPositions)
	}

	if utils.GameProgress == "play" && enemyCount == 0 { // 全部消灭
		utils.GameProgress = "next" // 下一关
	}

	// 游戏结束，检测按键消息
	tank.GameOverUpdate()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{240, 222, 180, 215})

	if utils.GameProgress == "init" || utils.GameProgress == "pass" {
		// 起始界面
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

	// 绘制战争迷雾 + 障碍物
	if utils.GameProgress == "play" {
		tank.DrawWarFog(screen, x, y, g.barriers)
	}

	// 绘制死亡名单
	tank.DrawNameList(screen)

	// 游戏结束界面
	tank.GameOverDraw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(monitor.ScreenWidth), int(monitor.ScreenHeight)
}

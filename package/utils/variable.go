package utils

var (
	GameProgress = "prepare" // init 初始界面  prepare 数据准备中  next 下一关  play 游戏进行中 over 游戏over  pass 通关
	GameLevel    = 0         // 游戏关卡
	KilledCount  = 0
)

type TankLevel struct {
	TankSpeed          float64
	TurrentRotateSpeed float64
}

var (
	TankLevels = []TankLevel{
		{TankSpeed: 0.4, TurrentRotateSpeed: 1.0},
		{TankSpeed: 0.4, TurrentRotateSpeed: 2.0},
		{TankSpeed: 0.4, TurrentRotateSpeed: 3.0},
		{TankSpeed: 0.4, TurrentRotateSpeed: 4.0},
		{TankSpeed: 0.4, TurrentRotateSpeed: 5.0},
		{TankSpeed: 0.4, TurrentRotateSpeed: 6.0},

		{TankSpeed: 1.0, TurrentRotateSpeed: 1.0},
		{TankSpeed: 1.0, TurrentRotateSpeed: 2.0},
		{TankSpeed: 1.0, TurrentRotateSpeed: 3.0},
		{TankSpeed: 1.0, TurrentRotateSpeed: 4.0},
		{TankSpeed: 1.0, TurrentRotateSpeed: 5.0},
		{TankSpeed: 1.0, TurrentRotateSpeed: 6.0},

		{TankSpeed: 1.5, TurrentRotateSpeed: 1.0},
		{TankSpeed: 1.5, TurrentRotateSpeed: 2.0},
		{TankSpeed: 1.5, TurrentRotateSpeed: 3.0},
		{TankSpeed: 1.5, TurrentRotateSpeed: 4.0},
		{TankSpeed: 1.5, TurrentRotateSpeed: 5.0},
		{TankSpeed: 1.5, TurrentRotateSpeed: 6.0},

		{TankSpeed: 2.0, TurrentRotateSpeed: 10.0},
		{TankSpeed: 3.0, TurrentRotateSpeed: 10.0},

		{TankSpeed: 5.0, TurrentRotateSpeed: 10.0},
		{TankSpeed: 8.0, TurrentRotateSpeed: 1.0},
	}
)

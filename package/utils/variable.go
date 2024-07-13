package utils

var (
	GameProgress = "prepare" // init 初始界面  prepare 数据准备中  next 下一关  play 游戏进行中 over 游戏over  pass 通关
	GameLevel    = 0         // 游戏关卡
	KilledCount  = 0
)


type TankLevel struct {
	Speed       float64
	RotateSpeed float64
}

var (
	TankLevels = []TankLevel{
		{Speed: 0.2, RotateSpeed: 1.0},
		{Speed: 0.2, RotateSpeed: 2.0},
		{Speed: 0.2, RotateSpeed: 3.0},
		{Speed: 0.2, RotateSpeed: 4.0},
		{Speed: 0.2, RotateSpeed: 5.0},
		{Speed: 0.2, RotateSpeed: 6.0},

		{Speed: 0.5, RotateSpeed: 1.0},
		{Speed: 0.5, RotateSpeed: 2.0},
		{Speed: 5.0, RotateSpeed: 3.0},
		{Speed: 0.5, RotateSpeed: 4.0},
		{Speed: 0.5, RotateSpeed: 5.0},
		{Speed: 0.5, RotateSpeed: 6.0},

		{Speed: 1.0, RotateSpeed: 1.0},
		{Speed: 1.0, RotateSpeed: 2.0},
		{Speed: 1.0, RotateSpeed: 3.0},
		{Speed: 1.0, RotateSpeed: 4.0},
		{Speed: 1.0, RotateSpeed: 5.0},
		{Speed: 1.0, RotateSpeed: 6.0},

		{Speed: 2.0, RotateSpeed: 10.0},
		{Speed: 3.0, RotateSpeed: 10.0},

		{Speed: 5.0, RotateSpeed: 10.0},
		{Speed: 8.0, RotateSpeed: 1.0},
	}
)

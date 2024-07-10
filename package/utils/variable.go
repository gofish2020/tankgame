package utils

var (
	GameProgress = "init" // init  prepare next play over pass
	GameLevel    = 1      // 游戏关卡
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
		{Speed: 0.5, RotateSpeed: 3.0},
		{Speed: 0.5, RotateSpeed: 4.0},
		{Speed: 1.0, RotateSpeed: 5.0},
		{Speed: 0.5, RotateSpeed: 6.0},

		{Speed: 1.0, RotateSpeed: 1.0},
		{Speed: 1.0, RotateSpeed: 2.0},
		{Speed: 1.0, RotateSpeed: 3.0},
		{Speed: 1.0, RotateSpeed: 4.0},
		{Speed: 1.0, RotateSpeed: 5.0},
		{Speed: 1.0, RotateSpeed: 6.0},

		{Speed: 1.0, RotateSpeed: 10.0},
		{Speed: 2.0, RotateSpeed: 10.0},
		{Speed: 3.0, RotateSpeed: 10.0},
		{Speed: 8.0, RotateSpeed: 1.0},
	}
)

package tank

import (
	"math"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/gofish2020/tankgame/package/utils"
)

type line struct {
	X1, Y1, X2, Y2 float64
}

func (l *line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

type object struct {
	Walls []line
}

func (o object) points() [][2]float64 {
	// Get one of the endpoints for all segments,
	// + the startpoint of the first one, for non-closed paths
	var points [][2]float64
	for _, wall := range o.Walls {
		points = append(points, [2]float64{wall.X2, wall.Y2})
	}
	p := [2]float64{o.Walls[0].X1, o.Walls[0].Y1}
	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, [2]float64{o.Walls[0].X1, o.Walls[0].Y1})
	}
	return points
}

// 矩形的四个边
func rect(x, y, w, h float64) []line {

	// 逆时针
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

type Image struct {
	X      int
	Y      int
	Width  int
	Height int
	Path   string
}

type Barrier struct {

	// 障碍物范围（位置 + 大小）
	X      float64
	Y      float64
	Width  float64
	Height float64

	// 障碍物的图片
	Image Image

	// 用来做 阴影计算
	Objects []object

	Border bool

	Destructible bool // 是否可破坏
	Collidable   bool // 是否可碰撞

	Health int
}

func addBoard(x, y float64) *Barrier {

	return &Barrier{
		X:      x,
		Y:      y,
		Width:  monitor.ScreenWidth,
		Height: monitor.ScreenHeight,
		Health: 100,
		Border: true,

		Objects: []object{{rect(0, 0, monitor.ScreenWidth, monitor.ScreenHeight)}},
	}
}

func addBarrier(x, y float64, barrierType string) *Barrier {
	b := Barrier{
		X:      x,
		Y:      y,
		Width:  64,
		Height: 64,
		Health: 100,

		// 默认砖块（为了做裁剪）
		Image: Image{
			X:      0,
			Y:      0,
			Height: 64,
			Width:  64,
		},

		// 为了阴影计算
		Objects: []object{{rect(x, y, 64, 64)}},
	}

	if barrierType == "b" { // 砖块
		b.Image.Path = "resource/Brick_Block_small.png"
		// 不可穿越，可毁坏
		b.Destructible = true
		b.Collidable = true

	} else if barrierType == "c" { // 沙漠
		b.Image.Path = "resource/camo_net.png"
		b.Objects = []object{{Walls: rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算

		// 可穿越，不可毁坏
		b.Destructible = false
		b.Collidable = false

	} else if barrierType == "w" { // 水
		b.Image.Path = "resource/water.png"
		// 说明不用做阴影计算
		b.Objects = []object{{Walls: rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算

		// 不可穿越，不可毁坏
		b.Destructible = false
		b.Collidable = true
	} else if barrierType == "l" { // 草地
		b.Image.Path = "resource/leaves.png"

		// 可穿越 不可毁坏
		b.Destructible = false
		b.Collidable = false
	} else if barrierType == "i" { // 铁块

		b.Health = math.MaxInt // 无敌
		b.Image.Path = "resource/iron.png"
		// 不可穿越 不可毁坏
		b.Destructible = true
		b.Collidable = true
	}
	return &b
}

// 创建地图
func NewMap() []*Barrier {

	freq := 8
	switch utils.GameLevel {
	case 1:
		freq = 25
	case 2:
		freq = 20
	case 3:
		freq = 15
	case 4:
		freq = 10
	}

	return createMap(100, 100, monitor.ScreenWidth-100, monitor.ScreenHeight-100, freq)
}

// x1, y1, x2, y2 绘制障碍物的范围 freq 障碍物出现的可能性
func createMap(x1, y1, x2, y2 float64, freq int) []*Barrier {
	var barriers []*Barrier
	// 初始边界
	barriers = append(barriers, addBoard(0, 0))

	for x := x1; x < x2; x = x + 64. {
		for y := y1; y < y2; y = y + 64. {

			s := ""
			switch r.Intn(freq) {
			case 0:
				s = "b"
			case 1:
				s = "c"
			case 2:
				s = "w"
			case 3:
				s = "l"
			case 4:
				s = "i"
			default:
				s = ""
			}

			if s != "" {
				barriers = append(barriers, addBarrier(x, y, s))
			}
		}
	}

	return barriers
}

func changeBarrier(barrier *Barrier, side string) {

	switch side {
	case "l":

		// 障碍物范围
		barrier.X += 32
		barrier.Width -= 32

		// 图片裁剪
		barrier.Image.X += 32

		// 墙体的边界 （逆时针的四条边）
		barrier.Objects[0].Walls[0].X1 += 32
		barrier.Objects[0].Walls[0].X2 += 32
		barrier.Objects[0].Walls[1].X1 += 32
		barrier.Objects[0].Walls[3].X2 += 32
	case "r":
		barrier.Width -= 32

		barrier.Image.X = 0
		barrier.Image.Width -= 32

		barrier.Objects[0].Walls[1].X2 -= 32
		barrier.Objects[0].Walls[2].X1 -= 32
		barrier.Objects[0].Walls[2].X2 -= 32
		barrier.Objects[0].Walls[3].X1 -= 32
	case "t":
		barrier.Y += 32
		barrier.Height -= 32

		barrier.Image.Y += 32

		barrier.Objects[0].Walls[0].Y1 += 32
		barrier.Objects[0].Walls[2].Y2 += 32
		barrier.Objects[0].Walls[3].Y1 += 32
		barrier.Objects[0].Walls[3].Y2 += 32
	case "b":
		barrier.Height -= 32

		barrier.Image.Height -= 32

		barrier.Objects[0].Walls[0].Y2 -= 32
		barrier.Objects[0].Walls[1].Y1 -= 32
		barrier.Objects[0].Walls[1].Y2 -= 32
		barrier.Objects[0].Walls[2].Y1 -= 32
	}

	// 当障碍物的宽/高为0，表示障碍物已经清理
	if barrier.Height == 0 || barrier.Width == 0 {
		barrier.Health = 0
	}
}

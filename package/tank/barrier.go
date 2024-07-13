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

	var points [][2]float64
	for _, wall := range o.Walls { // 每个线的 X2,Y2
		points = append(points, [2]float64{wall.X2, wall.Y2})
	}

	// 最后闭合最后一点
	p := [2]float64{o.Walls[0].X1, o.Walls[0].Y1}
	// 避免重复
	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, p)
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

type BarrierType string

var (
	BarrierTypeNone   BarrierType = "None"
	BarrierTypeBrick  BarrierType = "Brick"
	BarrierTypeWater  BarrierType = "Water"
	BarrierTypeIron   BarrierType = "Iron"
	BarrierTypeCamo   BarrierType = "Camo"
	BarrierTypeLeaves BarrierType = "Leaves"
	BarrierTypeBug    BarrierType = "Bug"
)

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

	Destructible bool // 是否可破坏
	Collidable   bool // 是否可碰撞

	Health int

	BarrierTypeVal BarrierType
}

func (b Barrier) getBarrierCollisionVectors() []Point {
	vectors := []Point{
		{b.X, b.Y},
		{b.X + b.Width, b.Y},
		{b.X + b.Width, b.Y + b.Height},
		{b.X, b.Y + b.Height},
	}
	return vectors
}

func addBoard(x, y float64) *Barrier {

	return &Barrier{
		X:              x,
		Y:              y,
		Width:          monitor.ScreenWidth,
		Height:         monitor.ScreenHeight,
		Health:         100,
		BarrierTypeVal: BarrierTypeNone,

		Objects: []object{{rect(0, 0, monitor.ScreenWidth, monitor.ScreenHeight)}},
	}
}

func addBarrier(x, y float64, BarrierTypeVal BarrierType) *Barrier {
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

		BarrierTypeVal: BarrierTypeVal,
	}

	if BarrierTypeVal == BarrierTypeBrick { // 砖块
		b.Image.Path = "resource/Brick_Block_small.png"
		// 不可穿越，可毁坏
		b.Destructible = true
		b.Collidable = true

	} else if BarrierTypeVal == BarrierTypeCamo { // 沙漠
		b.Image.Path = "resource/camo_net.png"
		b.Objects = []object{{Walls: rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算

		// 可穿越，不可毁坏
		b.Destructible = false
		b.Collidable = false

	} else if BarrierTypeVal == BarrierTypeWater { // 水
		b.Image.Path = "resource/water.png"
		// 说明不用做阴影计算
		b.Objects = []object{{Walls: rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算

		// 不可穿越，不可毁坏
		b.Destructible = false
		b.Collidable = true
	} else if BarrierTypeVal == BarrierTypeLeaves { // 草地
		b.Image.Path = "resource/leaves.png"

		// 可穿越 不可毁坏
		b.Destructible = false
		b.Collidable = false
	} else if BarrierTypeVal == BarrierTypeIron { // 铁块
		b.Image.Path = "resource/iron.png"
		// 不可穿越 可毁坏
		b.Destructible = true
		b.Collidable = true
	} else if BarrierTypeVal == BarrierTypeBug {

		b.Image.Path = "resource/bug.png"
		// 不可穿越 可毁坏
		b.Destructible = true
		b.Collidable = true
		b.Objects = []object{{Walls: rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算
	}

	return &b
}

const padding = 150

// 创建地图
func NewMap() []*Barrier {

	freq := 8
	switch utils.GameLevel {
	case 1:
		freq = 30
	case 2:
		freq = 24
	case 3:
		freq = 18
	case 4:
		freq = 12
	}

	return createMap(padding, padding, monitor.ScreenWidth-padding, monitor.ScreenHeight-padding, freq)
}

// x1, y1, x2, y2 绘制障碍物的范围 freq 障碍物出现的可能性
func createMap(x1, y1, x2, y2 float64, freq int) []*Barrier {

	bugCount := 1
	var barriers []*Barrier
	// 初始边界
	barriers = append(barriers, addBoard(0, 0))

	for x := x1; x < x2; x = x + 64. {
		for y := y1; y < y2; y = y + 64. {

			typeVal := BarrierTypeNone
			switch r.Intn(freq) {
			case 0:
				typeVal = BarrierTypeBrick
			case 1:
				typeVal = BarrierTypeCamo
			case 2:
				typeVal = BarrierTypeWater
			case 3:
				typeVal = BarrierTypeLeaves
			case 4:
				typeVal = BarrierTypeIron
			case 5:
				typeVal = BarrierTypeBug
			}

			if typeVal != BarrierTypeNone {
				if typeVal == BarrierTypeBug {
					if bugCount > 0 {
						bugCount--
					} else {
						typeVal = BarrierTypeBrick
					}
				}
				barriers = append(barriers, addBarrier(x, y, typeVal))
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

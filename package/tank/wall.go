package tank

import (
	"math"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/gofish2020/tankgame/package/utils"
)

type Line struct {
	X1, Y1, X2, Y2 float64
}

func (l *Line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

type Object struct {
	Walls []Line
}

func (o Object) points() [][2]float64 {
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
func Rect(x, y, w, h float64) []Line {
	return []Line{
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

	// 障碍物的位置
	X      float64
	Y      float64
	Width  float64
	Height float64
	// 障碍物的图片
	Image Image

	// 用来做 阴影计算
	Objects []Object

	Border bool

	Destructible bool // 是否可破坏
	Collidable   bool // 是否可碰撞
}

func AddBoard(x, y float64) Barrier {

	return Barrier{
		X:      x,
		Y:      y,
		Width:  monitor.ScreenWidth,
		Height: monitor.ScreenHeight,
		Border: true,

		Objects: []Object{{Rect(0, 0, monitor.ScreenWidth, monitor.ScreenHeight)}},
	}
}

func AddBarrier(x, y float64, barrierType string) Barrier {
	b := Barrier{
		X:      x,
		Y:      y,
		Width:  64,
		Height: 64,
		Image: Image{ // 默认砖块

			X:      0,
			Y:      0,
			Height: 64,
			Width:  64,
		},

		// 作为阴影计算
		Objects: []Object{{Rect(x, y, 64, 64)}},
	}

	if barrierType == "b" { // 砖块

		b.Image.Path = "resource/Brick_Block_small.png"

		// 不可穿越，可毁坏
		b.Destructible = true
		b.Collidable = true

	} else if barrierType == "c" { // 沙漠
		b.Image.Path = "resource/camo_net.png"
		b.Objects = []Object{{Walls: Rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算

		// 可穿越，不可毁坏
		b.Destructible = false
		b.Collidable = false

	} else if barrierType == "w" { // 水
		b.Image.Path = "resource/water.png"
		// 说明不用做阴影计算
		b.Objects = []Object{{Walls: Rect(float64(x), float64(y), 0.0, 0.0)}} // 说明不用做阴影计算

		// 不可穿越，不可毁坏
		b.Destructible = false
		b.Collidable = true
	} else if barrierType == "l" { // 草地
		b.Image.Path = "resource/leaves.png"

		// 可穿越 不可毁坏
		b.Destructible = false
		b.Collidable = false
	}
	return b
}

func NewMap() []Barrier {

	freq := 8
	switch utils.GameLevel {
	case 1:
		freq = 64
	case 2:
		freq = 32
	case 3:
		freq = 16
	case 4:
		freq = 8
	}

	return createMap(100, 100, monitor.ScreenWidth-100, monitor.ScreenHeight-100, freq)
}

func createMap(x1, y1, x2, y2 float64, freq int) []Barrier {
	var barriers []Barrier
	// 初始边界
	barriers = append(barriers, AddBoard(0, 0))

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
			default:
				s = ""
			}

			if s != "" {
				barriers = append(barriers, AddBarrier(x, y, s))
			}
		}
	}

	return barriers
}

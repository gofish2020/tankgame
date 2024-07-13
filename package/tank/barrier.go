package tank

import (
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/gofish2020/tankgame/package/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	// 顺时针的四个顶点
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

// 创建地图

const padding = 150

func NewMap() []*Barrier {

	freq := 0
	switch utils.GameLevel {
	case 1:
		freq = 30
	case 2:
		freq = 24
	case 3:
		freq = 18
	case 4:
		freq = 12
	default:
		freq = 10
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

////////////////////////// 光源照射 （阴影计算）////////////////////////

var (
	// 阴影
	shadowImage   = ebiten.NewImage(int(monitor.ScreenWidth), int(monitor.ScreenHeight))
	triangleImage = ebiten.NewImage(int(monitor.ScreenWidth), int(monitor.ScreenHeight))
)

func init() {
	triangleImage.Fill(color.White)
}

func DrawWarFogAndBarriers(screen *ebiten.Image, x, y float64, barriers []*Barrier) {

	if !utils.FullMap {
		drawFog(screen, x, y, barriers)
	}
	// 绘制障碍物
	drawBarrier(screen, x, y, barriers)

}

func drawBarrier(screen *ebiten.Image, x, y float64, barriers []*Barrier) {
	// 绘制障碍物
	for _, barrier := range barriers {
		if barrier.BarrierTypeVal == BarrierTypeNone || barrier.Health == 0 {
			continue
		}
		originalImg, _, _ := ebitenutil.NewImageFromFile(barrier.Image.Path)
		// 对图片 originalImg 进行裁剪
		subImg := originalImg.SubImage(image.Rect(barrier.Image.X, barrier.Image.Y,
			barrier.Image.Width, barrier.Image.Height)).(*ebiten.Image)
		// 绘制裁剪后的图片
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(barrier.X, barrier.Y)
		screen.DrawImage(subImg, options)
	}
}
func drawFog(screen *ebiten.Image, x, y float64, barriers []*Barrier) {
	shadowImage.Fill(color.Black)

	// x,y 相当于光源的位置
	rays := rayCasting(float64(x), float64(y), barriers)

	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.Blend = ebiten.BlendSourceOut
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]
		// 用三个点构成一个三角形
		v := rayVertices(float64(x), float64(y), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
		// 裁剪为白色
		shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, triangleImage, opt)
	}

	// 绘制迷雾最终效果
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(1.0)
	screen.DrawImage(shadowImage, op)
}

// intersection 计算给定的两条之间的交点
func intersection(l1, l2 line) (float64, float64, bool) {

	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
	uNum := -((l1.X1-l1.X2)*(l1.Y1-l2.Y1) - (l1.Y1-l1.Y2)*(l1.X1-l2.X1))

	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom
	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}

	x := l1.X1 + t*(l1.X2-l1.X1)
	y := l1.Y1 + t*(l1.Y2-l1.Y1)
	return x, y, true
}

func newRay(x, y, length, angle float64) line {
	return line{
		X1: x,
		Y1: y,
		X2: x + length*math.Cos(angle),
		Y2: y + length*math.Sin(angle),
	}
}

// rayCasting 返回从点 cx, cy 出发并与对象相交的直线切片
func rayCasting(cx, cy float64, barriers []*Barrier) []line {
	const rayLength = 10000 // something large enough to reach all objects

	var rays []line

	for _, bar := range barriers {

		if bar.Health > 0 { // 障碍物有血

			for _, obj := range bar.Objects {
				// 遍历每个对象中【点集合】
				for _, p := range obj.points() {
					// cx/cy 和 p[0],p[1] 构成一个线段
					l := line{cx, cy, p[0], p[1]}
					// 从 cx/cy 出发到 p[0]/p[1] 构成的线段和 x轴正方向的夹角
					angle := l.angle()

					for _, offset := range []float64{-0.005, 0.005} {
						points := [][2]float64{}

						// 从点 cx,cy 发出一束光，长度为rayLength，角度为 angle +/- 0.005
						ray := newRay(cx, cy, rayLength, angle+offset)

						// 将光线ray 和 所有对象的所有的边，求交点
						for _, bar := range barriers { // 所有的对象

							if bar.Health > 0 { // 障碍物有血

								for _, o := range bar.Objects {
									for _, wall := range o.Walls {
										if px, py, ok := intersection(ray, wall); ok { // 判断两个线段是否有交点
											points = append(points, [2]float64{px, py}) // 记录交点
										}
									}
								}
							}
						}

						// 只保留 和 cx/cy 距离最近的交点
						min := math.Inf(1) // 正无穷
						minI := -1
						for i, p := range points {
							d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1]) // 点 cx/cy 和 p[0]/p[1] 之间的距离的平方（勾股定理）
							if d2 < min {
								min = d2
								minI = i
							}
						}

						if minI != -1 {
							// 记录距离 cx/cy 和 最近的点，组成的线段
							rays = append(rays, line{cx, cy, points[minI][0], points[minI][1]})
						}
					}
				}
			}
		}

	}

	// Sort rays based on angle, otherwise light triangles will not come out right
	sort.Slice(rays, func(i int, j int) bool {
		return rays[i].angle() < rays[j].angle()
	})
	return rays
}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

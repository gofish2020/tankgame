package tank

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

// 坦克 + 子弹碰撞检测

func (t *Tank) CheckCollisions(tks []*Tank) {

	for _, projectile := range t.Projectiles {
		for _, tk := range tks {

			if !projectile.IsExplode {
				if isProjectileCollisionsTank(projectile.X, projectile.Y, tk, t) {
					projectile.IsExplode = true
				}
			}

		}
	}
}

func isProjectileCollisionsTank(x, y float64, t *Tank, origin *Tank) bool {

	if t == origin {
		return false
	}

	vertices := []Point{{t.CollisionX1, t.CollisionY1}, {t.CollisionX2, t.CollisionY2}, {t.CollisionX3, t.CollisionY3}, {t.CollisionX4, t.CollisionY4}}
	if checkCollision2(Point{x, y}, vertices) {
		fmt.Println(111)
		return true
	}
	return false
}

func checkCollision(pX, pY, x1, y1, x2, y2, x3, y3, x4, y4, tankAngle float64) bool {

	angleRad := tankAngle * math.Pi / 180.0

	rotatedPX := math.Cos(angleRad)*(pX-x1) - math.Sin(angleRad)*(pY-y1) + x1
	rotatedPY := math.Sin(angleRad)*(pX-x1) + math.Cos(angleRad)*(pY-y1) + y1

	// Calculate vectors from point 1 to the other corners of the rectangle
	vector1X := x2 - x1
	vector1Y := y2 - y1
	vector2X := x3 - x1
	vector2Y := y3 - y1

	// Calculate vectors from point 1 to the rotated projectile point
	vectorPX := rotatedPX - x1
	vectorPY := rotatedPY - y1

	// Calculate dot products
	dot1 := vectorPX*vector1X + vectorPY*vector1Y
	dot2 := vectorPX*vector2X + vectorPY*vector2Y

	// Check if the point is inside the rectangle
	return dot1 >= 0 && dot1 <= vector1X*vector1X+vector1Y*vector1Y &&
		dot2 >= 0 && dot2 <= vector2X*vector2X+vector2Y*vector2Y
}

func checkCollision1(pX, pY, x1, y1, x2, y2, x3, y3, x4, y4, cx, cy, tankAngle float64) bool {

	rotatedPX, rotatedPY := rotatePoint(pX, pY, -tankAngle, cx, cy)

	// 计算从中心点到其他顶点的向量
	vector1X := x2 - x1
	vector1Y := y2 - y1
	vector2X := x3 - x1
	vector2Y := y3 - y1

	// 计算从中心点到旋转后点的向量
	vectorPX := rotatedPX - x1
	vectorPY := rotatedPY - y1

	// 计算点积
	dot1 := vectorPX*vector1X + vectorPY*vector1Y
	dot2 := vectorPX*vector2X + vectorPY*vector2Y

	// 判断点是否在矩形内
	return dot1 >= 0 && dot1 <= vector1X*vector1X+vector1Y*vector1Y &&
		dot2 >= 0 && dot2 <= vector2X*vector2X+vector2Y*vector2Y
}

func checkCollision2(point Point, vertices []Point) bool {

	// 使用交叉乘积法判断点是否在多边形内
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		if !isLeft(vertices[i], vertices[next], point) {
			return false
		}
	}
	return true
}

// 判断点是否在边的左侧
func isLeft(p1, p2, p Point) bool {
	return ((p2.X-p1.X)*(p.Y-p1.Y) - (p2.Y-p1.Y)*(p.X-p1.X)) > 0
}

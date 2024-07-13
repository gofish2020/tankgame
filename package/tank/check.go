package tank

import (
	"math"
	"math/rand"
	"time"
)

type Point struct {
	X, Y float64
}

func (t *Tank) CheckCollisions(tks []*Tank, barriers []*Barrier) {

	//子弹碰撞检测 (坦克+障碍物)
	t.checkProjectileCollisionWithTankOrBarriers(tks, barriers)

	// 坦克和障碍物碰撞检测
	if t.hasActorCollided(barriers) {
		t.moveActorToPreviousPosition()

		if t.TkType == TankTypeNPC { // 避免npc tank 卡住

			// Check if enough time has passed since the last collision
			if time.Since(t.LastCollisionTime) > time.Second {
				// Randomly turn left or right
				if rand.Intn(2) == 0 {
					t.Angle += 90.0
				} else {
					t.Angle -= 90.0
				}

				// Update the last collision time
				t.LastCollisionTime = time.Now()
			}
		}
	}
}

func (t *Tank) moveActorToPreviousPosition() {
	t.X = t.PreX
	t.Y = t.PreY
}

func (t *Tank) hasActorCollided(barriers []*Barrier) bool {
	// 获取坦克的四个顶点
	tankVectors := t.getTankCollisionVectors()

	for _, barrier := range barriers {
		// 前提障碍物时可以被碰撞的（并且血量不为0）
		if !barrier.Border && barrier.Collidable && barrier.Health > 0 {

			// 获取障碍物的四个顶点
			objectVectors := barrier.getBarrierCollisionVectors()

			// 检测两个矩形是否相交
			if vectorsIntersect(tankVectors, objectVectors) {
				return true
			}
		}
	}

	// No collision detected
	return false
}

func (t Tank) getTankCollisionVectors() []Point {
	// Define tank's collision points as vectors
	vectors := []Point{
		{t.CollisionX1, t.CollisionY1},
		{t.CollisionX2, t.CollisionY2},
		{t.CollisionX3, t.CollisionY3},
		{t.CollisionX4, t.CollisionY4},
	}
	return vectors
}

func vectorsIntersect(vectors1, vectors2 []Point) bool {
	// Check for intersections between two sets of vectors

	// Check for intersections on each axis
	for _, axis := range getAxes(vectors1) {
		if !projectionOverlap(axis, vectors1, vectors2) {
			return false
		}
	}

	for _, axis := range getAxes(vectors2) {
		if !projectionOverlap(axis, vectors1, vectors2) {
			return false
		}
	}

	return true
}

// Project vectors onto an axis and check for overlap
func projectionOverlap(axis Point, vectors1, vectors2 []Point) bool {
	min1, max1 := projectOntoAxis(axis, vectors1)
	min2, max2 := projectOntoAxis(axis, vectors2)

	// Check for overlap on the axis
	return (min1 <= max2 && max1 >= min2) || (min2 <= max1 && max2 >= min1)
}

// Project vectors onto an axis and return the min and max values
func projectOntoAxis(axis Point, vectors []Point) (float64, float64) {
	min, max := dotProduct(axis, vectors[0]), dotProduct(axis, vectors[0])

	for _, point := range vectors[1:] {
		projection := dotProduct(axis, point)
		if projection < min {
			min = projection
		}
		if projection > max {
			max = projection
		}
	}

	return min, max
}

// 向量点积
func dotProduct(v1, v2 Point) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func getAxes(rectVectors []Point) []Point {
	axes := make([]Point, len(rectVectors))

	for i, point := range rectVectors {
		nextPoint := rectVectors[(i+1)%len(rectVectors)]
		edgeVector := Point{X: nextPoint.X - point.X, Y: nextPoint.Y - point.Y}
		// Get the perpendicular vector (normal) to the edge
		axes[i] = Point{X: -edgeVector.Y, Y: edgeVector.X}
	}
	return axes
}

// ////////////////// 子弹 和 （坦克+障碍物）碰撞检测 //////////////////
func (t *Tank) checkProjectileCollisionWithTankOrBarriers(tks []*Tank, barriers []*Barrier) {
	for _, projectile := range t.Projectiles {
		for _, tk := range tks {

			if !projectile.IsExplode { // 子弹正常情况下
				if isProjectileCollisionsTank(projectile.X, projectile.Y, tk, t) {
					projectile.IsExplode = true
				}
			}
		}

		if projectile.IsExplode { // 已碰撞
			return
		}

		for _, barrier := range barriers {
			if !projectile.IsExplode { // 子弹正常情况下
				if isProjectileCollisionsBarrier(projectile.X, projectile.Y, barrier) {
					projectile.IsExplode = true
				}
			}
		}
	}
}

func isProjectileCollisionsBarrier(x, y float64, barrier *Barrier) bool {

	if !barrier.Border && barrier.Destructible && barrier.Health > 0 {

		// 障碍物的范围
		left := barrier.X
		right := (barrier.X + barrier.Width)
		top := barrier.Y
		bottom := (barrier.Y + barrier.Height)

		// 在障碍物内
		if x >= left && x <= right && y >= top && y <= bottom {

			if barrier.Health != math.MaxInt { // 表示砖块
				//裁剪
				dx := math.Min(right-x, x-left) // 距离左右的最短距离
				dy := math.Min(bottom-y, y-top) // 距离上下的最短距离

				if dx < dy { // 说明距离左右的最小值 比距离上下的最小值【更小】
					if x-left < right-x { // 更靠近left
						changeBarrier(barrier, "l")
					} else { // 更靠近 right
						changeBarrier(barrier, "r")
					}
				} else {
					if y-top < bottom-y { // top
						changeBarrier(barrier, "t")
					} else { // bottom
						changeBarrier(barrier, "b")
					}
				}
			}

			return true
		}
	}
	return false
}

func isProjectileCollisionsTank(x, y float64, t *Tank, origin *Tank) bool {

	if t == origin {
		return false
	}

	// vertices := []Point{{t.CollisionX1, t.CollisionY1}, {t.CollisionX2, t.CollisionY2}, {t.CollisionX3, t.CollisionY3}, {t.CollisionX4, t.CollisionY4}}
	// if checkCollision1(Point{x, y}, vertices) {
	// 	t.HealthPoints -= 50 // 扣除血条
	// 	return true
	// }

	if checkCollision(Point{x, y}, t.X, t.Y, t.Width, t.Height, t.Angle) {
		t.HealthPoints -= 50 // 扣除血条
		return true
	}
	return false
}

func checkCollision(point Point, cx, cy float64, width, height float64, tankAngle float64) bool {

	// 坦克旋转 tankAngle角度，等价于 坦克不旋转，点 point 逆向旋转 -tankAngle
	rotatedPX, rotatedPY := rotatePoint(point.X, point.Y, -tankAngle, cx, cy)

	halfW, halfH := width/2, height/2

	xTop, yTop := cx-halfW, cy-halfH
	xBottom, yBottom := cx+halfW, cy+halfH

	// 就把旋转矩形变成不旋转的状态
	if xTop <= rotatedPX && rotatedPX <= xBottom && yTop <= rotatedPY && rotatedPY <= yBottom {
		return true
	}
	return false
}

func checkCollision1(point Point, vertices []Point) bool {

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

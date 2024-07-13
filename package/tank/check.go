package tank

import (
	"math"
	"math/rand"
	"time"

	"github.com/gofish2020/tankgame/package/utils"
)

type Point struct {
	X, Y float64
}

func (t *Tank) CheckCollisions(tks []*Tank, barriers []*Barrier) {

	//子弹碰撞检测 (坦克+障碍物)
	t.checkProjectileCollisionWithTankOrBarriers(tks, barriers)

	// 坦克和障碍物碰撞检测
	if t.hasTankCollided(barriers) {
		t.moveTankToPreviousPosition()

		if t.TkType == TankTypeNPC { // 避免npc tank 卡住

			// Check if enough time has passed since the last collision
			if time.Since(t.LastCollisionTime) > 1*time.Second {
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

func (t *Tank) moveTankToPreviousPosition() {
	t.X = t.PreX
	t.Y = t.PreY
}

func (t *Tank) hasTankCollided(barriers []*Barrier) bool {
	// 获取坦克的四个顶点
	tankVectors := t.getTankCollisionVectors()

	for _, barrier := range barriers {
		// 前提障碍物时可以被碰撞的（并且血量不为0）
		if barrier.BarrierTypeVal != BarrierTypeNone && barrier.Collidable && barrier.Health > 0 {

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

// 获取坦克的四个顶点
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

// 检测矩形相交
func vectorsIntersect(vectors1, vectors2 []Point) bool {

	// 每条边的垂直法向量
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

// 将矩形的四个顶点向量 投影到轴上并检查重叠
func projectionOverlap(axis Point, vectors1, vectors2 []Point) bool {
	// 矩形四个顶点，在 axis投影的范围
	min1, max1 := projectOntoAxis(axis, vectors1)
	// 矩形四个顶点，在 axis投影的范围
	min2, max2 := projectOntoAxis(axis, vectors2)

	// 投影范围有重叠
	return (min1 <= max2 && max1 >= min2) || (min2 <= max1 && max2 >= min1)
}

// 计算多边形在给定轴上的投影，并返回投影的最小值和最大值。
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

// 向量点积(结果是一个标量)表示v2在v1上的投影 |v2|Cos(θ) * |v1|的长度
func dotProduct(v1, v2 Point) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

// 这个 getAxes 函数计算的是多边形的法向量
func getAxes(rectVectors []Point) []Point {
	axes := make([]Point, len(rectVectors))

	for i, point := range rectVectors {
		nextPoint := rectVectors[(i+1)%len(rectVectors)]
		edgeVector := Point{X: nextPoint.X - point.X, Y: nextPoint.Y - point.Y}
		// 获取边的垂直向量（法向量）法向量是通过交换分量并改变一个分量的符号得到的垂直向量。
		axes[i] = Point{X: -edgeVector.Y, Y: edgeVector.X}
	}
	return axes
}

// ////////////////// 子弹 和 （坦克+障碍物）碰撞检测 //////////////////
func (t *Tank) checkProjectileCollisionWithTankOrBarriers(tks []*Tank, barriers []*Barrier) {
	for _, projectile := range t.Projectiles {
		for _, tk := range tks {

			if !projectile.IsExplode { // 坦克
				if isProjectileCollisionsTank(projectile.X, projectile.Y, tk, t) {
					projectile.IsExplode = true
					return
				}
			}
		}

		for _, barrier := range barriers {
			if !projectile.IsExplode { // 障碍物
				if isProjectileCollisionsBarrier(projectile.X, projectile.Y, barrier) {
					projectile.IsExplode = true
					return
				}
			}
		}
	}
}

func isProjectileCollisionsBarrier(x, y float64, barrier *Barrier) bool {

	if barrier.BarrierTypeVal != BarrierTypeNone && barrier.Destructible && barrier.Health > 0 {

		// 障碍物的范围
		left := barrier.X
		right := (barrier.X + barrier.Width)
		top := barrier.Y
		bottom := (barrier.Y + barrier.Height)

		// 在障碍物内
		if x >= left && x <= right && y >= top && y <= bottom {

			if barrier.BarrierTypeVal == BarrierTypeBrick { // 表示砖块

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
			} else if barrier.BarrierTypeVal == BarrierTypeBug {
				barrier.Health = 0
				utils.FullMap = !utils.FullMap
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

	// 利用差积判断是否在矩形内部
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

// 这里的 vertices 的坐标点，是逆时针的四个点
func checkCollision1(point Point, vertices []Point) bool {

	// 使用交叉乘积法判断点是否在多边形内
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		// 因为四条边的方向是逆时针，来进行向量计算，所以 point 也需要在处于逆时针的方向
		if isOut(vertices[i], vertices[next], point) {
			return false
		}
	}
	// 如果point 都在四条边的左边，说明point在 vertices 内部
	return true
}

// 叉积公式 u * v =  (x1,y1) * (x2,y2) = x1*y2 - y1*x2 如果大于0，表示 v 在u的右侧（顺时针方向）
func isOut(p1, p2, p Point) bool {

	// 从 p1 到p2 的向量 * 从 p1 到 p的向量
	return ((p2.X-p1.X)*(p.Y-p1.Y) - (p2.Y-p1.Y)*(p.X-p1.X)) > 0
}

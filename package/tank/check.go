package tank

import "math"

type Point struct {
	X, Y float64
}

// 坦克 + 子弹碰撞检测 / 障碍物 + 子弹 碰撞检测

func (t *Tank) CheckCollisions(tks []*Tank, barriers []*Barrier) {

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

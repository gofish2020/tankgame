package tank

type Point struct {
	X, Y float64
}

// 坦克 + 子弹碰撞检测

func (t *Tank) CheckCollisions(tks []*Tank) {

	for _, projectile := range t.Projectiles {
		for _, tk := range tks {

			if !projectile.IsExplode { // 子弹正常情况下
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

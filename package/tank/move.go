package tank

import "math"

// 坦克移动
type TankPosition struct {
	X  float64
	Y  float64
	TK *Tank
}

// 让npc坦克，检测在攻击范围内的 player 并且朝 player移动
func MoveAndFineEnemyTank(playerPosition TankPosition, npcPositions []TankPosition) {

	for _, npcPosition := range npcPositions {

		npcPosition.TK.Enemy = nil // 默认无敌人

		x := playerPosition.X - npcPosition.X
		y := playerPosition.Y - npcPosition.Y
		distance := math.Sqrt(x*x + y*y)

		if npcPosition.TK.Turrent.RangeDistance >= distance { // 在攻击范围内

			// 在视野内
			angle := math.Atan2(y, x) * 180 / math.Pi
			if angle < 0 {
				angle += 360.0
			}
			startAngle, endAngle := npcPosition.TK.Turrent.Angle-npcPosition.TK.Turrent.RangeAngle, npcPosition.TK.Turrent.Angle+npcPosition.TK.Turrent.RangeAngle

			if endAngle > 360 {
				endAngle -= 360
			}
			if startAngle < 0 {
				startAngle += 360
			}

			// 正常情况下 startAngle <= endAngle
			if startAngle <= endAngle {
				if startAngle <= angle && angle <= endAngle {
					npcPosition.TK.Enemy = playerPosition.TK
				}
			} else {
				// 如果处于 0 or 360的分割位置，startAngle > endAngle
				if angle <= endAngle || angle >= startAngle {
					npcPosition.TK.Enemy = playerPosition.TK
				}
			}
		}

		// 说明npc视野内没有敌人
		if npcPosition.TK.Enemy == nil {

			// 炮塔扫描
			npcPosition.TK.AddTurrentAngle(npcPosition.TK.Turrent.RotationSpeed)

			// 坦克移动方向：转向player的方向
			angle := math.Atan2(y, x) * 180 / math.Pi
			if angle < 0 {
				angle += 360.0
			}

			//	npcPosition.TK.Angle 表示 坦克 和 x 轴的夹角
			// angle 表示两个坦克连线 和 x轴的夹角
			if npcPosition.TK.Angle > angle {
				// 目的让 npcPosition.TK.Angle 往夹角小的方向移动，让炮台尽可能快的对准敌人
				if npcPosition.TK.Angle-angle > 180 {
					npcPosition.TK.AddTankAngle(1)
				} else {
					npcPosition.TK.AddTankAngle(-1)
				}
			} else if npcPosition.TK.Angle < angle {

				if angle-npcPosition.TK.Angle > 180 {
					npcPosition.TK.AddTankAngle(-1)
				} else {
					npcPosition.TK.AddTankAngle(1)
				}
			}

			npcPosition.TK.PreX, npcPosition.TK.PreY = npcPosition.TK.X, npcPosition.TK.Y
			// 移动坦克
			npcPosition.TK.X += npcPosition.TK.ForwardSpeed * math.Cos(npcPosition.TK.Angle*math.Pi/180)
			npcPosition.TK.Y += npcPosition.TK.ForwardSpeed * math.Sin(npcPosition.TK.Angle*math.Pi/180)
			// 更新碰撞盒子
			npcPosition.TK.updateTankCollisionBox()
		}
	}
}

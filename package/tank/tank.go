package tank

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/gofish2020/tankgame/package/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TankType int

const (
	ScreenToLogicScaleX = 5.12
	ScreenToLogicScaleY = 5.12

	minXCoordinates = 30
	minYCoordinates = 30

	TankTypePlayer TankType = iota
	TankTypeNPC
)

type Tank struct {
	X float64
	Y float64

	TkType    TankType
	ImagePath string

	// 🩸血量
	HealthPoints    int
	MaxHealthPoints int
	HealthBarWidth  float64
	HealthBarHeight float64

	// 炮弹装填
	ReloadTimer     int
	ReloadMaxTimer  int
	ReloadBarWidth  float64
	ReloadBarHeight float64

	// 旋转角度
	Angle float64
	// 角度变化速率
	RotationSpeed float64

	//前进速度
	ForwardSpeed float64
	// 后退速度
	BackwardSpeed float64

	// 炮塔参数
	Turrent Turret

	// 在攻击范围内的坦克
	Enemy *Tank
}

type Turret struct {
	Angle     float64
	ImagePath string

	// 炮塔旋转速度
	RotationSpeed float64

	//攻击范围
	RangeAngle    float64
	RangeDistance float64
}

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func NewTank(x, y float64, tankType TankType) *Tank {

	tank := Tank{

		X:         x,
		Y:         y,
		ImagePath: "resource/green_tank.png",

		TkType:        tankType,
		Angle:         0.0,
		RotationSpeed: 2.0,

		ReloadTimer:    0,
		ReloadMaxTimer: 100,

		ReloadBarWidth:  50,
		ReloadBarHeight: 5,

		HealthPoints:    100,
		MaxHealthPoints: 100,
		HealthBarWidth:  50,
		HealthBarHeight: 5,

		ForwardSpeed:  3.0,
		BackwardSpeed: 1.5,

		Turrent: Turret{
			Angle:         270.0,
			ImagePath:     "resource/green_tank_turret.png",
			RotationSpeed: 2.0,
		},
		Enemy: nil,
	}

	if tankType == TankTypePlayer {
		tank.Turrent.RangeAngle = 360.0
		tank.Turrent.RangeDistance = 300.0
	} else {
		tank.ImagePath = "resource/brown_tank.png"

		tank.Turrent.RangeAngle = 45.0
		tank.Turrent.RangeDistance = 100.0 + float64(r.Intn(300))
		tank.Turrent.ImagePath = "resource/brown_tank_turret.png"
	}

	return &tank
}

func (t *Tank) Update() {

	// 填充子弹
	if t.ReloadTimer < t.ReloadMaxTimer {
		t.ReloadTimer++
	}

	if t.TkType == TankTypePlayer { // 玩家坦克，手瞄

		if ebiten.IsKeyPressed(ebiten.KeySpace) && t.ReloadTimer == t.ReloadMaxTimer {
			t.ReloadTimer = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeyA) { // Press A
			t.Angle -= t.RotationSpeed
		} else if ebiten.IsKeyPressed(ebiten.KeyD) { // Press D
			t.Angle += t.RotationSpeed
		}

		if ebiten.IsKeyPressed(ebiten.KeyW) { // Press W

			t.X -= t.ForwardSpeed * math.Sin(-t.Angle*math.Pi/180)
			t.Y -= t.ForwardSpeed * math.Cos(-t.Angle*math.Pi/180)

		} else if ebiten.IsKeyPressed(ebiten.KeyS) { // Press S
			t.X += t.BackwardSpeed * math.Sin(-t.Angle*math.Pi/180)
			t.Y += t.BackwardSpeed * math.Cos(-t.Angle*math.Pi/180)

		}

		// 手动瞄准
		if ebiten.IsKeyPressed(ebiten.KeyJ) { // Press J
			t.Turrent.Angle -= t.Turrent.RotationSpeed
		} else if ebiten.IsKeyPressed(ebiten.KeyK) { // Press K
			t.Turrent.Angle += t.Turrent.RotationSpeed
		}

	} else { // npc tank 自瞄

		enemy := t.Enemy
		if enemy != nil { // 有敌人，自动瞄准

			x1, y1 := enemy.X, enemy.Y
			x2, y2 := t.X, t.Y

			// 计算夹角
			angle := float64(int(math.Atan2(y1-y2, x1-x2) / math.Pi * 180))
			// 角度限定在 [0,360]
			if angle < 0 {
				angle += 360
			}

			// 将 t.Turrent.Angle 限定在 [0,360]之间
			if t.Turrent.Angle >= 360 {
				t.Turrent.Angle -= 360
			} else if t.Turrent.Angle < 0 {
				t.Turrent.Angle += 360
			}

			// t.Turrent.Angle 表示炮塔和 x轴的夹角
			// angle 表示两个坦克连线 和 x轴的夹角
			if t.Turrent.Angle > angle {

				// 目的让t.Turrent.Angle 往哪个方向旋转（肯定是往夹角小的方向移动，让炮台尽可能快的对准敌人）
				if t.Turrent.Angle-angle > 180 {
					t.Turrent.Angle += 1
				} else {
					t.Turrent.Angle -= 1
				}
			} else if t.Turrent.Angle < angle {

				if angle-t.Turrent.Angle > 180 {
					t.Turrent.Angle -= 1
				} else {
					t.Turrent.Angle += 1
				}
			}
		}
	}

}

// 限制运行范围
func (t *Tank) LimitRange(minXCoordinates, minYCoordinates, maxXCoordinates, maxYCoordinates float64) {
	if t.X < minXCoordinates {
		t.X = minXCoordinates
	}
	if t.X > maxXCoordinates {
		t.X = maxXCoordinates
	}
	if t.Y < minYCoordinates {
		t.Y = minYCoordinates
	}
	if t.Y > maxYCoordinates {
		t.Y = maxYCoordinates
	}
}

// 绘制坦克各个元素
func (t *Tank) Draw(screen *ebiten.Image) {

	t.drawTank(screen)
	t.drawTurrent(screen)
	t.drawHealthBar(screen)
	t.drawReload(screen)
	t.drawAttackCircle(screen)
}

//........................基础元素绘制.....................

func (tk *Tank) drawAttackCircle(screen *ebiten.Image) {

	clr := color.RGBA{0, 255, 0, 128}
	if tk.Enemy != nil {
		clr = color.RGBA{255, 0, 0, 128}
	}

	if tk.TkType == TankTypePlayer {
		// player 才有提示圈
		vector.StrokeCircle(screen, float32(tk.X), float32(tk.Y), float32(tk.Turrent.RangeDistance), 1.0, clr, true)
	} else {
		startAngle, endAngle := (tk.Turrent.Angle-tk.Turrent.RangeAngle)*math.Pi/180, (tk.Turrent.Angle+tk.Turrent.RangeAngle)*math.Pi/180
		utils.DrawSector(screen, float32(tk.X), float32(tk.Y), 1.0, float32(tk.Turrent.RangeDistance), float32(startAngle), float32(endAngle), clr, false)
	}
}

// 坦克
func (tk *Tank) drawTank(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	// 加载图片
	tankBody, _, _ := ebitenutil.NewImageFromFile(tk.ImagePath)

	baseOffsetX := float64(tankBody.Bounds().Dx()) / 2 // hullBody.Bounds().Dx() = 256
	baseOffsetY := float64(tankBody.Bounds().Dy()) / 2 // hullBody.Bounds().Dy() = 256

	// 先平移图片（将图片的中心，移动到（0，0）位置）
	op.GeoM.Translate(-baseOffsetX, -baseOffsetY)
	// 旋转图片
	op.GeoM.Rotate(tk.Angle * math.Pi / 180.0)
	// 再平移图片到窗口的中心位置 （ 因为绘制收缩了，所以屏幕坐标需要增大）
	op.GeoM.Translate(tk.X*ScreenToLogicScaleX, tk.Y*ScreenToLogicScaleY)
	// 整个绘制收缩了（ 50 / 256）倍，即 1/5.12
	op.GeoM.Scale(1/ScreenToLogicScaleX, 1/ScreenToLogicScaleY)
	// 绘制图片
	screen.DrawImage(tankBody, op)

}

// 绘制炮塔
func (tk *Tank) drawTurrent(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	turrentBody, _, _ := ebitenutil.NewImageFromFile(tk.Turrent.ImagePath)

	baseOffsetX := float64(turrentBody.Bounds().Dx()) / 2 // hullBody.Bounds().Dx() = 256
	baseOffsetY := float64(turrentBody.Bounds().Dy()) / 2 // hullBody.Bounds().Dy() = 256
	// 先平移图片（将图片的中心，移动到（0，0）位置）
	op.GeoM.Translate(-baseOffsetX, -baseOffsetY)
	// 旋转图片
	op.GeoM.Rotate(tk.Turrent.Angle * math.Pi / 180.0)
	// 再平移图片到窗口的中心位置 （ 因为绘制收缩了，所以屏幕坐标需要增大）
	op.GeoM.Translate(tk.X*ScreenToLogicScaleX, tk.Y*ScreenToLogicScaleY)
	// 整个绘制收缩了（ 50 / 256）倍，即 1/5.12
	op.GeoM.Scale(1/ScreenToLogicScaleX, 1/ScreenToLogicScaleY)
	// 绘制图片
	screen.DrawImage(turrentBody, op)
}

// 血条
func (tk *Tank) drawHealthBar(screen *ebiten.Image) {

	// 血量百分比
	percentage := float64(tk.HealthPoints) / float64(tk.MaxHealthPoints)

	// 血量颜色
	var filledColor color.RGBA
	if percentage >= 0.60 {
		filledColor = color.RGBA{0, 255, 0, 255} // Green
	} else if percentage >= 0.40 {
		filledColor = color.RGBA{255, 165, 0, 255} // Orange
	} else if percentage > 0 {
		filledColor = color.RGBA{255, 0, 0, 255} // Red
	} else {
		filledColor = color.RGBA{0, 0, 0, 0} // Transparent
	}

	filledWidth := 1 + int(tk.HealthBarWidth*percentage)

	newImage := ebiten.NewImage(filledWidth, int(tk.HealthBarHeight))
	newImage.Fill(filledColor)

	op := &ebiten.DrawImageOptions{}
	// tk.X-25.5 左对齐坦卡边缘
	op.GeoM.Translate(tk.X-25.5, tk.Y+30)
	screen.DrawImage(newImage, op)

}

// 重新装弹
func (tk *Tank) drawReload(screen *ebiten.Image) {
	percentage := float64(tk.ReloadTimer) / float64(tk.ReloadMaxTimer)

	var filledColor color.RGBA = color.RGBA{128, 128, 128, 255} // grey

	if tk.ReloadTimer == tk.ReloadMaxTimer { // 满了
		filledColor = color.RGBA{255, 105, 180, 255}
	}

	filledWidth := 1 + int(tk.ReloadBarWidth*percentage)
	newImage := ebiten.NewImage(filledWidth, int(tk.ReloadBarHeight))
	newImage.Fill(filledColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(tk.X-25.5, tk.Y+35)
	screen.DrawImage(newImage, op)
}

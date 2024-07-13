package tank

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/gofish2020/tankgame/package/utils"
	"github.com/gofish2020/tankgame/package/utils/sound"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TankType int

const (
	ScreenToLogicScaleX = 5.12 // 图片是 256大小，希望缩为 50
	ScreenToLogicScaleY = 5.12

	MinXCoordinates = 30.0
	MinYCoordinates = 30.0

	TankTypePlayer TankType = iota
	TankTypeNPC
)

type Tank struct {
	X      float64
	Y      float64
	Width  float64 // 宽度
	Height float64 // 高度

	// 记录前一个位置，当做碰撞检测时候，回撤到前一个位置
	PreX float64
	PreY float64

	Name string

	TkType    TankType // 坦克的操作者
	ImagePath string   // 坦克图片

	// 🩸血量
	HealthPoints    int
	MaxHealthPoints int
	HealthBarWidth  float64
	HealthBarHeight float64

	// 炮弹装填
	ReloadTimer    int
	ReloadMaxTimer int
	ReloadSpeed    int

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

	// 四个角，旋转后的坐标（做碰撞检测）
	// 顺时针，左上
	CollisionX1 float64
	CollisionY1 float64
	// 右上
	CollisionX2 float64
	CollisionY2 float64
	// 右下
	CollisionX3 float64
	CollisionY3 float64
	// 左下
	CollisionX4 float64
	CollisionY4 float64

	// 炮塔参数
	Turrent Turret

	// 在攻击范围内的坦克
	Enemy *Tank

	Projectiles []*Projectile // 发射的炮弹

	LastCollisionTime time.Time
}

// 炮弹
type Projectile struct {
	X         float64 // 炮弹坐标 X
	Y         float64 // 炮弹坐标 Y
	Speed     float64 // 运行速率
	Angle     float64 // 移动方向
	Width     float64 // 宽度
	Height    float64 // 高度
	IsExplode bool    // 是否已碰撞

	Frame int // 爆炸图片遍历使用
}

// 炮塔
type Turret struct {
	Angle     float64
	ImagePath string

	// 炮塔旋转速度
	RotationSpeed float64

	//攻击范围
	RangeAngle    float64
	RangeDistance float64

	//子弹速率
	ProjectileSpeed float64
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

		Width:  50, // 坦克的宽
		Height: 50, // 坦克的高

		TkType:        tankType,
		Angle:         270.0,
		RotationSpeed: 2.0,

		ForwardSpeed:  5.0,
		BackwardSpeed: 3.5,

		ReloadTimer:    0,
		ReloadMaxTimer: 100,
		ReloadSpeed:    1.0,

		ReloadBarWidth:  50,
		ReloadBarHeight: 5,

		HealthPoints:    200,
		MaxHealthPoints: 200,
		HealthBarWidth:  50,
		HealthBarHeight: 5,

		Turrent: Turret{
			Angle:           270.0, // 默认指向上
			ImagePath:       "resource/green_tank_turret.png",
			RotationSpeed:   2.0,
			ProjectileSpeed: 30.0,
		},

		Projectiles: nil,
		Enemy:       nil,
	}

	if tankType == TankTypePlayer {
		tank.Turrent.RangeAngle = 360.0
		tank.Turrent.RangeDistance = 300.0
		tank.Name = "ikun"
		tank.ReloadSpeed = 2.0
	} else {

		var level utils.TankLevel // 随机坦克的速度
		if utils.GameLevel <= 3 {
			level = utils.TankLevels[r.Intn(utils.GameLevel*6)]
		} else {
			level = utils.TankLevels[r.Intn(len(utils.TankLevels))]
		}

		tank.ImagePath = "resource/brown_tank.png"
		tank.MaxHealthPoints = 50
		tank.HealthPoints = 50
		tank.Angle = 90.0
		tank.ForwardSpeed = level.TankSpeed                   // 前进速度
		tank.Turrent.RotationSpeed = level.TurrentRotateSpeed // 炮塔旋转速度

		tank.Turrent.RangeAngle = 45.0                            // 攻击视角
		tank.Turrent.RangeDistance = 100.0 + float64(r.Intn(300)) // 攻击范围
		tank.Turrent.ImagePath = "resource/brown_tank_turret.png"
		tank.Turrent.Angle = 90.0 // 敌人默认指向下
		tank.Name = enemyNames[r.Intn(len(enemyNames))]
	}
	// 更新坦克的四个顶点坐标
	tank.updateTankCollisionBox()
	return &tank
}

func (t *Tank) DeathSound() {

	soundName := strconv.Itoa(utils.KilledCount)
	if utils.KilledCount > 5 {
		soundName = "dead" + strconv.Itoa(rand.Intn(4)+1)
	}
	sound.PlaySound(soundName)
}
func (t *Tank) shot() {
	// 能量满，才能射击
	if t.ReloadTimer == t.ReloadMaxTimer {
		if t.TkType == TankTypePlayer { // player
			if utils.GameProgress == "pass" {
				sound.PlaySound("dog")
			} else {
				sound.PlaySound("boom")
			}
		}

		t.ReloadTimer = 0
		// 生成炮弹
		newProjectile := Projectile{
			X:         t.X,                       // 炮弹初始X
			Y:         t.Y,                       // 炮弹初始Y
			Angle:     t.Turrent.Angle,           // 初始角度（就是炮塔的角度）
			IsExplode: false,                     // 是否已经爆炸
			Speed:     t.Turrent.ProjectileSpeed, // 炮弹移动速度
		}
		t.Projectiles = append(t.Projectiles, &newProjectile)
	}

}

// 目的在于让 炮塔的角度始终使用 正度数 表示 [0,360]之间
func (t *Tank) AddTurrentAngle(duration float64) {

	t.Turrent.Angle += duration
	if t.Turrent.Angle >= 360.0 { // 超过360，转成360度范围
		t.Turrent.Angle -= 360.0
	} else if t.Turrent.Angle < 0 { // 负数转正数
		t.Turrent.Angle += 360.0
	}
}

func (t *Tank) AddTankAngle(duration float64) {

	t.Angle += duration
	if t.Angle >= 360.0 { // 超过360，转成360度范围
		t.Angle -= 360.0
	} else if t.Angle < 0 { // 负数转正数
		t.Angle += 360.0
	}
}

func (t *Tank) Update() {

	// 填充子弹
	if t.ReloadTimer < t.ReloadMaxTimer {
		t.ReloadTimer += t.ReloadSpeed
		if t.ReloadTimer > t.ReloadMaxTimer {
			t.ReloadTimer = t.ReloadMaxTimer
		}
	}

	if t.TkType == TankTypePlayer { // 玩家坦克，手瞄

		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			t.shot()
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) { // Press A
			t.AddTankAngle(-t.RotationSpeed)
			t.updateTankCollisionBox()
		} else if ebiten.IsKeyPressed(ebiten.KeyD) { // Press D

			t.AddTankAngle(t.RotationSpeed)
			t.updateTankCollisionBox()
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) { // Press W
			t.PreX, t.PreY = t.X, t.Y // 记录前一个位置，当做碰撞检测时候，来回撤到前一个位置

			t.X += t.ForwardSpeed * math.Cos(t.Angle*math.Pi/180)
			t.Y += t.ForwardSpeed * math.Sin(t.Angle*math.Pi/180)
			t.updateTankCollisionBox()
		} else if ebiten.IsKeyPressed(ebiten.KeyS) { // Press S

			t.PreX, t.PreY = t.X, t.Y

			t.Y -= t.BackwardSpeed * math.Sin(t.Angle*math.Pi/180)
			t.X -= t.BackwardSpeed * math.Cos(t.Angle*math.Pi/180)
			t.updateTankCollisionBox()
		}

		// 手动瞄准
		if ebiten.IsKeyPressed(ebiten.KeyJ) { // Press J
			t.AddTurrentAngle(-t.Turrent.RotationSpeed)
		} else if ebiten.IsKeyPressed(ebiten.KeyK) { // Press K
			t.AddTurrentAngle(t.Turrent.RotationSpeed)
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

			// t.Turrent.Angle 表示炮塔和 x轴的夹角
			// angle 表示两个坦克连线 和 x轴的夹角
			if t.Turrent.Angle > angle {
				// 目的让t.Turrent.Angle 往夹角小的方向移动，让炮台尽可能快的对准敌人
				if t.Turrent.Angle-angle > 180 {
					t.AddTurrentAngle(1)
				} else {
					t.AddTurrentAngle(-1)
				}
			} else if t.Turrent.Angle < angle {

				if angle-t.Turrent.Angle > 180 {
					t.AddTurrentAngle(-1)
				} else {
					t.AddTurrentAngle(1)
				}
			} else {
				// 这里精准瞄准，立刻射击
				t.shot()
			}

			//t.shot() // 不管是否瞄准，就射击
		}
	}

	// 更新炮弹的移动
	t.updateProjectile()

}

// 更新坦克的四个顶点边界
func (t *Tank) updateTankCollisionBox() {

	// 用来作为坦克四个角的初始坐标
	offsetX := float64(t.Width) / 2
	offsetY := float64(t.Height) / 2

	// t.X t.Y 矩形的中心点
	// 左上角 (x = -offsetX  y = -offsetY)
	t.CollisionX1, t.CollisionY1 = rotatePoint(t.X-offsetX, t.Y-offsetY, t.Angle, t.X, t.Y)
	//  右上角 (x = offsetX y = -offsetY )
	t.CollisionX2, t.CollisionY2 = rotatePoint(t.X+offsetX, t.Y-offsetY, t.Angle, t.X, t.Y)
	// 右下角 (x = offsetX y = offsetY)
	t.CollisionX3, t.CollisionY3 = rotatePoint(t.X+offsetX, t.Y+offsetY, t.Angle, t.X, t.Y)
	// 左下角 (x = -offsetX y=offsetY)
	t.CollisionX4, t.CollisionY4 = rotatePoint(t.X-offsetX, t.Y+offsetY, t.Angle, t.X, t.Y)

}

// 点 x/y 围绕点 cx/cy 旋转 angle 角度后的坐标
func rotatePoint(x, y, angle, cx, cy float64) (float64, float64) {

	// 角度转弧度
	angleRad := angle * math.Pi / 180
	cosAngle := math.Cos(angleRad)
	sinAngle := math.Sin(angleRad)

	// 平移点到原点
	x -= cx
	y -= cy

	/*
		矩阵旋转公式：
		x' = xCos(θ) - ySin(θ)
		y' = xSin(θ) + ycos(θ)
	*/
	xNew := x*cosAngle - y*sinAngle
	yNew := x*sinAngle + y*cosAngle

	// 平移回去
	xNew += cx
	yNew += cy

	return xNew, yNew
}

// 限制运行范围
func (t *Tank) LimitTankRange(minXCoordinates, minYCoordinates, maxXCoordinates, maxYCoordinates float64) {
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

// 更新炮弹的移动
func (t *Tank) updateProjectile() {

	for idx, projectile := range t.Projectiles {

		// 检查炮弹是否已经飞出去边界
		if projectile.X < 0 || projectile.X > monitor.ScreenWidth || projectile.Y < 0 || projectile.Y > monitor.ScreenHeight {
			// 删除炮弹
			t.removeProjectile(idx)
			continue
		}

		if projectile.IsExplode { // 炮弹已经爆炸
			if projectile.Frame > 16 { // 爆炸效果结束
				t.removeProjectile(idx) // 删除炮弹
			} else {
				projectile.Frame++ // 爆炸效果
			}
			continue
		}
		// 转为弧度
		angleRadians := projectile.Angle * math.Pi / 180.0
		// 水平和垂直分量计算
		offsetX := projectile.Speed * math.Cos(angleRadians)
		offsetY := projectile.Speed * math.Sin(angleRadians)
		// 累加
		projectile.X += offsetX
		projectile.Y += offsetY

	}
}

// 删除炮弹
func (t *Tank) removeProjectile(index int) {
	// Ensure the index is within bounds
	if index < 0 || index >= len(t.Projectiles) {
		return
	}
	t.Projectiles = append(t.Projectiles[:index], t.Projectiles[index+1:]...)
}

//........................基础元素绘制.....................

var (
	projectileImage, _, _ = ebitenutil.NewImageFromFile("resource/projectile.png")
	explosionImg, _, _    = ebitenutil.NewImageFromFile("resource/explosion.png")
)

// 绘制坦克各个元素
func (t *Tank) Draw(screen *ebiten.Image) {

	t.drawTank(screen)
	t.drawTurrent(screen)
	t.drawHealthBar(screen)
	t.drawReload(screen)
	t.drawAttackCircle(screen)
	t.drawProjectile(screen)

}

// 绘制炮弹
func (tk *Tank) drawProjectile(screen *ebiten.Image) {

	frameOX := 0
	frameOY := 0
	frameWidth := 64
	frameHeight := 64
	frameCount := 16
	for _, projectile := range tk.Projectiles {

		if projectile.IsExplode { // 绘制爆炸特效

			frameIndex := projectile.Frame % frameCount
			if frameIndex < 0 || frameIndex >= frameCount {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(projectile.X, projectile.Y)
			// 按照一列一列显示图片
			sy := frameOY + (frameIndex/4)*frameHeight
			sx := frameOX + (frameIndex%4)*frameWidth
			// 裁剪图片
			subImg := explosionImg.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
			screen.DrawImage(subImg, op)

		} else { // 绘制炮弹正常飞行
			op := &ebiten.DrawImageOptions{}

			baseOffsetX := float64(projectileImage.Bounds().Dx()) / 2
			baseOffsetY := float64(projectileImage.Bounds().Dy()) / 2

			// 先平移图片（将图片的中心，移动到（0，0）位置）
			op.GeoM.Translate(-baseOffsetX, -baseOffsetY)
			// 旋转图片
			op.GeoM.Rotate(projectile.Angle * math.Pi / 180.0)

			// 再平移图片到窗口的中心位置 （ 因为绘制收缩了，所以屏幕坐标需要增大）
			op.GeoM.Translate(projectile.X, projectile.Y)
			// 绘制图片
			screen.DrawImage(projectileImage, op)
		}

	}
}

func (tk *Tank) drawAttackCircle(screen *ebiten.Image) {

	clr := color.RGBA{255, 248, 220, 100}
	if tk.Enemy != nil {
		clr = color.RGBA{255, 69, 0, 100}
	}

	if tk.TkType == TankTypePlayer {
		// player 圆圈
		//vector.StrokeCircle(screen, float32(tk.X), float32(tk.Y), float32(tk.Turrent.RangeDistance), 1.0, clr, true)
	} else {
		startAngle, endAngle := (tk.Turrent.Angle-tk.Turrent.RangeAngle)*math.Pi/180, (tk.Turrent.Angle+tk.Turrent.RangeAngle)*math.Pi/180
		utils.DrawSector(screen, float32(tk.X), float32(tk.Y), 1.0, float32(tk.Turrent.RangeDistance), float32(startAngle), float32(endAngle), clr, true)
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
	// 整个绘制收缩了（ 50 / 256）倍，即 1/5.12
	op.GeoM.Scale(1/ScreenToLogicScaleX, 1/ScreenToLogicScaleY)
	// 再平移图片到窗口的中心位置 （ 因为绘制收缩了，所以屏幕坐标需要增大）
	op.GeoM.Translate(tk.X, tk.Y)
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

	// 整个绘制收缩了（ 50 / 256）倍，即 1/5.12
	op.GeoM.Scale(1/ScreenToLogicScaleX, 1/ScreenToLogicScaleY)
	// 再平移图片到窗口的中心位置 （ 因为绘制收缩了，所以屏幕坐标需要增大）
	op.GeoM.Translate(tk.X, tk.Y)
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

package main

import (
	"log"

	"github.com/gofish2020/tankgame/package/game"
	"github.com/gofish2020/tankgame/package/monitor"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

func main() {

	ebiten.SetRunnableOnUnfocused(true) // 游戏界面不显示，依然运行
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetTPS(50)            // 窗口刷新频率
	ebiten.SetVsyncEnabled(true) // 垂直同步
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowTitle("Tank Shoot")
	ebiten.SetWindowSize(int(monitor.ScreenWidth), int(monitor.ScreenHeight))
	ebiten.SetWindowFloating(true)          // 置顶显示
	ebiten.SetWindowMousePassthrough(false) // 鼠标穿透

	// 需要提前调用一下，不然没有声音
	audio.NewContext(44100)
	audio.CurrentContext().NewPlayerFromBytes([]byte{}).Play() // 类似于预热的感觉（可能是库有bug）

	game := game.NewGame()
	err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		InitUnfocused:     true, // 启动时候，窗体不聚焦
		ScreenTransparent: true, // 窗体透明
		SkipTaskbar:       true, // 图片不显示在任务栏
		X11ClassName:      "Tank Shoot",
		X11InstanceName:   "Tank Shoot",
	})
	if err != nil {
		log.Fatal(err)
	}
}

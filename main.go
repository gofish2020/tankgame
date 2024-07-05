package main

import (
	"log"

	"github.com/gofish2020/tankgame/package/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	ebiten.SetRunnableOnUnfocused(true) // 游戏界面不显示，依然运行
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetTPS(50)            // 窗口刷新频率
	ebiten.SetVsyncEnabled(true) // 垂直同步
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowTitle("Tank Shot")
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowFloating(true) // 置顶显示
	ebiten.SetWindowMousePassthrough(true)

	game := game.NewGame()
	err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		InitUnfocused:     true, // 启动时候，窗体不聚焦
		ScreenTransparent: true, // 窗体透明
		SkipTaskbar:       true, // 图片不显示在任务栏
		X11ClassName:      "Tank Shot",
		X11InstanceName:   "Tank Shot",
	})
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	config "yikong/internal/app"
	constant "yikong/internal/constants"
	"yikong/internal/logging"
	"yikong/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	if !config.InspectIsADBExisted() {
		logging.Error("未找到ADB, 请先安装ADB")
		return
	}

	application := app.New()
	window := application.NewWindow(constant.AppName)

	mainUi := ui.MainUI(window)

	window.Resize(fyne.NewSize(1200, 600))
	window.SetFixedSize(true)
	window.SetContent(mainUi)
	window.ShowAndRun()
}

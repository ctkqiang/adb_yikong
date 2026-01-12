package main

import (
	config "yikong/internal/app"
	constant "yikong/internal/constants"
	"yikong/internal/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	if !config.InspectIsADBExisted() {
		logging.Error("未找到ADB, 请先安装ADB")
		return
	}

	application := app.New()
	window := application.NewWindow(constant.AppName)

	hello := widget.NewLabel(constant.AppName)

	window.Resize(fyne.NewSize(1200, 600))
	window.SetFixedSize(true)

	window.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	window.ShowAndRun()

}

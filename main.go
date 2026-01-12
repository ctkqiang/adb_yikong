package main

import (
	"yikong/config"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	application := app.New()
	window := application.NewWindow(config.AppName)

	hello := widget.NewLabel(config.AppName)

	window.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	window.ShowAndRun()

}

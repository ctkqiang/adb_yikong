package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MainUI(window fyne.Window) fyne.CanvasObject {
	hello := widget.NewLabel("ADB调试器")

	sidebar := container.NewVBox(
		widget.NewButton("设备", func() { hello.SetText("设备管理") }),
		widget.NewButton("应用", func() { hello.SetText("应用管理") }),
		widget.NewButton("日志", func() { hello.SetText("日志查看") }),
		widget.NewButton("文件", func() { hello.SetText("文件传输") }),
	)

	return container.NewHSplit(
		container.NewBorder(
			widget.NewLabel("功能菜单"),
			nil, nil, nil,
			sidebar,
		),
		container.NewVBox(
			hello,
			widget.NewButton("Hi!", func() {
				hello.SetText("欢迎使用ADB调试器 :)")
			}),
		),
	)
}

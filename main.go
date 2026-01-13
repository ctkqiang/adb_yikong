package main

import (
	_ "embed"
	"os"
	config "yikong/internal/app"
	constant "yikong/internal/constants"
	"yikong/internal/logging"
	"yikong/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/logo.png
var iconBytes []byte

func main() {
	if _, err := os.Stat("assets/logo.png"); os.IsNotExist(err) {
		logging.Error("图标文件不存在: %s", err)
	} else {
		info, _ := os.Stat("assets/logo.png")
		logging.Info("图标文件大小: %d bytes", info.Size())
	}

	logging.Info("iconBytes embed 结果: 长度=%d, 是否为nil=%v", len(iconBytes), iconBytes == nil)

	if !config.InspectIsADBExisted() {
		logging.Error("未找到ADB, 请先安装ADB")
		return
	}

	application := app.New()
	window := application.NewWindow(constant.AppName)

	var appIcon fyne.Resource

	if len(iconBytes) == 0 {
		logging.Warn("图标文件为空，使用默认图标")
		appIcon = theme.ComputerIcon()
	} else {
		appIcon = &fyne.StaticResource{
			StaticName:    "logo.png",
			StaticContent: iconBytes,
		}
	}

	application.SetIcon(appIcon)
	window.SetIcon(appIcon)

	mainUi := ui.MainUI(window)

	if err := ui.SystemTray(application, window); err != nil {
		logging.Error("%s", "系统托盘初始化失败: "+err.Error())
	}

	window.Resize(fyne.NewSize(1200, 600))
	window.SetFixedSize(true)
	window.SetContent(mainUi)

	window.SetCloseIntercept(func() {
		window.Hide()
	})

	window.ShowAndRun()
}

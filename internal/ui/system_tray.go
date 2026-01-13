package ui

import (
	"yikong/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

func SystemTray(application fyne.App, window fyne.Window) error {
	if desk, ok := application.(desktop.App); ok {

		icon := theme.ComputerIcon()

		menu := fyne.NewMenu(constants.AppName,
			fyne.NewMenuItem("显示", func() {
				window.Show()
			}),
			fyne.NewMenuItem("退出", func() {
				application.Quit()
			}),
		)

		desk.SetSystemTrayIcon(icon)
		desk.SetSystemTrayMenu(menu)
	}
	return nil
}

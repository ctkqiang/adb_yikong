package ui

import (
	"yikong/internal/constants"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var customIcon fyne.Resource

func SetCustomIcon(icon fyne.Resource) {
	customIcon = icon
}

func InfoDialog(window fyne.Window, message string, isAbout bool) error {
	if isAbout {
		content := widget.NewLabel(message)
		aboutDialog := dialog.NewCustom(constants.AppName, "确定", content, window)
		if customIcon != nil {
			aboutDialog.SetIcon(customIcon)
		}
		aboutDialog.Show()
	} else {
		dialog.ShowInformation(constants.AppName, message, window)
	}

	return nil
}

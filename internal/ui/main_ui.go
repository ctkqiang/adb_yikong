package ui

import (
	"image/color"
	"time"
	"yikong/internal/adb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	window         fyne.Window
	mainContainer  *fyne.Container
	deviceList     *widget.List
	devices        []adb.Device
	selectedDevice string
}

func MainUI(window fyne.Window) fyne.CanvasObject {
	ui := &UI{
		window: window,
	}
	return ui.createUI()
}

func (ui *UI) createUI() fyne.CanvasObject {
	devicePage := ui.createDevicePage()
	functionPage := ui.createFunctionPage()
	functionPage.Hide()

	ui.mainContainer = container.NewStack(devicePage, functionPage)

	return ui.mainContainer
}

func (ui *UI) createDevicePage() fyne.CanvasObject {
	title := widget.NewLabel("已连接设备")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	subtitle := widget.NewLabel("选择一个设备以继续")
	subtitle.TextStyle = fyne.TextStyle{Italic: true}
	subtitle.Alignment = fyne.TextAlignCenter

	header := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
	)

	refreshBtn := widget.NewButton("刷新设备", func() {
		ui.refreshDevices()
	})
	refreshBtn.Importance = widget.HighImportance

	refreshBg := canvas.NewRectangle(color.White)
	refreshBg.CornerRadius = 20
	refreshBtnContainer := container.NewStack(refreshBg, refreshBtn)

	ui.deviceList = widget.NewList(
		func() int {
			return len(ui.devices)
		},
		func() fyne.CanvasObject {
			card := canvas.NewRectangle(color.White)
			card.CornerRadius = 20

			deviceIcon := widget.NewIcon(theme.ComputerIcon())
			deviceIcon.Resize(fyne.NewSize(32, 32))

			deviceName := widget.NewLabel("设备名称")
			deviceName.TextStyle = fyne.TextStyle{Bold: true}

			deviceID := widget.NewLabel("设备ID")
			deviceID.TextStyle = fyne.TextStyle{Bold: false}

			deviceInfo := container.NewVBox(deviceName, deviceID)
			content := container.NewBorder(nil, nil, deviceIcon, nil, deviceInfo)

			selectBtn := widget.NewButton("选择", nil)
			selectBtn.Importance = widget.HighImportance

			return container.NewStack(card, container.NewBorder(nil, selectBtn, nil, nil, container.NewPadded(content)))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			stack := o.(*fyne.Container)
			border := stack.Objects[1].(*fyne.Container)
			paddedContent := border.Objects[0].(*fyne.Container)
			content := paddedContent.Objects[0].(*fyne.Container)
			deviceInfo := content.Objects[0].(*fyne.Container)
			deviceName := deviceInfo.Objects[0].(*widget.Label)
			deviceID := deviceInfo.Objects[1].(*widget.Label)
			selectBtn := border.Objects[1].(*widget.Button)
			card := stack.Objects[0].(*canvas.Rectangle)

			deviceName.SetText(ui.devices[i].Name)
			deviceID.SetText(ui.devices[i].ID)

			if ui.selectedDevice == ui.devices[i].ID {
				card.FillColor = color.White
				card.StrokeColor = color.White
				card.StrokeWidth = 3
			} else {
				card.FillColor = color.White
				card.StrokeColor = color.White
				card.StrokeWidth = 2
			}
			card.Refresh()

			selectBtn.OnTapped = func() {
				ui.selectDevice(ui.devices[i].ID)
			}
		},
	)

	ui.refreshDevices()

	mainContent := container.NewBorder(
		header,
		refreshBtnContainer,
		nil, nil,
		ui.deviceList,
	)

	return container.NewStack(
		container.NewPadded(mainContent),
	)
}

func (ui *UI) createFunctionPage() fyne.CanvasObject {
	backBtn := widget.NewButton("返回设备列表", func() {
		ui.showDevicePage()
	})
	backBtn.Importance = widget.MediumImportance

	backBg := canvas.NewRectangle(color.NRGBA{R: 255, G: 182, B: 193, A: 255})
	backBg.CornerRadius = 20
	backBtnContainer := container.NewStack(backBg, backBtn)

	deviceInfo := widget.NewCard(
		"当前设备",
		"",
		widget.NewLabel("等待选择功能..."),
	)

	featureButtons := []struct {
		name     string
		icon     fyne.Resource
		callback func()
	}{
		{"设备管理", theme.SettingsIcon(), func() { deviceInfo.SetContent(widget.NewLabel("设备管理功能")) }},
		{"应用管理", theme.DocumentSaveIcon(), func() { deviceInfo.SetContent(widget.NewLabel("应用管理功能")) }},
		{"日志查看", theme.DocumentPrintIcon(), func() { deviceInfo.SetContent(widget.NewLabel("日志查看功能")) }},
		{"文件传输", theme.MailSendIcon(), func() { deviceInfo.SetContent(widget.NewLabel("文件传输功能")) }},
		{"设置", theme.SettingsIcon(), func() { deviceInfo.SetContent(widget.NewLabel("设置功能")) }},
	}

	buttonGrid := container.NewGridWithColumns(2,
		container.NewVBox(),
		container.NewVBox(),
	)

	for idx, feature := range featureButtons {
		btnContainer := ui.createFeatureButton(feature.name, feature.icon, feature.callback)
		if idx < 3 {
			buttonGrid.Objects[0].(*fyne.Container).Add(btnContainer)
		} else {
			buttonGrid.Objects[1].(*fyne.Container).Add(btnContainer)
		}
	}

	mainContent := container.NewBorder(
		container.NewVBox(
			backBtnContainer,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		container.NewVBox(
			deviceInfo,
			widget.NewSeparator(),
			buttonGrid,
		),
	)

	return container.NewStack(
		container.NewPadded(mainContent),
	)
}

func (ui *UI) createFeatureButton(text string, icon fyne.Resource, onTap func()) *fyne.Container {
	btn := widget.NewButtonWithIcon(text, icon, onTap)
	btn.Importance = widget.HighImportance

	bg := canvas.NewRectangle(color.White)
	bg.CornerRadius = 20

	decoratedBtn := container.NewStack(bg, btn)

	btn.OnTapped = func() {
		ripple := canvas.NewCircle(color.NRGBA{R: 255, G: 182, B: 193, A: 100})
		ripple.Resize(fyne.NewSize(20, 20))
		ripple.Hide()

		decoratedBtn.Add(ripple)
		ripple.Show()

		go func() {
			for i := 0; i < 10; i++ {
				time.Sleep(50 * time.Millisecond)
				size := float32(i*5 + 20)
				fyne.Do(func() {
					ripple.Resize(fyne.NewSize(size, size))
					ripple.Move(fyne.NewPos(btn.Size().Width/2-size/2, btn.Size().Height/2-size/2))
				})
			}
			fyne.Do(func() {
				ripple.Hide()
				decoratedBtn.Remove(ripple)
			})
		}()

		if onTap != nil {
			onTap()
		}
	}

	return decoratedBtn
}

func (ui *UI) refreshDevices() {
	devices, err := adb.GetDevices()
	if err != nil {
		errorCard := widget.NewCard("错误", "获取设备失败: "+err.Error(),
			widget.NewButton("确定", func() {}),
		)
		dialog := widget.NewModalPopUp(errorCard, ui.window.Canvas())
		dialog.Show()
		return
	}

	ui.devices = devices
	ui.deviceList.Refresh()
}

func (ui *UI) selectDevice(deviceID string) {
	ui.selectedDevice = deviceID
	ui.showFunctionPage()
}

func (ui *UI) showDevicePage() {
	currentPage := ui.mainContainer.Objects[1]
	newPage := ui.mainContainer.Objects[0]

	if rect, ok := currentPage.(*fyne.Container).Objects[0].(*canvas.Rectangle); ok {
		go func() {
			for i := 255; i >= 0; i -= 15 {
				rect.FillColor = color.NRGBA{R: 255, G: 240, B: 245, A: uint8(i)}
				rect.Refresh()
				time.Sleep(20 * time.Millisecond)
			}
			currentPage.Hide()
			newPage.Show()
			ui.mainContainer.Refresh()
		}()
	} else {
		for _, obj := range ui.mainContainer.Objects {
			obj.Hide()
		}
		newPage.Show()
		ui.mainContainer.Refresh()
	}
}

func (ui *UI) showFunctionPage() {
	for _, obj := range ui.mainContainer.Objects {
		obj.Hide()
	}
	ui.mainContainer.Objects[1].Show()
	ui.mainContainer.Refresh()
}

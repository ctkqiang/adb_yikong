package ui

import (
	"image/color"
	"strconv"
	"time"
	"yikong/internal/adb"
	"yikong/internal/constants"

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

	iconMap := map[string]func() fyne.Resource{
		"SettingsIcon":      theme.SettingsIcon,
		"DocumentSaveIcon":  theme.DocumentSaveIcon,
		"DocumentPrintIcon": theme.DocumentPrintIcon,
		"MailSendIcon":      theme.MailSendIcon,
	}

	buttonGrid := container.NewGridWithColumns(2,
		container.NewVBox(),
		container.NewVBox(),
	)

	featureOrder := []string{"device_management", "app_management", "log_viewing", "file_transfer", "settings"}

	for idx, featureID := range featureOrder {
		config, exists := constants.FeatureMap[featureID]
		if !exists {
			// 如果配置不存在，使用默认值
			defaultConfig := constants.FeatureConfig{
				ID:           featureID,
				Name:         "未知功能",
				Description:  "功能配置不存在",
				IconName:     "SettingsIcon",
				CommandGroup: []string{},
				DefaultLabel: "功能未配置",
			}
			config = defaultConfig
		}

		iconFunc, iconExists := iconMap[config.IconName]

		var icon fyne.Resource

		if iconExists {
			icon = iconFunc()
		} else {
			icon = theme.SettingsIcon() // 默认图标
		}

		// 创建回调函数，显示功能详情
		callback := func(featureConfig constants.FeatureConfig) func() {
			return func() {
				// 创建功能详情界面
				content := ui.createFeatureDetailContent(featureConfig)
				deviceInfo.SetContent(content)
			}
		}(config)

		btnContainer := ui.createFeatureButton(config.Name, icon, callback)
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

func (ui *UI) createDeviceManagementPage() fyne.CanvasObject {
	// 创建标题
	title := widget.NewLabel("设备管理")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// 创建设备列表容器
	devicesLabel := widget.NewLabel("已连接设备:")
	devicesLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 创建设备列表显示
	var deviceListContainer *fyne.Container
	deviceListContainer = container.NewVBox()

	// 刷新设备列表的函数
	refreshDeviceList := func() {
		devices, err := adb.GetDevices()
		deviceListContainer.RemoveAll()

		if err != nil {
			errorLabel := widget.NewLabel("获取设备失败: " + err.Error())
			errorLabel.TextStyle = fyne.TextStyle{Italic: true}
			deviceListContainer.Add(errorLabel)
			return
		}

		if len(devices) == 0 {
			noDeviceLabel := widget.NewLabel("没有检测到连接的设备")
			noDeviceLabel.TextStyle = fyne.TextStyle{Italic: true}
			deviceListContainer.Add(noDeviceLabel)
			return
		}

		for _, device := range devices {
			deviceCard := canvas.NewRectangle(color.White)
			deviceCard.CornerRadius = 10
			deviceCard.StrokeColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
			deviceCard.StrokeWidth = 1

			deviceIcon := widget.NewIcon(theme.ComputerIcon())
			deviceName := widget.NewLabel(device.Name)
			deviceName.TextStyle = fyne.TextStyle{Bold: true}
			deviceID := widget.NewLabel("ID: " + device.ID)
			deviceID.TextStyle = fyne.TextStyle{Italic: true}

			deviceInfo := container.NewVBox(deviceName, deviceID)
			deviceContent := container.NewBorder(nil, nil, deviceIcon, nil, deviceInfo)
			paddedContent := container.NewPadded(deviceContent)

			deviceContainer := container.NewStack(deviceCard, paddedContent)
			deviceListContainer.Add(deviceContainer)
		}
	}

	// 初始刷新设备列表
	refreshDeviceList()

	// 刷新按钮
	refreshBtn := widget.NewButtonWithIcon("刷新设备列表", theme.ViewRefreshIcon(), func() {
		refreshDeviceList()
	})

	// 操作按钮区域
	buttonsLabel := widget.NewLabel("设备操作:")
	buttonsLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 创建设备操作按钮
	deviceOperations := []struct {
		name string
		cmd  string
	}{
		{"重启设备", constants.ADBReboot},
		{"进入Recovery模式", constants.ADBRebootRecovery},
		{"进入Bootloader", constants.ADBRebootBootloader},
		{"获取设备序列号", constants.ADBGetSerialNo},
		{"获取设备状态", constants.ADBGetState},
	}

	buttonsContainer := container.NewVBox()
	for _, op := range deviceOperations {
		btn := widget.NewButton(op.name, func(cmd string) func() {
			return func() {
				// 创建执行命令的对话框
				messageCard := widget.NewCard("执行命令", cmd,
					widget.NewButton("确定", func() {}),
				)
				dialog := widget.NewModalPopUp(messageCard, ui.window.Canvas())
				dialog.Show()
			}
		}(op.cmd))
		buttonsContainer.Add(btn)
	}

	// 组合所有组件
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		devicesLabel,
		deviceListContainer,
		refreshBtn,
		widget.NewSeparator(),
		buttonsLabel,
		buttonsContainer,
	)

	return container.NewPadded(content)
}

func (ui *UI) createAppManagementPage() fyne.CanvasObject {
	return widget.NewLabel("应用管理功能正在开发中...")
}

func (ui *UI) createLogViewingPage() fyne.CanvasObject {
	return widget.NewLabel("日志查看功能正在开发中...")
}

func (ui *UI) createFileTransferPage() fyne.CanvasObject {
	return widget.NewLabel("文件传输功能正在开发中...")
}

func (ui *UI) createSettingsPage() fyne.CanvasObject {
	return widget.NewLabel("设置功能正在开发中...")
}

func (ui *UI) createFeatureDetailContent(config constants.FeatureConfig) fyne.CanvasObject {
	// 根据功能ID选择不同的页面实现
	switch config.ID {
	case "device_management":
		return ui.createDeviceManagementPage()
	case "app_management":
		return ui.createAppManagementPage()
	case "log_viewing":
		return ui.createLogViewingPage()
	case "file_transfer":
		return ui.createFileTransferPage()
	case "settings":
		return ui.createSettingsPage()
	}

	// 默认实现：显示命令列表
	description := widget.NewLabel(config.Description)
	description.Wrapping = fyne.TextWrapWord

	commandsTitle := widget.NewLabel("相关ADB命令:")
	commandsTitle.TextStyle = fyne.TextStyle{Bold: true}

	var commandItems []fyne.CanvasObject
	for _, cmdKey := range config.CommandGroup {
		// 获取实际命令字符串
		cmdStr, exists := constants.CommandMap[cmdKey]
		displayText := cmdKey // 默认显示键名
		if exists {
			displayText = cmdStr
		}
		cmdLabel := widget.NewLabel("• " + displayText)
		cmdLabel.Wrapping = fyne.TextWrapWord
		commandItems = append(commandItems, cmdLabel)
	}

	// 如果命令组为空，显示提示
	if len(commandItems) == 0 {
		noCommands := widget.NewLabel("暂无相关命令")
		noCommands.TextStyle = fyne.TextStyle{Italic: true}
		commandItems = append(commandItems, noCommands)
	}

	commandsContainer := container.NewVBox(commandItems...)

	// 创建分隔线
	separator := widget.NewSeparator()

	// 创建状态标签
	statusText := "功能ID: " + config.ID + " | 命令数量: " + strconv.Itoa(len(config.CommandGroup))
	statusLabel := widget.NewLabel(statusText)
	statusLabel.TextStyle = fyne.TextStyle{Italic: true}
	statusLabel.Alignment = fyne.TextAlignCenter

	// 组合所有组件
	content := container.NewVBox(
		description,
		separator,
		commandsTitle,
		commandsContainer,
		widget.NewSeparator(),
		statusLabel,
	)

	return container.NewPadded(content)
}

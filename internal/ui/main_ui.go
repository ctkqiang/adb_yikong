package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"
	"yikong/internal/adb"
	"yikong/internal/constants"
	"yikong/internal/logging"

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

// executeCommandDialog 显示命令执行对话框并执行ADB命令
func (ui *UI) executeCommandDialog(commandKey string, commandDisplay string, params map[string]string) {
	// 检查是否选择了设备（对于需要设备的命令）
	// 注意：有些命令不需要设备，如adb devices
	if ui.selectedDevice == "" && commandKey != "ADBDevices" && commandKey != "ADBDevicesL" {
		var dialog *widget.PopUp
		btn := widget.NewButton("确定", func() {
			if dialog != nil {
				dialog.Hide()
			}
		})
		errorCard := widget.NewCard("错误", "请先选择一个设备", btn)
		dialog = widget.NewModalPopUp(errorCard, ui.window.Canvas())
		dialog.Show()
		return
	}

	// 创建执行对话框
	title := "执行命令: " + commandDisplay
	outputDisplay := widget.NewMultiLineEntry()
	outputDisplay.Wrapping = fyne.TextWrapWord
	outputDisplay.SetPlaceHolder("命令输出将显示在这里...")
	outputDisplay.Disable()

	scrollContainer := container.NewScroll(outputDisplay)
	scrollContainer.SetMinSize(fyne.NewSize(500, 300))

	statusLabel := widget.NewLabel("准备执行...")
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	closeBtn := widget.NewButton("关闭", nil)
	closeBtn.Disable() // 初始禁用，执行完成后启用

	var dialog *widget.PopUp

	// 更新输出显示的函数
	updateOutput := func(text string) {
		fyne.Do(func() {
			currentText := outputDisplay.Text
			if currentText == "" {
				outputDisplay.SetText(text)
			} else {
				outputDisplay.SetText(currentText + "\n" + text)
			}
			scrollContainer.ScrollToBottom()
		})
	}

	// 执行命令
	go func() {
		fyne.Do(func() {
			statusLabel.SetText("正在执行命令...")
		})

		// 设置超时时间
		timeout := 30 * time.Second

		// 执行命令
		result, err := adb.ExecuteCommandFromConstants(
			commandKey,
			ui.selectedDevice, // 对于不需要设备的命令，adb.go会忽略空设备ID
			params,
			timeout,
			updateOutput,
		)

		fyne.Do(func() {
			if err != nil {
				statusLabel.SetText("执行失败")
				updateOutput("错误: " + err.Error())
			} else if result != nil {
				if result.Success {
					statusLabel.SetText("执行成功")
					updateOutput("命令执行完成，退出码: " + strconv.Itoa(result.ExitCode))
				} else {
					statusLabel.SetText("执行失败")
					updateOutput("命令执行失败，退出码: " + strconv.Itoa(result.ExitCode))
					if result.ErrorOutput != "" {
						updateOutput("错误输出: " + result.ErrorOutput)
					}
				}
				updateOutput("执行时间: " + result.Duration.String())
			}
			closeBtn.Enable()
		})
	}()

	// 设置关闭按钮回调
	closeBtn.OnTapped = func() {
		if dialog != nil {
			dialog.Hide()
		}
	}

	// 创建对话框内容
	content := container.NewVBox(
		statusLabel,
		widget.NewSeparator(),
		scrollContainer,
		widget.NewSeparator(),
		closeBtn,
	)

	card := widget.NewCard(title, "", content)
	dialog = widget.NewModalPopUp(card, ui.window.Canvas())
	dialog.Show()
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
		name    string
		cmdKey  string
		display string
	}{
		{"列出设备", "ADBDevices", "adb devices"},
		{"列出设备(长格式)", "ADBDevicesL", "adb devices -l"},
		{"重启设备", "ADBReboot", "adb reboot"},
		{"进入Recovery模式", "ADBRebootRecovery", "adb reboot recovery"},
		{"进入Bootloader", "ADBRebootBootloader", "adb reboot-bootloader"},
		{"获取设备序列号", "ADBGetSerialNo", "adb get-serialno"},
		{"获取设备状态", "ADBGetState", "adb get-state"},
		{"以root权限运行", "ADBRoot", "adb root"},
	}

	buttonsContainer := container.NewVBox()
	for _, op := range deviceOperations {
		btn := widget.NewButton(op.name, func(cmdKey string, display string) func() {
			return func() {
				ui.executeCommandDialog(cmdKey, display, nil)
			}
		}(op.cmdKey, op.display))
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
	// 创建标题
	title := widget.NewLabel("日志查看")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// 创建设备信息标签
	var deviceInfoText string
	if ui.selectedDevice == "" {
		deviceInfoText = "未选择设备"
	} else {
		deviceInfoText = "当前设备: " + ui.selectedDevice
	}
	deviceInfo := widget.NewLabel(deviceInfoText)
	deviceInfo.TextStyle = fyne.TextStyle{Italic: true}
	deviceInfo.Alignment = fyne.TextAlignCenter

	// 创建状态标签
	statusLabel := widget.NewLabel("")
	statusLabel.TextStyle = fyne.TextStyle{Monospace: true}
	statusLabel.Alignment = fyne.TextAlignCenter
	statusLabel.Wrapping = fyne.TextWrapWord

	// 创建日志显示区域
	logDisplay := widget.NewMultiLineEntry()
	logDisplay.Wrapping = fyne.TextWrapWord
	logDisplay.SetPlaceHolder("日志将显示在这里...")
	logDisplay.Disable() // 设置为只读

	// 创建滚动容器
	scrollContainer := container.NewScroll(logDisplay)
	scrollContainer.SetMinSize(fyne.NewSize(600, 400))

	// 声明按钮变量，以便在闭包中引用
	var fetchLogBtn, clearLogBtn, saveLogBtn *widget.Button

	fetchLogBtn = widget.NewButton("获取日志", func() {
		if ui.selectedDevice == "" {
			var dialog *widget.PopUp
			btn := widget.NewButton("确定", func() {
				if dialog != nil {
					dialog.Hide()
				}
			})
			errorCard := widget.NewCard("错误", "请先选择一个设备", btn)
			dialog = widget.NewModalPopUp(errorCard, ui.window.Canvas())
			dialog.Show()
			return
		}

		// 显示加载中
		logDisplay.SetText("正在获取日志...")
		statusLabel.SetText("状态: 正在获取日志...")
		fetchLogBtn.Disable()

		go func() {
			logs, err := adb.GetLogcat(ui.selectedDevice)
			fyne.Do(func() {
				fetchLogBtn.Enable()
				if err != nil {
					logging.Error("获取日志失败: deviceID=%s, 错误: %v", ui.selectedDevice, err)
					var dialog *widget.PopUp
					btn := widget.NewButton("确定", func() {
						if dialog != nil {
							dialog.Hide()
						}
					})
					errorCard := widget.NewCard("错误", "获取日志失败: "+err.Error(), btn)
					dialog = widget.NewModalPopUp(errorCard, ui.window.Canvas())
					dialog.Show()
					logDisplay.SetText("")
					return
				}
				logging.Info("获取到日志，长度: %d 字节", len(logs))

				// 更新状态标签
				var logSizeStr string
				if len(logs) >= 1024*1024 {
					logSizeStr = fmt.Sprintf("%.1f MB", float64(len(logs))/(1024*1024))
				} else if len(logs) >= 1024 {
					logSizeStr = fmt.Sprintf("%.1f KB", float64(len(logs))/1024)
				} else {
					logSizeStr = fmt.Sprintf("%d 字节", len(logs))
				}
				statusLabel.SetText(fmt.Sprintf("状态: 获取到日志，长度: %s", logSizeStr))

				// 处理大日志文件
				const maxDisplaySize = 10 * 1024 * 1024 // 10MB
				const headSize = 1 * 1024 * 1024        // 显示开头1MB
				const tailSize = 4 * 1024 * 1024        // 显示结尾4MB

				var displayText string
				if len(logs) > maxDisplaySize {
					logging.Warn("日志过大 (%d 字节)，进行智能截断显示", len(logs))
					// 更新状态标签
					statusLabel.SetText(fmt.Sprintf("状态: 日志过大 (%s)，已进行智能截断显示", logSizeStr))

					// 格式化字节大小
					var sizeStr string
					if len(logs) >= 1024*1024 {
						sizeStr = fmt.Sprintf("%.1f MB", float64(len(logs))/(1024*1024))
					} else if len(logs) >= 1024 {
						sizeStr = fmt.Sprintf("%.1f KB", float64(len(logs))/1024)
					} else {
						sizeStr = fmt.Sprintf("%d 字节", len(logs))
					}

					// 智能截断：显示开头和结尾部分
					headPart := logs[:headSize]
					tailPart := logs[len(logs)-tailSize:]
					// 格式化headSize和tailSize用于显示
					var headSizeStr, tailSizeStr string
					if headSize >= 1024*1024 {
						headSizeStr = fmt.Sprintf("%d MB", headSize/(1024*1024))
					} else if headSize >= 1024 {
						headSizeStr = fmt.Sprintf("%d KB", headSize/1024)
					} else {
						headSizeStr = fmt.Sprintf("%d 字节", headSize)
					}
					if tailSize >= 1024*1024 {
						tailSizeStr = fmt.Sprintf("%d MB", tailSize/(1024*1024))
					} else if tailSize >= 1024 {
						tailSizeStr = fmt.Sprintf("%d KB", tailSize/1024)
					} else {
						tailSizeStr = fmt.Sprintf("%d 字节", tailSize)
					}

					displayText = headPart + "\n\n... [日志过长，已智能截断。显示开头" + headSizeStr + "和结尾" + tailSizeStr + "，完整日志大小: " + sizeStr + "] ...\n\n" + tailPart
				} else {
					displayText = logs
				}

				logDisplay.SetText(displayText)
				logging.Info("日志已设置到UI显示组件，显示长度: %d 字节", len(displayText))

				// 更新状态标签显示最终状态
				var displaySizeStr string
				if len(displayText) >= 1024*1024 {
					displaySizeStr = fmt.Sprintf("%.1f MB", float64(len(displayText))/(1024*1024))
				} else if len(displayText) >= 1024 {
					displaySizeStr = fmt.Sprintf("%.1f KB", float64(len(displayText))/1024)
				} else {
					displaySizeStr = fmt.Sprintf("%d 字节", len(displayText))
				}
				currentStatus := statusLabel.Text
				if len(logs) > maxDisplaySize {
					statusLabel.SetText(fmt.Sprintf("%s | 显示长度: %s (已截断)", currentStatus, displaySizeStr))
				} else {
					statusLabel.SetText(fmt.Sprintf("%s | 显示长度: %s", currentStatus, displaySizeStr))
				}

				scrollContainer.ScrollToBottom()
				logging.Info("已滚动到底部")
			})
		}()
	})
	fetchLogBtn.Importance = widget.HighImportance

	clearLogBtn = widget.NewButton("清除日志", func() {
		if ui.selectedDevice == "" {
			var dialog *widget.PopUp
			btn := widget.NewButton("确定", func() {
				if dialog != nil {
					dialog.Hide()
				}
			})
			errorCard := widget.NewCard("错误", "请先选择一个设备", btn)
			dialog = widget.NewModalPopUp(errorCard, ui.window.Canvas())
			dialog.Show()
			return
		}

		clearLogBtn.Disable()
		go func() {
			err := adb.ClearLogcat(ui.selectedDevice)
			fyne.Do(func() {
				clearLogBtn.Enable()
				if err != nil {
					var dialog *widget.PopUp
					btn := widget.NewButton("确定", func() {
						if dialog != nil {
							dialog.Hide()
						}
					})
					errorCard := widget.NewCard("错误", "清除日志失败: "+err.Error(), btn)
					dialog = widget.NewModalPopUp(errorCard, ui.window.Canvas())
					dialog.Show()
					return
				}
				// 显示成功消息
				logDisplay.SetText("日志缓冲区已清除")
			})
		}()
	})
	clearLogBtn.Importance = widget.MediumImportance

	saveLogBtn = widget.NewButton("保存日志", func() {
		// 保存日志功能暂未实现
		infoCard := widget.NewCard("提示", "保存日志功能正在开发中",
			widget.NewButton("确定", func() {

			}),
		)
		dialog := widget.NewModalPopUp(infoCard, ui.window.Canvas())
		dialog.Show()
	})
	saveLogBtn.Importance = widget.MediumImportance

	// 按钮容器
	buttonContainer := container.NewHBox(
		fetchLogBtn,
		clearLogBtn,
		saveLogBtn,
	)

	// 组合所有组件
	content := container.NewVBox(
		title,
		deviceInfo,
		statusLabel,
		widget.NewSeparator(),
		scrollContainer,
		buttonContainer,
	)

	return container.NewPadded(content)
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

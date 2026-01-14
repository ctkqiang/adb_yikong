package ui

import (
	"yikong/internal/adb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	// 第一页：设备列表
	devicePage := ui.createDevicePage()

	// 第二页：功能页面（初始隐藏）
	functionPage := ui.createFunctionPage()
	functionPage.Hide()

	ui.mainContainer = container.NewStack(devicePage, functionPage)

	return ui.mainContainer
}

func (ui *UI) createDevicePage() fyne.CanvasObject {
	title := widget.NewLabel("已连接设备")
	title.TextStyle = fyne.TextStyle{Bold: true}

	// 刷新按钮
	refreshBtn := widget.NewButton("刷新设备", func() {
		ui.refreshDevices()
	})

	// 创建设备列表
	ui.deviceList = widget.NewList(
		func() int {
			return len(ui.devices)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			btn := o.(*widget.Button)
			btn.SetText(ui.devices[i].Name)
			btn.OnTapped = func() {
				ui.selectDevice(ui.devices[i].ID)
			}
		},
	)

	// 初始加载设备
	ui.refreshDevices()

	return container.NewBorder(
		title,
		refreshBtn,
		nil, nil,
		ui.deviceList,
	)
}

func (ui *UI) createFunctionPage() fyne.CanvasObject {
	backBtn := widget.NewButton("返回设备列表", func() {
		ui.showDevicePage()
	})

	deviceInfo := widget.NewLabel("")

	// 功能按钮区域
	functionBtns := container.NewVBox(
		widget.NewButton("设备管理", func() { deviceInfo.SetText("设备管理功能") }),
		widget.NewButton("应用管理", func() { deviceInfo.SetText("应用管理功能") }),
		widget.NewButton("日志查看", func() { deviceInfo.SetText("日志查看功能") }),
		widget.NewButton("文件传输", func() { deviceInfo.SetText("文件传输功能") }),
		widget.NewButton("设置", func() { deviceInfo.SetText("设置功能") }),
	)

	return container.NewBorder(
		backBtn,
		nil, nil, nil,
		container.NewVBox(
			deviceInfo,
			functionBtns,
		),
	)
}

func (ui *UI) refreshDevices() {
	devices, err := adb.GetDevices()
	if err != nil {
		ui.window.Content().(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("获取设备失败: " + err.Error())
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
	for _, obj := range ui.mainContainer.Objects {
		obj.Hide()
	}
	ui.mainContainer.Objects[0].Show()
	ui.mainContainer.Refresh()
}

func (ui *UI) showFunctionPage() {
	for _, obj := range ui.mainContainer.Objects {
		obj.Hide()
	}
	ui.mainContainer.Objects[1].Show()
	ui.mainContainer.Refresh()
}

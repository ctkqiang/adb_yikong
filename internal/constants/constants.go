package constants

// 应用常量
const (
	AppName      = "ADB易控"
	Version      = "1.0.0"
	Author       = "钟智强"
	AuthorEmail  = "johnmelodymel@qq.com"
	AuthorGithub = "https://github.com/ctkqiang"

	WINDOW_ADB_URL = "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"
	UNIX_ADB_URL   = ""
)

// ADB基础命令
const (
	ADB        = "adb"
	ADBHelp    = "adb help"
	ADBVersion = "adb version"
)

// ADB服务器命令
const (
	ADBKillServer  = "adb kill-server"
	ADBStartServer = "adb start-server"
)

// 设备管理
const (
	ADBDevices     = "adb devices"
	ADBDevicesL    = "adb devices -l" // 详细设备信息
	ADBConnect     = "adb connect %s" // 需要IP地址参数
	ADBUSB         = "adb usb"
	ADBGetSerialNo = "adb get-serialno"
	ADBGetState    = "adb get-state"
)

// 重启命令
const (
	ADBReboot           = "adb reboot"
	ADBRebootRecovery   = "adb reboot recovery"
	ADBRebootBootloader = "adb reboot-bootloader"
	ADBRoot             = "adb root"
)

// Shell命令
const (
	ADBShell = "adb shell"
	// Shell子命令
	ShellGetProp        = "adb shell getprop %s" // 需要属性名参数
	ShellAndroidVersion = "adb shell getprop ro.build.version.release"
	ShellPWD            = "adb shell pwd"
	ShellLS             = "adb shell ls"
	ShellLSRecursive    = "adb shell ls -R"
)

// 文件操作
const (
	ADBPush = "adb push %s %s" // 需要源文件和目标路径
	ADBPull = "adb pull %s %s" // 需要设备路径和本地路径
)

// 应用管理
const (
	ADBInstall    = "adb install %s"      // 需要APK路径
	ADBInstallR   = "adb install -r %s"   // 重新安装保留数据
	ADBInstallK   = "adb install -k %s"   // 保留数据安装
	ADBUninstall  = "adb uninstall %s"    // 需要包名
	ADBUninstallK = "adb uninstall -k %s" // 卸载保留数据
)

// 包管理器命令
const (
	PMUninstall        = "adb shell pm uninstall %s" // 需要包名
	PMClear            = "adb shell pm clear %s"     // 需要包名
	PMListPackages     = "adb shell pm list packages"
	PMListPackages3    = "adb shell pm list packages -3" // 第三方应用
	PMListPackagesS    = "adb shell pm list packages -s" // 系统应用
	PMGrant            = "adb shell pm grant %s %s"      // 包名和权限
	PMRevoke           = "adb shell pm revoke %s %s"     // 包名和权限
	PMResetPermissions = "adb shell pm reset-permissions"
)

// 日志命令
const (
	ADBLogcat    = "adb logcat"
	ADBLogcatC   = "adb logcat -c"      // 清除日志
	ADBLogcatD   = "adb logcat -d > %s" // 保存到文件
	ADBBugreport = "adb bugreport > %s" // 保存到文件
)

// 屏幕操作
const (
	Screenshot          = "adb shell screencap -p %s" // 需要保存路径
	Screenrecord        = "adb shell screenrecord %s" // 需要保存路径
	ScreenrecordVerbose = "adb shell screenrecord --verbose %s"
)

// 输入和按键
const (
	InputText     = "adb shell input text '%s'"            // 需要文本内容
	InputKeyEvent = "adb shell input keyevent %d"          // 需要按键码
	InputTap      = "adb shell input tap %d %d"            // 需要x,y坐标
	InputSwipe    = "adb shell input swipe %d %d %d %d %d" // 需要起始坐标和时长
)

// Activity管理器
const (
	AMStart     = "adb shell am start %s" // 需要Intent参数
	AMStartHome = "adb shell am start -W -c android.intent.category.HOME -a android.intent.action.MAIN"
	AMStartView = "adb shell am start -a android.intent.action.VIEW"
	AMStartCall = "adb shell am start -a android.intent.action.CALL -d tel:%s" // 需要电话号码
	AMStartSMS  = "adb shell am start -a android.intent.action.SENDTO -d sms:%s --es sms_body \"%s\" --ez exit_on_sent false"
	AMBroadcast = "adb shell am broadcast -a '%s'" // 需要action名称
)

// 窗口管理器
const (
	WMSize         = "adb shell wm size %dx%d" // 需要宽高
	WMSizeReset    = "adb shell wm size reset"
	WMDensity      = "adb shell wm density %d" // 需要密度值
	WMDensityReset = "adb shell wm density reset"
)

// 备份和恢复
const (
	ADBBackup         = "adb backup -apk -all -f %s" // 需要备份文件路径
	ADBBackupShared   = "adb backup -apk -shared -all -f %s"
	ADBBackupNoSystem = "adb backup -apk -nosystem -all -f %s"
	ADBRestore        = "adb restore %s"  // 需要备份文件路径
	ADBSideload       = "adb sideload %s" // 需要ROM/zip文件路径
)

// 设备信息
const (
	DumpsysBattery          = "adb shell dumpsys battery"
	DumpsysBatterySetLevel  = "adb shell dumpsys battery set level %d"
	DumpsysBatterySetStatus = "adb shell dumpsys battery set status %d"
	DumpsysBatteryReset     = "adb shell dumpsys battery reset"
	DumpsysBatterySetUSB    = "adb shell dumpsys battery set usb %s"
	DumpsysIPhoneSubInfo    = "adb shell dumpsys iphonesubinfo" // 获取IMEI
	DumpsysWindow           = "adb shell dumpsys window windows"
	DumpsysPackage          = "adb shell dumpsys package packages"
	DumpsysActivity         = "adb shell dumpsys activity %s/%s" // 包名/Activity名
	Netstat                 = "adb shell netstat"
	PS                      = "adb shell ps"
)

// Monkey测试
const (
	MonkeyTest = "adb shell monkey -p %s -v %d -s %d" // 包名、事件数、种子
)

// Shared Preferences
const (
	SPPut    = "adb shell 'am broadcast -a %s.sp.PUT --es key %s --es value \"%s\"'"
	SPRemove = "adb shell 'am broadcast -a %s.sp.REMOVE --es key %s'"
	SPClear  = "adb shell 'am broadcast -a %s.sp.CLEAR --es key %s'"
)

// Fastboot命令
const (
	FastbootDevices = "fastboot devices"
)

// 常用路径常量
const (
	PathDataData    = "/data/data/%s"               // 应用数据目录
	PathSharedPrefs = "/data/data/%s/shared_prefs/" // SharedPreferences目录
	PathDatabases   = "/data/data/%s/databases/"    // 数据库目录
	PathDataApp     = "/data/app/"                  // 用户安装的APK
	PathSystemApp   = "/system/app/"                // 系统APK
	PathSDCard      = "/sdcard/"                    // SD卡路径
	PathExternalSD  = "/mnt/sdcard/external_sd/"    // 外部SD卡
)

// 按键码常量
const (
	KeyCodeHome           = 3
	KeyCodeBack           = 4
	KeyCodeCall           = 5
	KeyCodeEndCall        = 6
	KeyCodePower          = 26
	KeyCodeCamera         = 27
	KeyCodeClear          = 28
	KeyCodeEnter          = 66
	KeyCodeDelete         = 67
	KeyCodeVolumeUp       = 24
	KeyCodeVolumeDown     = 25
	KeyCodeMenu           = 82
	KeyCodeSearch         = 84
	KeyCodeMediaPlayPause = 85
	KeyCodeMediaStop      = 86
	KeyCodeMediaNext      = 87
	KeyCodeMediaPrev      = 88
	KeyCodeMute           = 91
	KeyCodePageUp         = 92
	KeyCodePageDown       = 93
	KeyCodeExplorer       = 64 // 打开浏览器
	KeyCodeEnvelope       = 65 // 打开邮件
)

// 按键码映射表（用于显示）
var KeyCodeMap = map[int]string{
	0:  "KEYCODE_0",
	1:  "KEYCODE_SOFT_LEFT",
	2:  "KEYCODE_SOFT_RIGHT",
	3:  "KEYCODE_HOME",
	4:  "KEYCODE_BACK",
	5:  "KEYCODE_CALL",
	6:  "KEYCODE_ENDCALL",
	26: "KEYCODE_POWER",
	27: "KEYCODE_CAMERA",
	64: "KEYCODE_EXPLORER",
	66: "KEYCODE_ENTER",
	67: "KEYCODE_DEL",
	82: "KEYCODE_MENU",
	84: "KEYCODE_SEARCH",
	85: "KEYCODE_MEDIA_PLAY_PAUSE",
	86: "KEYCODE_MEDIA_STOP",
	87: "KEYCODE_MEDIA_NEXT",
	88: "KEYCODE_MEDIA_PREVIOUS",
	91: "KEYCODE_MUTE",
}

// 常用快捷键组合
const (
	KeyCombinationScreenshot = "Volume Down + Power"
	KeyCombinationPowerMenu  = "Power (长按)"
	KeyCombinationRecentApps = "Home (双击) 或 Recent Apps 键"
)

// 功能配置结构体
type FeatureConfig struct {
	ID           string   // 功能唯一标识
	Name         string   // 显示名称
	Description  string   // 功能描述
	IconName     string   // 图标名称（对应fyne主题图标）
	CommandGroup []string // 相关的ADB命令常量组
	DefaultLabel string   // 默认显示的标签文本
}

// 功能映射表（用于UI功能按钮与ADB命令的映射）
var FeatureMap = map[string]FeatureConfig{
	"device_management": {
		ID:           "device_management",
		Name:         "设备管理",
		Description:  "查看和管理已连接的Android设备",
		IconName:     "SettingsIcon",
		CommandGroup: []string{"ADBDevices", "ADBDevicesL", "ADBConnect", "ADBUSB", "ADBGetSerialNo", "ADBGetState", "ADBReboot", "ADBRebootRecovery", "ADBRebootBootloader", "ADBRoot"},
		DefaultLabel: "设备管理功能",
	},
	"app_management": {
		ID:           "app_management",
		Name:         "应用管理",
		Description:  "安装、卸载和管理Android应用",
		IconName:     "DocumentSaveIcon",
		CommandGroup: []string{"ADBInstall", "ADBInstallR", "ADBInstallK", "ADBUninstall", "ADBUninstallK", "PMListPackages", "PMListPackages3", "PMListPackagesS", "PMUninstall", "PMClear", "PMGrant", "PMRevoke", "PMResetPermissions"},
		DefaultLabel: "应用管理功能",
	},
	"log_viewing": {
		ID:           "log_viewing",
		Name:         "日志查看",
		Description:  "查看和保存设备日志",
		IconName:     "DocumentPrintIcon",
		CommandGroup: []string{"ADBLogcat", "ADBLogcatC", "ADBLogcatD", "ADBBugreport"},
		DefaultLabel: "日志查看功能",
	},
	"file_transfer": {
		ID:           "file_transfer",
		Name:         "文件传输",
		Description:  "在设备和电脑之间传输文件",
		IconName:     "MailSendIcon",
		CommandGroup: []string{"ADBPush", "ADBPull", "ShellLS", "ShellLSRecursive", "ShellPWD"},
		DefaultLabel: "文件传输功能",
	},
	"settings": {
		ID:           "settings",
		Name:         "设置",
		Description:  "应用程序设置和ADB配置",
		IconName:     "SettingsIcon",
		CommandGroup: []string{"ADBKillServer", "ADBStartServer", "ADBHelp", "ADBVersion"},
		DefaultLabel: "设置功能",
	},
}

// 命令映射表（ADB命令常量名到实际命令字符串的映射）
var CommandMap = map[string]string{
	// 设备管理命令
	"ADBDevices":          ADBDevices,
	"ADBDevicesL":         ADBDevicesL,
	"ADBConnect":          ADBConnect,
	"ADBUSB":              ADBUSB,
	"ADBGetSerialNo":      ADBGetSerialNo,
	"ADBGetState":         ADBGetState,
	"ADBReboot":           ADBReboot,
	"ADBRebootRecovery":   ADBRebootRecovery,
	"ADBRebootBootloader": ADBRebootBootloader,
	"ADBRoot":             ADBRoot,
	// 应用管理命令
	"ADBInstall":         ADBInstall,
	"ADBInstallR":        ADBInstallR,
	"ADBInstallK":        ADBInstallK,
	"ADBUninstall":       ADBUninstall,
	"ADBUninstallK":      ADBUninstallK,
	"PMListPackages":     PMListPackages,
	"PMListPackages3":    PMListPackages3,
	"PMListPackagesS":    PMListPackagesS,
	"PMUninstall":        PMUninstall,
	"PMClear":            PMClear,
	"PMGrant":            PMGrant,
	"PMRevoke":           PMRevoke,
	"PMResetPermissions": PMResetPermissions,
	// 日志查看命令
	"ADBLogcat":    ADBLogcat,
	"ADBLogcatC":   ADBLogcatC,
	"ADBLogcatD":   ADBLogcatD,
	"ADBBugreport": ADBBugreport,
	// 文件传输命令
	"ADBPush":          ADBPush,
	"ADBPull":          ADBPull,
	"ShellLS":          ShellLS,
	"ShellLSRecursive": ShellLSRecursive,
	"ShellPWD":         ShellPWD,
	// 设置命令
	"ADBKillServer":  ADBKillServer,
	"ADBStartServer": ADBStartServer,
	"ADBHelp":        ADBHelp,
	"ADBVersion":     ADBVersion,
}

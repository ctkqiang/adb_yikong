package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"yikong/internal/http"
	"yikong/internal/logging"
	"yikong/internal/utilities"
)

// InspectIsADBExisted 检查系统中是否存在ADB工具
// 该函数会通过多种方式查找ADB，包括PATH环境变量、常见安装路径等
//
// 返回值:
//   - bool: 如果找到ADB则返回true，否则返回false
func InspectIsADBExisted() bool {
	cmd := exec.Command("adb", "version")

	if err := cmd.Run(); err == nil {
		logging.Info("已通过PATH环境变量找到ADB")
		return true
	}

	// 遍历PATH环境变量中的所有目录查找adb可执行文件
	pathDirs := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))

	for _, dir := range pathDirs {
		adbPath := filepath.Join(dir, "adb")
		if runtime.GOOS == "windows" {
			adbPath += ".exe"
		}

		if _, err := os.Stat(adbPath); err == nil {
			return true
		}
	}

	// 在常见的ADB安装路径中查找
	commonPaths := []string{
		"C:\\platform-tools\\adb.exe",
		"/usr/local/bin/adb",
		"/opt/homebrew/bin/adb",
		filepath.Join(os.Getenv("HOME"), "platform-tools", "adb"),
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// SetupADB 根据当前操作系统设置ADB（Android Debug Bridge）环境
// 该函数会自动检测操作系统类型，并在不同平台上安装ADB工具
// 对于Windows系统，会下载并解压ADB工具包，然后将其添加到系统PATH环境变量
// 对于macOS系统，使用Homebrew进行安装
// 对于Linux系统，调用专门的安装函数
//
// 返回值:
//
//	error: 安装过程中出现错误时返回错误信息，成功时返回nil
func SetupADB() error {
	operating_system := runtime.GOOS

	logging.Info("操作系统: %s", operating_system)

	switch operating_system {
	case "windows":
		// Windows平台ADB安装流程：下载ZIP包 -> 解压到指定目录 -> 添加到系统PATH -> 清理临时文件
		win_url := "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"
		tempDir := os.Getenv("TEMP")
		zipPath := filepath.Join(tempDir, "adb.zip")
		extractPath := "C:\\platform-tools"

		logging.Debug("正在从 %s 下载...", win_url)

		if err := http.DownloadFile(win_url, zipPath); err != nil {
			logging.Error("下载失败: %v", err)
			os.Exit(1)
		}

		logging.Info("正在解压到 %s...\n", extractPath)

		if err := utilities.ExtractZipWindows(zipPath, extractPath); err != nil {
			logging.Error("解压失败: %v", err)
			os.Exit(1)
		}

		logging.Info("正在添加到系统 PATH 环境变量...")

		if err := utilities.Add_TO_WINDOW_PATH(extractPath); err != nil {
			logging.Error("添加到 PATH 环境变量失败: %v", err)
			os.Exit(1)
		}

		os.Remove(zipPath)
		logging.Info("ADB 安装成功！")
		logging.Info("请重启终端以使 PATH 更改生效。")

	case "darwin":
		utilities.InstallViaHomebrew()
	case "linux":
		utilities.InstallADBLinux()
	default:
		logging.Error("不支持的操作系统: %s", operating_system)
		os.Exit(1)
	}

	return nil
}

package config

import (
	"os"
	"path/filepath"
	"runtime"
	"yikong/internal/http"
	"yikong/internal/logging"
	"yikong/internal/utilities"
)

func SetupADB() error {
	operating_system := runtime.GOOS

	logging.Info("操作系统:", operating_system)

	switch operating_system {
	case "windows":
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

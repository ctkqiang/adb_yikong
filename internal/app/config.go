package config

import (
	"os"
	"path/filepath"
	"runtime"
	"yikong/internal/http"
	"yikong/internal/logging"
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

	}

	return nil
}

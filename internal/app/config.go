package config

import (
	"os"
	"path/filepath"
	"runtime"
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

		logging.Debug("Downloading from %s...", win_url)

	}

	return nil
}

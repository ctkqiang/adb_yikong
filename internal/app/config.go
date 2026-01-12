package config

import (
	"runtime"
	"yikong/internal/logging"
)

func SetupADB() error {
	operating_system := runtime.GOOS

	logging.Info("操作系统:", operating_system)

	return nil
}

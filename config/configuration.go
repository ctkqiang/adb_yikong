package config

import (
	"fmt"
	"runtime"
)

func SetupADB() error {
	operating_system := runtime.GOOS

	fmt.Println("操作系统:", operating_system)

	return nil
}

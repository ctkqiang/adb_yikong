package logging

import (
	"fmt"
	config "yikong/internal/constants"
)

func Info(format string, a ...any) (int, error) {
	message := fmt.Sprintf(format, a...)
	return fmt.Println("\033[32m" + config.AppName + " [信息]" + message + "\033[0m")
}
func Debug(format string, a ...any) (int, error) {
	message := fmt.Sprintf(format, a...)
	return fmt.Println("\033[33m" + config.AppName + " [调试]" + message + "\033[0m")
}

func Error(format string, a ...any) (int, error) {
	message := fmt.Sprintf(format, a...)
	return fmt.Println("\033[31m" + config.AppName + " [错误]" + message + "\033[0m")
}

func Warn(format string, a ...any) (int, error) {
	message := fmt.Sprintf(format, a...)
	return fmt.Println("\033[33m" + config.AppName + " [警告]" + message + "\033[0m")
}

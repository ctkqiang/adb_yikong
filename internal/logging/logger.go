package logging

import (
	"fmt"
	"strings"
	config "yikong/internal/constants"
)

func Info(a ...any) (int, error) {
	return fmt.Println("\033[32m" + config.AppName + " [信息]" + strings.Join(convertToStrings(a), "") + "\033[0m")
}
func Debug(a ...any) (int, error) {
	return fmt.Println("\033[33m" + config.AppName + " [调试]" + strings.Join(convertToStrings(a), "") + "\033[0m")
}

func Error(a ...any) (int, error) {
	return fmt.Println("\033[31m" + config.AppName + " [错误]" + strings.Join(convertToStrings(a), "") + "\033[0m")
}

func Warn(a ...any) (int, error) {
	return fmt.Println("\033[33m" + config.AppName + " [警告]" + strings.Join(convertToStrings(a), "") + "\033[0m")
}

func convertToStrings(a []any) []string {
	s := make([]string, len(a))

	for i, v := range a {
		s[i] = fmt.Sprint(v)
	}

	return s
}

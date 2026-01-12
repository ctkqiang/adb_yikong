package logging

import (
	"fmt"
	"strings"
	config "yikong/internal/constants"
)

func Info(a ...any) (int, error) {
	return fmt.Println(config.AppName + " [信息]" + strings.Join(convertToStrings(a), ""))
}
func Debug(a ...any) (int, error) {
	return fmt.Println(config.AppName + " [调试]" + strings.Join(convertToStrings(a), ""))
}

func Error(a ...any) (int, error) {
	return 0, fmt.Errorf(config.AppName+" [错误]%s", strings.Join(convertToStrings(a), ""))
}

func Warn(a ...any) (int, error) {
	return fmt.Println(config.AppName + " [警告]" + strings.Join(convertToStrings(a), ""))
}

func convertToStrings(a []any) []string {
	s := make([]string, len(a))
	for i, v := range a {
		s[i] = fmt.Sprint(v)
	}
	return s
}

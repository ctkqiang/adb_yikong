package adb

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"yikong/internal/constants"
)

type Device struct {
	ID   string
	Name string
}

type CommandResult struct {
	Output      string
	ErrorOutput string
	ExitCode    int
	Duration    time.Duration
	Success     bool
}

func GetDevices() ([]Device, error) {
	cmd := exec.Command("adb", "devices")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out.String(), "\n")
	var devices []Device

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "*") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == "device" {
			devices = append(devices, Device{
				ID:   parts[0],
				Name: parts[0],
			})
		}
	}

	return devices, nil
}

// ExecuteCommand 执行任意命令并返回结果
func ExecuteCommand(command string, args []string, timeout time.Duration) (*CommandResult, error) {
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	result := &CommandResult{
		Output:      stdout.String(),
		ErrorOutput: stderr.String(),
		Duration:    duration,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Success = false
		log.Printf("命令执行失败: %s %v, 错误: %v, 退出码: %d, 执行时间: %v", command, args, err, result.ExitCode, duration)
		return result, err
	}

	result.ExitCode = 0
	result.Success = true
	log.Printf("命令执行成功: %s %v, 执行时间: %v", command, args, duration)
	return result, nil
}

// ExecuteADBCommand 执行ADB命令
func ExecuteADBCommand(args []string, timeout time.Duration) (*CommandResult, error) {
	return ExecuteCommand("adb", args, timeout)
}

// ExecuteADBCommandWithDevice 对特定设备执行ADB命令
func ExecuteADBCommandWithDevice(deviceID string, args []string, timeout time.Duration) (*CommandResult, error) {
	deviceArgs := []string{"-s", deviceID}
	deviceArgs = append(deviceArgs, args...)
	return ExecuteADBCommand(deviceArgs, timeout)
}

// ExecuteADBCommandString 执行字符串形式的ADB命令（方便使用）
func ExecuteADBCommandString(command string, timeout time.Duration) (*CommandResult, error) {
	// 简单的命令解析，将字符串分割为参数
	args := strings.Fields(command)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	// 移除命令名称"adb"如果存在
	if args[0] == "adb" {
		args = args[1:]
	}

	return ExecuteADBCommand(args, timeout)
}

// GetDeviceName 获取设备的友好名称（尝试获取设备型号）
func GetDeviceName(deviceID string) (string, error) {
	result, err := ExecuteADBCommandWithDevice(deviceID, []string{"shell", "getprop", "ro.product.model"}, 10*time.Second)
	if err != nil {
		return deviceID, err // 返回设备ID作为后备名称
	}

	name := strings.TrimSpace(result.Output)
	if name == "" {
		return deviceID, nil
	}
	return name, nil
}

// GetLogcat 获取设备日志（一次性转储）
func GetLogcat(deviceID string) (string, error) {
	log.Printf("获取设备日志: deviceID=%s", deviceID)
	result, err := ExecuteADBCommandWithDevice(deviceID, []string{"logcat", "-d"}, 30*time.Second)
	if err != nil {
		log.Printf("获取设备日志失败: deviceID=%s, 错误: %v", deviceID, err)
		return "", err
	}
	log.Printf("获取设备日志成功: deviceID=%s, 输出长度: %d", deviceID, len(result.Output))
	return result.Output, nil
}

// ClearLogcat 清除设备日志缓冲区
func ClearLogcat(deviceID string) error {
	log.Printf("清除设备日志: deviceID=%s", deviceID)
	_, err := ExecuteADBCommandWithDevice(deviceID, []string{"logcat", "-c"}, 10*time.Second)
	if err != nil {
		log.Printf("清除设备日志失败: deviceID=%s, 错误: %v", deviceID, err)
		return err
	}
	log.Printf("清除设备日志成功: deviceID=%s", deviceID)
	return nil
}

// ExecuteCommandStream 执行命令并实时输出流
func ExecuteCommandStream(command string, args []string, timeout time.Duration, outputCallback func(line string)) (*CommandResult, error) {
	log.Printf("开始执行命令: %s %v (超时: %v)", command, args, timeout)
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 创建管道用于实时读取
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// 实时读取输出
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if outputCallback != nil {
				outputCallback(line)
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if outputCallback != nil {
				outputCallback(line)
			}
		}
	}()

	// 等待命令完成
	startTime := time.Now()
	err = cmd.Wait()
	duration := time.Since(startTime)

	result := &CommandResult{
		Output:      stdout.String(),
		ErrorOutput: stderr.String(),
		Duration:    duration,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Success = false
		return result, err
	}

	result.ExitCode = 0
	result.Success = true
	return result, nil
}

// ExecuteADBCommandStream 执行ADB命令并实时输出流
func ExecuteADBCommandStream(args []string, timeout time.Duration, outputCallback func(line string)) (*CommandResult, error) {
	return ExecuteCommandStream("adb", args, timeout, outputCallback)
}

// ExecuteADBCommandWithDeviceStream 对特定设备执行ADB命令并实时输出流
func ExecuteADBCommandWithDeviceStream(deviceID string, args []string, timeout time.Duration, outputCallback func(line string)) (*CommandResult, error) {
	deviceArgs := []string{"-s", deviceID}
	deviceArgs = append(deviceArgs, args...)
	return ExecuteADBCommandStream(deviceArgs, timeout, outputCallback)
}

// ExecuteCommandFromConstants 执行来自constants的命令
// commandKey: constants.CommandMap中的键，如"ADBReboot"
// deviceID: 可选的设备ID，如果为空则不对特定设备执行
// params: 用于替换命令中的占位符，如map[string]string{"IP": "192.168.1.100"}
// timeout: 超时时间
// outputCallback: 可选的实时输出回调函数
func ExecuteCommandFromConstants(commandKey string, deviceID string, params map[string]string, timeout time.Duration, outputCallback func(line string)) (*CommandResult, error) {
	log.Printf("执行常量命令: key=%s, deviceID=%s, params=%v, timeout=%v, hasOutputCallback=%v",
		commandKey, deviceID, params, timeout, outputCallback != nil)
	// 获取命令字符串
	cmdStr, exists := constants.CommandMap[commandKey]
	if !exists {
		// 如果不是CommandMap中的键，尝试直接作为命令字符串使用
		cmdStr = commandKey
		log.Printf("命令键未在CommandMap中找到，使用原始字符串: %s", cmdStr)
	}

	// 替换参数占位符
	for key, value := range params {
		placeholder := "%s" // 目前constants中使用%s作为占位符
		if strings.Contains(cmdStr, placeholder) {
			cmdStr = strings.Replace(cmdStr, placeholder, value, 1)
		}
		// 也支持{KEY}格式的占位符
		bracketPlaceholder := "{" + key + "}"
		if strings.Contains(cmdStr, bracketPlaceholder) {
			cmdStr = strings.ReplaceAll(cmdStr, bracketPlaceholder, value)
		}
	}

	log.Printf("参数替换后的命令字符串: %s", cmdStr)

	// 解析命令字符串为参数列表
	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command after parsing")
	}

	// 移除命令名称"adb"如果存在
	if args[0] == "adb" {
		args = args[1:]
		log.Printf("移除adb前缀后的参数: %v", args)
	}

	log.Printf("最终执行参数: %v, 设备ID: %s, 使用流式输出: %v", args, deviceID, outputCallback != nil)

	// 根据是否指定设备ID选择执行函数
	if outputCallback != nil {
		// 使用流式执行
		if deviceID != "" {
			return ExecuteADBCommandWithDeviceStream(deviceID, args, timeout, outputCallback)
		}
		return ExecuteADBCommandStream(args, timeout, outputCallback)
	} else {
		// 使用非流式执行
		if deviceID != "" {
			return ExecuteADBCommandWithDevice(deviceID, args, timeout)
		}
		return ExecuteADBCommand(args, timeout)
	}
}

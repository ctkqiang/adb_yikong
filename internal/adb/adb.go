package adb

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
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
		return result, err
	}

	result.ExitCode = 0
	result.Success = true
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

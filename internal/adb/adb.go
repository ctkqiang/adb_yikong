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

	// 实时读取标准输出并写入缓冲区
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stdout.WriteString(line + "\n")
			if outputCallback != nil {
				outputCallback(line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("扫描标准输出时发生错误: %v", err)
		}
	}()

	// 实时读取标准错误并写入缓冲区
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			stderr.WriteString(line + "\n")
			if outputCallback != nil {
				outputCallback(line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("扫描标准错误时发生错误: %v", err)
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

// LogcatStream 表示一个流式日志会话
type LogcatStream struct {
	deviceID   string
	ctx        context.Context
	cancel     context.CancelFunc
	isRunning  bool
	outputChan chan string
	errorChan  chan error
}

// StartLogcatStream 启动实时日志流
// deviceID: 设备ID
// outputCallback: 日志输出回调函数，每行日志都会调用此函数
// filterArgs: 可选的过滤参数，如 []string{"*:V"} 表示显示所有级别的日志
// 返回一个LogcatStream实例，可以用于停止日志流
func StartLogcatStream(deviceID string, outputCallback func(line string), filterArgs ...string) (*LogcatStream, error) {
	log.Printf("启动实时日志流: deviceID=%s, filterArgs=%v", deviceID, filterArgs)

	ctx, cancel := context.WithCancel(context.Background())

	stream := &LogcatStream{
		deviceID:   deviceID,
		ctx:        ctx,
		cancel:     cancel,
		isRunning:  true,
		outputChan: make(chan string, 1000), // 缓冲通道，避免阻塞
		errorChan:  make(chan error, 1),
	}

	// 构建命令参数
	args := []string{"-s", deviceID, "logcat"}
	if len(filterArgs) > 0 {
		args = append(args, filterArgs...)
	}

	// 启动命令
	go func() {
		defer func() {
			stream.isRunning = false
			close(stream.outputChan)
			close(stream.errorChan)
			log.Printf("日志流已停止: deviceID=%s", deviceID)
		}()

		cmd := exec.CommandContext(ctx, "adb", args...)

		// 创建管道用于实时读取
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			stream.errorChan <- err
			return
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			stream.errorChan <- err
			return
		}

		// 启动命令
		if err := cmd.Start(); err != nil {
			stream.errorChan <- err
			return
		}

		log.Printf("实时日志流已启动: deviceID=%s, PID=%d", deviceID, cmd.Process.Pid)

		// 实时读取标准输出
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				select {
				case stream.outputChan <- line:
					// 成功发送到通道
				case <-ctx.Done():
					return
				default:
					// 通道已满，丢弃旧的日志行，保留新的
					// 这是为了在高频日志下避免内存泄漏
					select {
					case <-stream.outputChan: // 丢弃一条旧日志
						stream.outputChan <- line // 插入新日志
					default:
						// 如果还是满的，继续尝试
					}
				}
			}

			if err := scanner.Err(); err != nil {
				log.Printf("读取日志输出错误: %v", err)
			}
		}()

		// 读取标准错误
		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				line := scanner.Text()
				log.Printf("日志流错误输出: %s", line)
			}
		}()

		// 处理输出通道，调用回调函数
		go func() {
			for {
				select {
				case line, ok := <-stream.outputChan:
					if !ok {
						return
					}
					if outputCallback != nil {
						outputCallback(line)
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		// 等待命令完成或上下文取消
		err = cmd.Wait()
		if err != nil {
			// 如果是上下文取消导致的错误，是正常的
			if ctx.Err() == context.Canceled {
				log.Printf("日志流被正常停止: deviceID=%s", deviceID)
				return
			}
			log.Printf("日志流异常终止: deviceID=%s, 错误: %v", deviceID, err)
			stream.errorChan <- err
		}
	}()

	return stream, nil
}

// Stop 停止日志流
func (s *LogcatStream) Stop() {
	if s == nil || !s.isRunning {
		return
	}

	log.Printf("正在停止日志流: deviceID=%s", s.deviceID)
	s.cancel()
	s.isRunning = false
}

// IsRunning 检查日志流是否正在运行
func (s *LogcatStream) IsRunning() bool {
	if s == nil {
		return false
	}
	return s.isRunning
}

// StreamLogcat 启动实时日志流（简化接口）
// 这是一个方便使用的函数，内部调用StartLogcatStream
func StreamLogcat(deviceID string, outputCallback func(line string), filterArgs ...string) (*LogcatStream, error) {
	return StartLogcatStream(deviceID, outputCallback, filterArgs...)
}

package utilities

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"yikong/internal/logging"
)

func ExtractZipWindows(zipPath, extractPath string) error {
	cmd := exec.Command("tar", "-xf", zipPath, "-C", "C:\\")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Add_TO_WINDOW_PATH(path string) error {

	currentPath := os.Getenv("PATH")

	if strings.Contains(strings.ToLower(currentPath), strings.ToLower(path)) {
		logging.Warn("路径已存在于 PATH 环境变量中")
		return nil
	}

	cmd := exec.Command("setx", "PATH", currentPath+";"+path, "/M")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logging.Error("注意：这需要管理员权限。")
	logging.Error("如果失败，请以管理员身份运行程序。")

	return cmd.Run()
}

func InstallViaHomebrew() error {
	logging.Info("正在通过 Homebrew 安装 ADB...")

	if _, err := exec.LookPath("brew"); err != nil {
		logging.Error("未找到 Homebrew, 正在安装 Homebrew...")
	}

	cmd := exec.Command("brew", "install", "--cask", "android-platform-tools")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logging.Error("Homebrew 安装失败: %v\n", err)
		logging.Info("\nTMD尝试手动安装把...")

		return err
	}

	logging.Info("\n 已通过 Homebrew 安装 ADB!")

	return nil
}

func SetupMacAdbPath(installation_path string) error {
	logging.Info("\n 提示：如需将 ADB 添加到 PATH，请执行以下操作:")

	shellConfig := ""
	shellName := ""

	shell := os.Getenv("SHELL")

	if strings.Contains(shell, "zsh") {
		shellConfig = "~/.zshrc"
		shellName = "zsh"
	} else {
		shellConfig = "~/.bash_profile"
		shellName = "bash"
	}

	logging.Info("\n对于 %s (%s):\n", shellName, shellConfig)
	logging.Info(`  echo 'export PATH="$PATH:%s"' >> %s`, installation_path, shellConfig)
	logging.Info("\n  source %s\n", shellConfig)

	currentPath := os.Getenv("PATH")

	if !strings.Contains(currentPath, installation_path) {
		os.Setenv("PATH", currentPath+":"+installation_path)
	}

	createSymlink(installation_path)
	verifyInstallation()

	return nil
}

func createSymlink(installPath string) {
	adbBinary := filepath.Join(installPath, "adb")
	targetPath := "/usr/local/bin/adb"
	if os.Geteuid() == 0 {
		os.Remove(targetPath)
		os.Symlink(adbBinary, targetPath)
		logging.Info("\n🔗 已创建符号链接: %s -> %s\n", targetPath, adbBinary)
	} else {
		logging.Info("\n💡 提示：如需创建全局符号链接（需 sudo）:")
		logging.Info("  sudo ln -sf %s/adb /usr/local/bin/adb\n", installPath)
	}
}

func verifyInstallation() {
	logging.Info("\n🔍 正在验证安装...")

	cmd := exec.Command("adb", "version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		logging.Error(" 当前 PATH 中未找到 ADB。你可能需要：\n")
		logging.Error("  1. 重启终端")
		logging.Error("  2. 重新加载 shell 配置文件")
		logging.Error("  3. 或使用 adb 的完整路径")
	} else {
		logging.Info(" 安装成功！ADB 版本：\n%s\n", output)
	}
}

func InstallADBMacOS() error {
	logging.Info("正在macOS系统上安装ADB")

	if _, err := exec.LookPath("brew"); err != nil {
		logging.Info("未找到Homebrew，正在安装Homebrew...")
		installHomebrew()
	}

	logging.Info("通过Homebrew安装android-platform-tools...")
	cmd := exec.Command("brew", "install", "--cask", "android-platform-tools")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logging.Error("Homebrew安装失败: %v", err)
		return installMacOSManual()
	}

	logging.Info("已通过Homebrew安装ADB")
	setupMacOSPath()
	return nil
}

// Linux安装ADB
func InstallADBLinux() error {
	logging.Info("正在Linux系统上安装ADB")

	if isCommandAvailable("adb") {
		logging.Info("ADB已经安装")
		return nil
	}

	distro := detectLinuxDistro()
	logging.Info("检测到Linux发行版: %s", distro)

	switch distro {
	case "ubuntu", "debian", "linuxmint", "pop":
		return installLinuxAPT()
	case "fedora", "centos", "rhel":
		return installLinuxDNF()
	case "arch", "manjaro":
		return installLinuxPacman()
	default:
		return installLinuxManual()
	}
}

// Windows安装ADB
func InstallADBWindows() error {
	logging.Info("正在Windows系统上安装ADB")

	tempDir := os.Getenv("TEMP")
	zipPath := filepath.Join(tempDir, "platform-tools.zip")
	extractPath := "C:\\platform-tools"

	url := "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"

	logging.Info("下载平台工具...")
	if err := downloadFile(url, zipPath); err != nil {
		logging.Error("下载失败: %v", err)
		return err
	}

	logging.Info("解压文件...")
	if err := ExtractZipWindows(zipPath, extractPath); err != nil {
		logging.Error("解压失败: %v", err)
		return err
	}

	logging.Info("添加到系统PATH...")
	if err := Add_TO_WINDOW_PATH(extractPath); err != nil {
		logging.Warn("添加PATH失败: %v", err)
		logging.Info("请手动添加PATH: %s", extractPath)
	}

	os.Remove(zipPath)
	logging.Info("ADB安装完成")
	return nil
}

// 辅助函数
func installHomebrew() {
	cmd := exec.Command("/bin/bash", "-c",
		`curl -fsSL "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh"`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logging.Error("Homebrew安装失败: %v", err)
	}
}

func installMacOSManual() error {
	homeDir, _ := os.UserHomeDir()
	installPath := filepath.Join(homeDir, "platform-tools")
	url := "https://dl.google.com/android/repository/platform-tools-latest-darwin.zip"

	logging.Info("手动下载ADB...")
	os.RemoveAll(installPath)
	os.MkdirAll(installPath, 0755)

	if err := downloadAndExtract(url, installPath); err != nil {
		logging.Error("下载解压失败: %v", err)
		return err
	}

	setupMacOSPathWithDir(installPath)
	return nil
}

func setupMacOSPath() {
	adbPath := "/opt/homebrew/bin/adb"

	if runtime.GOARCH == "amd64" {
		adbPath = "/usr/local/bin/adb"
	}

	if !isCommandAvailable("adb") {
		logging.Info("将ADB添加到PATH:")
		logging.Info("  echo 'export PATH=\"$PATH:%s\"' >> ~/.zshrc", filepath.Dir(adbPath))
		logging.Info("  source ~/.zshrc")
	}
}

func setupMacOSPathWithDir(installPath string) {
	currentPath := os.Getenv("PATH")
	if !strings.Contains(currentPath, installPath) {
		os.Setenv("PATH", currentPath+":"+installPath)
	}
	logging.Info("已将ADB添加到当前会话PATH")
	logging.Info("永久添加请运行:")
	logging.Info("  echo 'export PATH=\"$PATH:%s\"' >> ~/.zshrc", installPath)
}

func detectLinuxDistro() string {
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}
	if _, err := os.Stat("/etc/arch-release"); err == nil {
		return "arch"
	}
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "fedora"
	}
	return "unknown"
}

func installLinuxAPT() error {
	logging.Info("使用APT安装...")

	runCommand("sudo", "apt", "update")
	if err := runCommand("sudo", "apt", "install", "-y", "android-tools-adb"); err != nil {
		logging.Error("APT安装失败: %v", err)
		return installLinuxManual()
	}

	setupLinuxUdevRules()
	return nil
}

func installLinuxDNF() error {
	logging.Info("使用DNF安装...")

	if err := runCommand("sudo", "dnf", "install", "-y", "android-tools"); err != nil {
		logging.Error("DNF安装失败: %v", err)
		return installLinuxManual()
	}

	setupLinuxUdevRules()
	return nil
}

func installLinuxPacman() error {
	logging.Info("使用Pacman安装...")

	if err := runCommand("sudo", "pacman", "-S", "--noconfirm", "android-tools"); err != nil {
		logging.Error("Pacman安装失败: %v", err)
		return installLinuxManual()
	}

	setupLinuxUdevRules()
	return nil
}

func installLinuxManual() error {
	homeDir, _ := os.UserHomeDir()
	installPath := filepath.Join(homeDir, "platform-tools")
	url := "https://dl.google.com/android/repository/platform-tools-latest-linux.zip"

	logging.Info("手动下载ADB...")
	os.RemoveAll(installPath)
	os.MkdirAll(installPath, 0755)

	if err := downloadAndExtract(url, installPath); err != nil {
		logging.Error("下载解压失败: %v", err)
		return err
	}

	setupLinuxPath(installPath)
	setupLinuxUdevRules()
	return nil
}

func setupLinuxPath(installPath string) {
	shell := os.Getenv("SHELL")
	shellConfig := "~/.bashrc"

	if strings.Contains(shell, "zsh") {
		shellConfig = "~/.zshrc"
	}

	currentPath := os.Getenv("PATH")
	if !strings.Contains(currentPath, installPath) {
		os.Setenv("PATH", currentPath+":"+installPath)
	}

	logging.Info("已将ADB添加到当前会话PATH")
	logging.Info("永久添加请运行:")
	logging.Info("  echo 'export PATH=\"$PATH:%s\"' >> %s", installPath, shellConfig)
}

func setupLinuxUdevRules() {
	if os.Geteuid() != 0 {
		logging.Info("需要root权限设置udev规则")
		logging.Info("请手动运行以下命令:")
		logging.Info("  sudo cp 51-android.rules /etc/udev/rules.d/")
		logging.Info("  sudo udevadm control --reload-rules")
		logging.Info("  sudo udevadm trigger")
		return
	}

	// 创建udev规则文件
	rulesPath := "/etc/udev/rules.d/51-android.rules"
	rulesContent := `SUBSYSTEM=="usb", ATTR{idVendor}=="0bb4", MODE="0666"
SUBSYSTEM=="usb", ATTR{idVendor}=="18d1", MODE="0666"
SUBSYSTEM=="usb", ATTR{idVendor}=="04e8", MODE="0666"`

	os.WriteFile(rulesPath, []byte(rulesContent), 0644)
	runCommand("sudo", "udevadm", "control", "--reload-rules")
	runCommand("sudo", "udevadm", "trigger")

	logging.Info("已设置udev规则")
}

// 通用函数
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func downloadAndExtract(url, dest string) error {
	tempFile, err := os.CreateTemp("", "adb-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if err := downloadFile(url, tempFile.Name()); err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return ExtractZipWindows(tempFile.Name(), dest)
	}

	cmd := exec.Command("unzip", "-o", tempFile.Name(), "-d", dest)
	return cmd.Run()
}

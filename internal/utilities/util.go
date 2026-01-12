package utilities

import (
	"os"
	"os/exec"
	"path/filepath"
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
	logging.Info("\n📝 To add ADB to your PATH:")

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

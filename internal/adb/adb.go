package adb

import (
	"bytes"
	"os/exec"
	"strings"
)

type Device struct {
	ID   string
	Name string
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

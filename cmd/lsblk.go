package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// BlockDevice 结构体表示一个块设备
type BlockDevice struct {
	Name       string
	MajorMinor string
	Size       string
	Type       string
	MountPoint string
}

// ListBlockDevices 函数用于列出所有块设备的信息
func ListBlockDevices() ([]BlockDevice, error) {
	// 执行 lsblk 命令，获取块设备信息
	cmd := exec.Command("lsblk", "-o", "NAME,MAJ:MIN,SIZE,TYPE,MOUNTPOINT")
	output := bytes.Buffer{}
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// 解析 lsblk 输出，提取块设备信息
	lines := strings.Split(output.String(), "\n")
	devices := []BlockDevice{}
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		device := BlockDevice{
			Name:       fields[0],
			MajorMinor: fields[1],
			Size:       fields[2],
			Type:       fields[3],
		}
		if len(fields) > 4 {
			device.MountPoint = fields[4]
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// Lsblk 实现lsblk
func Lsblk() {
	devices, err := ListBlockDevices()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("%-8s %6s %10s %6s %s\n", "NAME", "MAJ:MIN", "SIZE", "TYPE", "MOUNTPOINT")
	for _, device := range devices {
		fmt.Printf("%-8s %6s %10s %6s %s\n", device.Name, device.MajorMinor, device.Size, device.Type, device.MountPoint)
	}
}

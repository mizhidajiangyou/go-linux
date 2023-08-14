package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// ScsiDevice 结构体表示一个 SCSI 设备的信息
type ScsiDevice struct {
	Host    string
	Channel string
	Target  string
	Lun     string
	Vendor  string
	Model   string
	Rev     string
	Size    string
	Type    string
	Mount   string
}

// ListScsiDevices 函数用于列出所有 SCSI 设备的信息
func ListScsiDevices() ([]ScsiDevice, error) {
	// 执行 lsscsi 命令，获取 SCSI 设备信息
	cmd := exec.Command("lsscsi", "-s")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析 lsscsi 输出，提取 SCSI 设备信息
	lines := strings.Split(string(output), "\n")
	devices := []ScsiDevice{}
	for _, line := range lines[:len(lines)-1] {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		device := ScsiDevice{
			Host:    fields[0],
			Channel: fields[1],
			Target:  fields[2],
			Lun:     fields[3],
			Vendor:  fields[4],
			Model:   fields[5],
			Rev:     fields[6],
			Size:    fields[7],
			Type:    fields[8],
		}
		if len(fields) > 9 {
			device.Mount = fields[9]
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// Lsscsi 实现Lsscsi
func Lsscsi() {
	devices, err := ListScsiDevices()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%-4s %-8s %-6s %-6s %-8s %-16s %-8s %-8s %-10s %s\n", "HOST", "CHANNEL", "TARGET", "LUN", "VENDOR", "MODEL", "REV", "SIZE", "TYPE", "MOUNT")
	for _, device := range devices {
		fmt.Printf("%-4s %-8s %-6s %-6s %-8s %-16s %-8s %-8s %-10s %s\n", device.Host, device.Channel, device.Target, device.Lun, device.Vendor, device.Model, device.Rev, device.Size, device.Type, device.Mount)
	}
}

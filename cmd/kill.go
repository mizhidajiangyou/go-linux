package cmd

import (
	"fmt"
	"os"
	"syscall"
)

// KillProcess 杀死指定进程
func KillProcess(pid int) error {
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to kill process %d: %v", pid, err)
	}

	return nil
}

// ForceKillProcess 强制杀死指定进程
func ForceKillProcess(pid int) error {
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		return fmt.Errorf("failed to kill process %d: %v", pid, err)
	}

	return nil
}

func IsProcess(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}

	return false
}

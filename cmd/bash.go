package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func BashFile(scriptName string) error {
	// 将路径中的环境变量替换为其对应的值
	scriptName = os.ExpandEnv(scriptName)

	// 关键修改：只提取脚本路径（第一个空格之前的部分）
	scriptPath := strings.Split(scriptName, " ")[0]

	// 检查脚本文件是否存在（只检查脚本路径部分）
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script file '%s' does not exist", scriptPath)
	}

	// 运行脚本（使用原始完整命令行，但只检查脚本路径）
	cmd := exec.Command("/bin/bash", scriptName)
	err := cmd.Run()

	// 如果出现错误，返回错误信息
	if err != nil {
		return fmt.Errorf("failed to run script '%s': %s", scriptPath, err.Error())
	}

	// 如果执行成功，返回 nil
	if cmd.ProcessState.Success() {
		return nil
	}

	// 否则返回错误信息
	return fmt.Errorf("script '%s' not ok, return code is %d", scriptPath, cmd.ProcessState.ExitCode())
}

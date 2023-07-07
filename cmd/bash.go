package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func BashFile(scriptName string) error {
	// 将路径中的环境变量替换为其对应的值
	scriptName = os.ExpandEnv(scriptName)
	// 检查脚本文件是否存在
	if _, err := os.Stat(scriptName); os.IsNotExist(err) {
		return fmt.Errorf("script file '%s' does not exist", scriptName)
	}

	// 运行脚本文件
	cmd := exec.Command("/bin/bash", scriptName)
	err := cmd.Run()

	// 如果出现错误，返回错误信息
	if err != nil {
		return fmt.Errorf("failed to run script '%s': %s", scriptName, err.Error())
	}

	// 如果执行成功，返回 nil
	if cmd.ProcessState.Success() {
		return nil
	}

	// 否则返回错误信息
	return fmt.Errorf("script '%s' not ok , return code is %d", scriptName, cmd.ProcessState.ExitCode())
}

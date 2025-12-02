package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BashFile 执行包含脚本路径和参数的命令
func BashFile(scriptCommand string) error {
	// 1. 展开环境变量（如 $HOME）
	scriptCommand = os.ExpandEnv(scriptCommand)

	// 2. 安全分割命令参数（处理空格分隔的参数）
	parts := strings.Fields(scriptCommand)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// 3. 获取脚本路径（第一个参数）
	scriptPath := parts[0]

	// 4. 规范化路径（确保绝对路径）
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("invalid script path '%s': %w", scriptPath, err)
	}

	// 5. 检查脚本文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("script file '%s' does not exist", absPath)
	}

	// 6. 构建命令参数（/bin/bash + 脚本路径 + 其他参数）
	args := append([]string{absPath}, parts[1:]...)

	// 7. 创建并执行命令
	cmd := exec.Command("/bin/bash", args...)
	cmd.Env = os.Environ()

	// 8. 执行并返回错误（包含原始错误信息）
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run script '%s': %w", absPath, err)
	}

	return nil
}

type prefixWriter struct {
	writer io.Writer
	prefix string
}

func (w *prefixWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	output := string(p)
	prefixedOutput := w.prefix + strings.ReplaceAll(output, "\n", "\n"+w.prefix)
	if !strings.HasSuffix(prefixedOutput, "\n") {
		prefixedOutput += "\n"
	}
	return w.writer.Write([]byte(prefixedOutput))
}

func BashFileConsole(scriptCommand string) error {
	scriptCommand = os.ExpandEnv(scriptCommand)
	parts := strings.Fields(scriptCommand)
	if len(parts) == 0 {
		return errors.New("脚本命令为空")
	}

	// 路径检查与 kubeconfig 验证
	scriptPath := parts[0]
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("脚本路径规范化失败: %s (原始路径: %s)", err, scriptPath)
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("脚本文件不存在: %s", absPath)
	}

	// 构建命令
	args := append([]string{absPath}, parts[1:]...)
	cmd := exec.Command("/bin/bash", args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	stdoutPrefix := fmt.Sprintf("[脚本输出][%s] ", filepath.Base(absPath))
	stderrPrefix := fmt.Sprintf("[脚本错误][%s] ", filepath.Base(absPath))
	cmd.Stdout = &prefixWriter{writer: os.Stdout, prefix: stdoutPrefix}
	cmd.Stderr = &prefixWriter{writer: os.Stderr, prefix: stderrPrefix}

	// 执行命令
	log.Printf("[INFO] 开始执行脚本: %s", absPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("脚本执行失败: %s (脚本路径: %s, PID: %d)", err, absPath, cmd.Process.Pid)
	}
	return nil
}

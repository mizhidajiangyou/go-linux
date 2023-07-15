package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func ExecuteCommand(ip, user, password, command string) error {
	// 创建 SSH 连接
	client, err := ssh.Dial("tcp", ip+":22", &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	})
	if err != nil {
		return err
	}
	defer client.Close()

	// 创建 SSH 会话
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput(command)
	if err != nil {
		return fmt.Errorf("failed to execute command: %s", err)
	}

	fmt.Println(string(output))

	return nil
}

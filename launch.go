package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func Launch(cli string, env []string, args []string) error {
	binPath, err := exec.LookPath(cli)
	if err != nil {
		return fmt.Errorf("%s 未找到，请确认已安装并在 PATH 中", cli)
	}
	argv := append([]string{cli}, args...)
	return syscall.Exec(binPath, argv, env)
}

package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func Launch(cli string, env []string, args []string) error {
	binPath, err := exec.LookPath(cli)
	if err != nil {
		return fmt.Errorf(tr(msgCLINotFound), cli)
	}
	argv := append([]string{cli}, args...)
	return syscall.Exec(binPath, argv, env)
}

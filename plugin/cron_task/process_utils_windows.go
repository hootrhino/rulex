//go:build windows

package cron_task

import (
	"os/exec"
	"syscall"
)

func GetSysProcAttr() *syscall.SysProcAttr {
	return nil
}

func KillProcess(proc *exec.Cmd) error {
	return proc.Process.Kill()
}

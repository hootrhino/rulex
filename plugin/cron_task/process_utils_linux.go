//go:build linux

package cron_task

import (
	"os/exec"
	"syscall"
)

func GetSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func KillProcess(proc *exec.Cmd) error {
	return syscall.Kill(-proc.Process.Pid, syscall.SIGKILL)
}

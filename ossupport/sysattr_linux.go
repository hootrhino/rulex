package ossupport

import (
	"syscall"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid:     true,
		Pdeathsig:  syscall.SIGTERM,
		Cloneflags: syscall.CLONE_NEWUTS,
	}
}

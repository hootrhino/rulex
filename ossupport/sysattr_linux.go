package ossupport

import (
	// "os"
	"syscall"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		// Setsid: true,
		Setpgid: true,
		// Credential: &syscall.Credential{
		// 	Uid: 0,
		// 	Gid: 0,
		// },
	}
}

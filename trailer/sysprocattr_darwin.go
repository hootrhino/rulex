package trailer

import "syscall"

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: false,
	}
}

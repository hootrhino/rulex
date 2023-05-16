package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hootrhino/rulex/glogger"
)


/*
*
* Get Local IP
*
 */
func HostNameI() ([]string, error) {
	return []string{}, nil
}

/*
*
* 获取设备树
*
 */
type WindowsDevices struct {
	Uarts  []string `json:"uarts"`
	Videos []string `json:"videos"`
	Audios []string `json:"audios"`
}

func GetLinuxDevices() (WindowsDevices, error) {
	WindowsDevices := WindowsDevices{
		Uarts:  []string{},
		Videos: []string{},
		Audios: []string{},
	}
	return WindowsDevices, nil
}

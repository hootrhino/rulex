package utils

import (
	"net"
	"os"
)

func HostNameI() ([]string, error) {
	// ws://192.168.150.100:2580/ws
	host, _ := os.Hostname()
	addrs, _ := net.LookupHost(host)
	addrsL := []string{}
	for _, addr := range addrs {
		if len(addr) <= 28 {
			addrsL = append(addrsL, addr)
		}
	}
	return addrsL, nil
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

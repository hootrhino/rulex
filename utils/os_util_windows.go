package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

/*
*
* 获取IP地址
*
 */
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

func GetSystemDevices() (SystemDevices, error) {
	SystemDevices := SystemDevices{
		Uarts:  []string{},
		Videos: []string{},
		Audios: []string{},
	}
	return SystemDevices, nil
}
func CatOsRelease() (map[string]string, error) {
	return map[string]string{"os": "windows"}, nil
}

/*
*
* 执行系统命令
*
 */
func OsExecute(name string, arg ...string) error {
	nmcliCmd := exec.Command(name, arg...)
	output, err := nmcliCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error when executing[%s %v]:%s", name, arg, err.Error()+", output:"+string(output))
	}
	return nil
}

package utils

import (
	"os"
	"os/exec"

	"strings"
)

/*
*
* Get Local IP
*
 */
func HostNameI() ([]string, error) {
	cmd := exec.Command("hostname", "-I")
	data, err1 := cmd.Output()
	if err1 != nil {
		return []string{}, err1
	}
	ss := []string{}
	for _, s := range strings.Split(string(data), " ") {
		if s != "\n" {
			ss = append(ss, s)
		}
	}
	return ss, nil
}

/*
*
* 获取设备树
*
 */
type LinuxDevices struct {
	Uarts  []string `json:"uarts"`
	Videos []string `json:"videos"`
	Audios []string `json:"audios"`
}

func GetLinuxDevices() (LinuxDevices, error) {
	LinuxDevices := LinuxDevices{
		Uarts:  []string{},
		Videos: []string{},
		Audios: []string{},
	}
	f, err := os.Open("/dev/")
	if err != nil {
		return LinuxDevices, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return LinuxDevices, err
	}

	for _, d := range list {
		if !d.IsDir() {
			if strings.Contains(d.Name(), "ttyS") {
				LinuxDevices.Uarts = append(LinuxDevices.Uarts, d.Name())
			}
			if strings.Contains(d.Name(), "ttyACM") {
				LinuxDevices.Uarts = append(LinuxDevices.Uarts, d.Name())
			}
			if strings.Contains(d.Name(), "ttyUSB") {
				LinuxDevices.Uarts = append(LinuxDevices.Uarts, d.Name())
			}
			if strings.Contains(d.Name(), "video") {
				LinuxDevices.Videos = append(LinuxDevices.Videos, d.Name())
			}
			if strings.Contains(d.Name(), "audio") {
				LinuxDevices.Audios = append(LinuxDevices.Audios, d.Name())
			}
		}
	}
	return LinuxDevices, nil
}

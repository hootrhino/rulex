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
	dist, _ := GetOSDistribution()
	if dist == "openwrt" {
		line := `ip addr show | awk '/inet / {print $2}' | awk 'BEGIN{FS="/"} {split($0, arr, "/"); print arr[1]}'`
		cmd := exec.Command("sh", "-c", line)
		output, err := cmd.Output()
		if err != nil {
			return []string{}, err
		}
		result := strings.TrimSpace(string(output))
		ips := []string{}
		for _, v := range strings.Split(result, "\n") {
			if v != "127.0.0.1" {
				ips = append(ips, v)
			}
		}
		return ips, nil
	}
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
func GetSystemDevices() (SystemDevices, error) {
	SystemDevices := SystemDevices{
		Uarts:  []string{},
		Videos: []string{},
		Audios: []string{},
	}
	f, err := os.Open("/dev/")
	if err != nil {
		return SystemDevices, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return SystemDevices, err
	}

	for _, d := range list {
		if !d.IsDir() {
			if strings.Contains(d.Name(), "ttyS") {
				SystemDevices.Uarts = append(SystemDevices.Uarts, d.Name())
			}
			if strings.Contains(d.Name(), "ttyACM") {
				SystemDevices.Uarts = append(SystemDevices.Uarts, d.Name())
			}
			if strings.Contains(d.Name(), "ttyUSB") {
				SystemDevices.Uarts = append(SystemDevices.Uarts, d.Name())
			}
			if strings.Contains(d.Name(), "video") {
				SystemDevices.Videos = append(SystemDevices.Videos, d.Name())
			}
			if strings.Contains(d.Name(), "audio") {
				SystemDevices.Audios = append(SystemDevices.Audios, d.Name())
			}
		}
	}
	return SystemDevices, nil
}

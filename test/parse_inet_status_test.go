// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// DEVICE           TYPE      STATE      CONNECTION
// eth0             ethernet  connected  eth0
// usb0             ethernet  connected  usb0
// wlx0cc6551c5026  wifi      connected  iotlab4072
// eth1             ethernet  connected  eth1
// lo               loopback  unmanaged  --

// go test -timeout 30s -run ^Test_parse_network_status github.com/hootrhino/rulex/test -v -count=1

// NetworkDevice 表示网络设备
type NetworkDevice struct {
	Device     string
	Type       string
	State      string
	Connection string
}

func Test_parse_network_status(t *testing.T) {
	input := `DEVICE           TYPE      STATE      CONNECTION
eth0             ethernet  connected  eth0
usb0             ethernet  connected  usb0
wlx0cc6551c5026  wifi      connected  iotlab4072
eth1             ethernet  connected  eth1
lo               loopback  unmanaged  --`

	devices, err := parseNetworkDevices(input)
	if err != nil {
		fmt.Println("Error parsing network devices:", err)
		return
	}

	// 打印解析结果
	for _, device := range devices {
		fmt.Printf("Device: %s, Type: %s, State: %s, Connection: %s\n", device.Device, device.Type, device.State, device.Connection)
	}
}

// parseNetworkDevices 解析网络设备信息
func parseNetworkDevices(input string) ([]NetworkDevice, error) {
	var devices []NetworkDevice

	// 将输入按行分割
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "DEVICE") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		device := NetworkDevice{
			Device:     fields[0],
			Type:       fields[1],
			State:      fields[2],
			Connection: fields[3],
		}
		devices = append(devices, device)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return devices, nil
}

// NetworkDevice 表示网络设备
type NetworkDeviceDetail struct {
	Device      string
	Type        string
	HWAddr      string
	MTU         int
	State       string
	Connection  string
	Carrier     string
	IPv4Addr    string
	IPv4Gateway string
	IPv4DNS     string
	IPv6Addr    string
	IPv6Gateway string
}

// go test -timeout 30s -run ^Test_parse_nmcliDeviceShowOutput github.com/hootrhino/rulex/test -v -count=1
func Test_parse_nmcliDeviceShowOutput(t *testing.T) {

	// 替换为你的 nmcli 输出
	nmcliOutput := `
GENERAL.DEVICE:                         eth0
GENERAL.TYPE:                           ethernet
GENERAL.HWADDR:                         02:81:5E:DF:D4:81
GENERAL.MTU:                            1500
GENERAL.STATE:                          100 (connected)
GENERAL.CONNECTION:                     eth0
WIRED-PROPERTIES.CARRIER:               on
IP4.ADDRESS[1]:                         192.168.1.185/24
IP4.GATEWAY:                            192.168.1.1
IP4.DNS[1]:                             192.168.1.1
IP6.ADDRESS[1]:                         fe80::9460:7480:61a9:cbd2/64
IP6.GATEWAY:                            --
`

	// 解析 nmcli 输出
	device, err := parseNmcliDeviceShow(nmcliOutput)
	if err != nil {
		fmt.Println("Error parsing nmcli output:", err)
		return
	}

	// 打印解析结果
	fmt.Printf("Device: %s\n", device.Device)
	fmt.Printf("Type: %s\n", device.Type)
	fmt.Printf("HWAddr: %s\n", device.HWAddr)
	fmt.Printf("MTU: %d\n", device.MTU)
	fmt.Printf("State: %s\n", device.State)
	fmt.Printf("Connection: %s\n", device.Connection)
	fmt.Printf("Carrier: %s\n", device.Carrier)
	fmt.Printf("IPv4Addr: %s\n", device.IPv4Addr)
	fmt.Printf("IPv4Gateway: %s\n", device.IPv4Gateway)
	fmt.Printf("IPv4DNS: %s\n", device.IPv4DNS)
	fmt.Printf("IPv6Addr: %s\n", device.IPv6Addr)
	fmt.Printf("IPv6Gateway: %s\n", device.IPv6Gateway)
}

// parseNMCLIOutput 解析 nmcli 输出
// nmcli device show
func parseNmcliDeviceShow(output string) (*NetworkDeviceDetail, error) {
	lines := strings.Split(output, "\n")

	device := &NetworkDeviceDetail{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "GENERAL.DEVICE:":
			device.Device = fields[1]
		case "GENERAL.TYPE:":
			device.Type = fields[1]
		case "GENERAL.HWADDR:":
			device.HWAddr = fields[1]
		case "GENERAL.MTU:":
			device.MTU = parseInt(fields[1])
		case "GENERAL.STATE:":
			device.State = fields[1]
		case "GENERAL.CONNECTION:":
			device.Connection = fields[1]
		case "WIRED-PROPERTIES.CARRIER:":
			device.Carrier = fields[1]
		case "IP4.ADDRESS[1]:":
			device.IPv4Addr = fields[1]
		case "IP4.GATEWAY:":
			device.IPv4Gateway = fields[1]
		case "IP4.DNS[1]:":
			device.IPv4DNS = fields[1]
		case "IP6.ADDRESS[1]:":
			device.IPv6Addr = fields[1]
		case "IP6.GATEWAY:":
			device.IPv6Gateway = fields[1]
		}
	}

	return device, nil
}

// parseInt 将字符串转换为整数，如果失败返回 0
func parseInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}

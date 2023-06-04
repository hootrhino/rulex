package utils

import "encoding/json"

/*
*
* 磁盘
*
 */
type DiskUsage struct {
	DeviceID  string  `json:"deviceID"`
	FreeSpace float64 `json:"freeSpace"`
	Size      float64 `json:"size"`
}

func (m DiskUsage) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

/*
*
* CPU
*
 */
type CpuUsage struct {
	Name  string `json:"name"`
	Usage uint64 `json:"usage"`
}

func (m CpuUsage) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

/*
*
* 网卡
*
 */
type NetworkInterfaceUsage struct {
	Name                string
	CurrentBandwidth    uint64
	BytesTotalPerSec    uint64
	BytesReceivedPerSec uint64
	BytesSentPerSec     uint64
	PacketsPerSec       uint64
}

func (m NetworkInterfaceUsage) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

/*
*
* 系统设备
*
 */
type SystemDevices struct {
	Uarts  []string `json:"uarts"`
	Videos []string `json:"videos"`
	Audios []string `json:"audios"`
}

func (m SystemDevices) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

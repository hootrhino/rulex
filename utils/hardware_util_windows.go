package utils

import (
	"fmt"
	"strconv"

	"github.com/hootrhino/wmi"
)

/*
*
* Get-WmiObject -query "SELECT Name, PercentIdleTime FROM Win32_PerfFormattedData_PerfOS_Processor"
*
 */
func GetCpuUsage() ([]CpuUsage, error) {
	type Model struct {
		Name            string
		PercentIdleTime uint64
	}
	var PerfOS_Processor []Model
	err := wmi.Query(`SELECT Name, PercentIdleTime FROM Win32_PerfFormattedData_PerfOS_Processor`,
		&PerfOS_Processor)
	if err != nil {
		return nil, err
	}
	Usages := []CpuUsage{}
	for _, v := range PerfOS_Processor {
		Usages = append(Usages, CpuUsage{
			Name:  "CPU:" + v.Name,
			Usage: 100 - v.PercentIdleTime,
		})
	}
	return Usages, nil
}

/*
*
Get-WmiObject -query "SELECT * FROM Win32_logicalDisk"

DeviceID     : C:
DriveType    : 3
ProviderName :
FreeSpace    : 219526295552
Size         : 999360032768
VolumeName   :
*/
func GetDiskUsage() ([]DiskUsage, error) {
	type Model struct {
		DeviceID  string
		FreeSpace uint64
		Size      uint64
	}
	var models []Model
	err := wmi.Query(`SELECT DeviceID, FreeSpace, Size FROM Win32_logicalDisk`,
		&models)
	if err != nil {
		return nil, err
	}
	data := []DiskUsage{}
	for _, v := range models {
		data = append(data, DiskUsage{
			DeviceID:  v.DeviceID,
			FreeSpace: Decimal(float64(v.FreeSpace) / 1024 / 1024 / 1024), // Mb
			Size:      Decimal(float64(v.Size) / 1024 / 1024 / 1024),      // Mb
		})
	}
	return data, nil
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

/*
*
* 获取网卡实时速率

  - SELECT
    Name,
    CurrentBandwidth,
    BytesTotalPerSec,
    BytesReceivedPerSec,
    BytesSentPerSec,
    PacketsPerSec
    FROM Win32_PerfFormattedData_Tcpip_NetworkInterface

    -----------------

    BytesReceivedPersec : 1114
    BytesSentPersec     : 608
    BytesTotalPersec    : 1722
    CurrentBandwidth    : 1000000000
    Name                : Intel[R] Ethernet Connection [12] I219-V
    PacketsPersec       : 23
    PSComputerName      :

*
*/
func NetInterfaceUsage() ([]NetworkInterfaceUsage, error) {
	sql := `
SELECT
Name,
CurrentBandwidth,
BytesTotalPerSec,
BytesReceivedPerSec,
BytesSentPerSec,
PacketsPerSec
FROM Win32_PerfFormattedData_Tcpip_NetworkInterface
`
	var model []NetworkInterfaceUsage
	err := wmi.Query(sql, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

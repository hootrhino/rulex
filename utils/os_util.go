package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/i4de/rulex/glogger"
)

/*
*
* GetPwd
*
 */
func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	return dir
}

/*
*
* Byte to Mbyte
*
 */
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

/*
*
* Get Local IP
*
 */
func HostNameI() (string, error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("hostname", "-I")
		data, err1 := cmd.Output()
		if err1 != nil {
			return "", err1
		}
		return string(data), nil
	}
	return "[0.0.0.0]only support unix-like OS", nil
}

func toMegaBytes(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

/*
*
* DEBUG使用
*
 */
func TraceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	var info [7]float64
	info[0] = float64(ms.HeapObjects)
	info[1] = toMegaBytes(ms.HeapAlloc)
	info[2] = toMegaBytes(ms.TotalAlloc)
	info[3] = toMegaBytes(ms.HeapSys)
	info[4] = toMegaBytes(ms.HeapIdle)
	info[5] = toMegaBytes(ms.HeapReleased)
	info[6] = toMegaBytes(ms.HeapIdle - ms.HeapReleased)

	for _, v := range info {
		fmt.Printf("%v,\t", v)
	}
	fmt.Println()
}

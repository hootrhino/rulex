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
func HostNameI() ([]string, error) {
	if runtime.GOOS == "linux" {
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
	return []string{}, nil
}

func BtoMB(bytes uint64) float64 {
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
	info[1] = BtoMB(ms.HeapAlloc)
	info[2] = BtoMB(ms.TotalAlloc)
	info[3] = BtoMB(ms.HeapSys)
	info[4] = BtoMB(ms.HeapIdle)
	info[5] = BtoMB(ms.HeapReleased)
	info[6] = BtoMB(ms.HeapIdle - ms.HeapReleased)

	for _, v := range info {
		fmt.Printf("%v,\t", v)
	}
	fmt.Println()
}

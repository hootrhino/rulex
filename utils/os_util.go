package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

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
	var result = make([]float64, 7)
	result[0] = float64(ms.HeapObjects)
	result[1] = toMegaBytes(ms.HeapAlloc)
	result[2] = toMegaBytes(ms.TotalAlloc)
	result[3] = toMegaBytes(ms.HeapSys)
	result[4] = toMegaBytes(ms.HeapIdle)
	result[5] = toMegaBytes(ms.HeapReleased)
	result[6] = toMegaBytes(ms.HeapIdle - ms.HeapReleased)

	fmt.Printf("%d\t", time.Now().Unix())
	for _, v := range result {
		fmt.Printf("%.2f\t", v)
	}
	fmt.Printf("\n")
	time.Sleep(1 * time.Second)
}

func TraceMem() {
	TraceMemStats()
	var container [200 * 1024 * 1024]byte
	for i := 0; i < 200*1024*1024; i++ {
		container[i] = 0
	}
	TraceMemStats()
	container[0] = 0
	log.Printf("%d", len(container))
}

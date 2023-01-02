package utils

import (
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

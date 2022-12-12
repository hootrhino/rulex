package utils

import (
	"os"

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
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

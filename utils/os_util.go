package utils

import (
	"log"
	"os"
)

/*
*
* GetPwd
*
 */
func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

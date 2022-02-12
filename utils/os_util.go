package utils

import (
	"log"
	"os"
)

/*
*
* pwd
*
 */
func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

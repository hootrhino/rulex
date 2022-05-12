package utils

import (
	"log"
	"os"
)

/*
#include "utils.h"
*/
import "C"

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

func GetPid() int {
	return int(C.GetPid())
}

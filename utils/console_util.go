package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ngaut/log"
)

//
// show banner
//
func ShowBanner() {
	//
	defaultBanner :=
		`
---------------------------------------------
,------. ,--. ,--.,--.   ,------.,--.   ,--.
|  .--. '|  | |  ||  |   |  .---' \  '.'  /
|  '--'.'|  | |  ||  |   |  '--,   .'    \
|  |\  \ '  '-'  '|  '--.|  '---. /  .'.  \
'--' '--' '-----' '-----''------''--'   '--'
---------------------------------------------
`
	file, err := os.Open("conf/banner.txt")
	if err != nil {
		log.Warn("No banner found, print default banner")
		log.Info(defaultBanner)
	} else {
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Warn("No banner found, print default banner")
			log.Info(defaultBanner)
			fmt.Println("\n", defaultBanner)
		} else {
			fmt.Println("\n", string(data))
		}
	}
}

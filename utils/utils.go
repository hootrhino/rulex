package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/ngaut/log"
)

//
// MakeUUID
//
func MakeUUID(prefix string) string {
	return prefix + "_" + uuid.NewString()
}

//
//
//

func Post(data interface{}, api string) (string, error) {
	bites, errs1 := json.Marshal(data)
	if errs1 != nil {
		log.Error(errs1)
		return "", errs1
	}
	r, errs2 := http.Post(api, "application/json",
		bytes.NewBuffer(bites))
	if errs2 != nil {
		log.Error(errs2)
		return "", errs2
	}
	defer r.Body.Close()
	body, errs3 := ioutil.ReadAll(r.Body)
	if errs3 != nil {
		log.Error(errs3)
		return "", errs3
	}
	return string(body), nil
}

//
// show banner
//
func ShowBanner() {
	//
	defaultBanner :=
		`
-----------------------------------------------------------
~~~/=====\       ██████╗ ██╗   ██╗██╗     ███████╗██╗  ██╗
~~~||\\\||--->o  ██╔══██╗██║   ██║██║     ██╔════╝╚██╗██╔╝
~~~||///||--->o  ██████╔╝██║   ██║██║     █████╗   ╚███╔╝ 
~~~||///||--->o  ██╔══██╗██║   ██║██║     ██╔══╝   ██╔██╗ 
~~~||\\\||--->o  ██║  ██║╚██████╔╝███████╗███████╗██╔╝ ██╗
~~~\=====/       ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝
-----------------------------------------------------------
`
	file, err := os.Open("conf/banner.txt")
	defer file.Close()
	if err != nil {
		log.Warn("No banner found, print default banner")
		log.Info(defaultBanner)
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Warn("No banner found, print default banner")
			log.Info(defaultBanner)
		} else {
			log.Info("\n", string(data))
		}
	}
	log.Info("Rulex start successfully")
}

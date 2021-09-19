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

func post(data map[string]interface{}, api string) string {
	p, errs1 := json.Marshal(data)
	if errs1 != nil {
		log.Error(errs1)
	}
	r, errs2 := http.Post(api, "application/json",
		bytes.NewBuffer(p))
	if errs2 != nil {
		log.Error(errs2)
	}
	defer r.Body.Close()

	body, errs5 := ioutil.ReadAll(r.Body)
	if errs5 != nil {
		log.Error(errs5)
	}
	return string(body)
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
	log.Info("rulex start successfully")
}

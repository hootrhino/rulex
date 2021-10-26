package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
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

//
// JSON String to a struct
//
func TransformConfig(s1 []byte, s2 interface{}) error {
	if err := json.Unmarshal(s1, &s2); err != nil {
		return err
	}
	if err := validator.New().Struct(s2); err != nil {
		return err
	}
	return nil
}

//
// Bind config to struct
// config: a Map, s: a struct variable
//
func BindResourceConfig(config *map[string]interface{}, s interface{}) error {
	configBytes, err0 := json.Marshal(config)
	if err0 != nil {
		return err0
	}
	if err := json.Unmarshal(configBytes, &s); err != nil {
		return err
	}
	if err := validator.New().Struct(s); err != nil {
		return err
	}
	return nil
}

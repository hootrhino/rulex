package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"unicode"

	"github.com/go-playground/validator/v10"

	"github.com/google/uuid"
	"github.com/ngaut/log"
)

//
// MakeUUID
//
func InUuid() string {
	return MakeUUID("INEND")
}

//
// MakeUUID
//
func OutUuid() string {
	return MakeUUID("OUTEND")
}

//
// MakeUUID
//
func RuleUuid() string {
	return MakeUUID("RULE")
}

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
func BindResourceConfig(config map[string]interface{}, s interface{}) error {
	configBytes, err0 := json.Marshal(&config)
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

//
// 去掉\u0000字符
//
func TrimZero(s string) string {
	str := make([]rune, 0, len(s))
	for _, v := range s {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) {
			continue
		}
		str = append(str, v)
	}
	return string(str)
}

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

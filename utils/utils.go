package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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

func Post(client http.Client, data interface{}, api string, headers map[string]string) (string, error) {

	bites, errs1 := json.Marshal(data)
	if errs1 != nil {
		log.Error(errs1)
		return "", errs1
	}
	body := strings.NewReader(string(bites))
	request, _ := http.NewRequest("POST", api, body)
	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	response, err2 := client.Do(request)
	if err2 != nil {
		return "", err2
	}
	if response.StatusCode != 200 {
		bytes0, err3 := ioutil.ReadAll(response.Body)
		if err3 != nil {
			return "", err3
		}
		return "", fmt.Errorf("Error:%v", string(bytes0))
	}
	var r []byte
	response.Body.Read(r)
	bytes1, err3 := ioutil.ReadAll(response.Body)
	if err3 != nil {
		return "", err3
	}
	return string(bytes1), nil
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
}

//
// JSON String to a struct, (can't validate map!!!)
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
func BindConfig(config map[string]interface{}, s interface{}) error {
	return BindResourceConfig(config, s)
}
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

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ngaut/log"
)

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

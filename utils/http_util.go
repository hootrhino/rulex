package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hootrhino/rulex/glogger"
)

/*
*
* HTTP POST
*
 */
func Post(client http.Client, data interface{},
	url string, headers map[string]string) (string, error) {
	bites, errs1 := json.Marshal(data)
	if errs1 != nil {
		glogger.GLogger.Error(errs1)
		return "", errs1
	}
	body := strings.NewReader(string(bites))
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	response, err2 := client.Do(request)
	if err2 != nil {
		return "", err2
	}
	if response.StatusCode != 200 {
		bytes0, err3 := io.ReadAll(response.Body)
		if err3 != nil {
			return "", err3
		}
		return "", fmt.Errorf("Error:%v", string(bytes0))
	}
	var r []byte
	response.Body.Read(r)
	bytes1, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		return "", err3
	}
	return string(bytes1), nil
}

/*
*
* HTTP GET
*
 */
func Get(client http.Client, url string) string {
	var err error
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}

	response, err := client.Do(request)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	return string(body)
}

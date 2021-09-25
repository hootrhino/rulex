package cloud

import (
	"io/ioutil"
	"net/http"

	"github.com/ngaut/log"
)

//
// 请求云端
//

func GetCloud(token string, secret string, api string) string {
	var err error
	request, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Error(err)
		return ""
	}
	request.Header.Set("token", token)
	request.Header.Set("secret", secret)
	response, err := (&http.Client{}).Do(request)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(body)
}

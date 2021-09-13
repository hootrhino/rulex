package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
)

var router *gin.Engine

type Service struct {
	Name        string
	Doc         string
	Description string
}

func T1() {
	router = gin.Default()
	router.GET("/services", func(c *gin.Context) {
		datas := []Service{}
		datas = append(datas, Service{
			Name:        "语音识别",
			Doc:         "http://www.google.com",
			Description: "语音识别接口",
		})
		datas = append(datas, Service{
			Name:        "人脸识别",
			Doc:         "http://www.google.com",
			Description: "人脸识别接口",
		})
		datas = append(datas, Service{
			Name:        "活体识别",
			Doc:         "http://www.google.com",
			Description: "活体识别接口",
		})
		c.JSON(http.StatusOK, datas)

	})
	router.Run(":9990")
}

func GetCloudApi(token string, secret string, api string) string {
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

func TestCloudApi(t *testing.T) {
	r := GetCloudApi("49ba59abbe56e057", "49ba59abbe56e057", "https://atomicservices.000webhostapp.com/data.json")
	t.Log("-------------------------->\n", r)

}

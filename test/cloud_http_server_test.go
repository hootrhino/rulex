package test

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

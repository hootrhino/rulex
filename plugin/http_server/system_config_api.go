package httpserver

import "github.com/gin-gonic/gin"

// 主要是针对WIFI、时区、IP地址设置

/*
*
* WIFI
*
 */
func SetWifi(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置时间、时区
*
 */
func SetTime(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

/*
*
* 设置静态网络IP等
*
 */
func SetStaticNetwork(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
	}

}

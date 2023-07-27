package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/typex"
	"gorm.io/gorm"
)

/*
*
* API Server
*
 */
type RulexApiServer struct {
	ginEngine  *gin.Engine
	ruleEngine typex.RuleX
}

/*
*
* 新建路由
*
 */
func (ras *RulexApiServer) AddRoute(method, path string, handler func(ctx *gin.Context)) {
	if method == "GET" {
		ras.ginEngine.GET(path, handler)
	}
	if method == "POST" {
		ras.ginEngine.POST(path, handler)
	}
	if method == "PUT" {
		ras.ginEngine.PUT(path, handler)
	}
	if method == "DELETE" {
		ras.ginEngine.DELETE(path, handler)
	}
}

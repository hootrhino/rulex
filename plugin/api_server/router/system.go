package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/api_server/service"
	"github.com/hootrhino/rulex/typex"
)

func InitSystemRouter(ruleEngine typex.RuleX, Router *gin.RouterGroup) {
	registerRouter := Router.Group("system")
	system := service.NewSystemService(ruleEngine)
	registerRouter.GET("/plugins", system.GetPlugins) // 获取数据

}

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/api_server/service"
)

func InitRuleRouter(Router *gin.RouterGroup) {
	registerRouter := Router.Group("rule")
	newAccount := service.NewDemoService()
	registerRouter.GET("/list", newAccount.GetRuleList) // 获取数据

}

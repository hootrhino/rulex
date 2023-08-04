package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/api_server/service"
	"github.com/hootrhino/rulex/typex"
)

func InitRuleRouter(ruleEngine typex.RuleX, Router *gin.RouterGroup) {
	registerRouter := Router.Group("rule")
	newAccount := service.NewRuleService(ruleEngine)
	registerRouter.GET("/list", newAccount.GetRuleList) // 获取数据

}

package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/api_server/middleware"
	"github.com/hootrhino/rulex/plugin/api_server/router"
	"github.com/hootrhino/rulex/typex"
)

func Routers(ruleEngine typex.RuleX) *gin.Engine {
	Router := gin.Default()
	// 注册全局中间件,根据实际业务需求注册
	Router.Use(
		middleware.CorsMiddleWare(), // 跨域中间件
	)
	// 配置全局路径
	ApiGroup := Router.Group("/api/v2/")
	// 注册路由
	router.InitRuleRouter(ruleEngine, ApiGroup)
	router.InitSystemRouter(ruleEngine, ApiGroup)
	return Router
}

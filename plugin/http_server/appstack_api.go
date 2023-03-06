package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/i4de/rulex/typex"
)

// 列表
func Apps(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

// 新建
func CreateApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

// 更新
func UpdateApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

// 停止
func StopApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

// 删除
func RemoveApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData("ok"))
}

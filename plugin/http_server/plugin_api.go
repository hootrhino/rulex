package httpserver

import (
	"github.com/hootrhino/rulex/typex"

	"github.com/gin-gonic/gin"
)

/*
*
* 插件的服务接口
*
 */

func PluginService(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID string `json:"uuid" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	plugin, ok := e.AllPlugins().Load(form.UUID)
	if ok {
		result := plugin.(typex.XPlugin).Service(typex.ServiceArg{
			Name: form.Name,
		})
		c.JSON(200, OkWithData(result.Out))
		return
	}
	c.JSON(200, Error("plugin not exists"))
}

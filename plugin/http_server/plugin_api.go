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

func PluginService(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		UUID string      `json:"uuid" binding:"required"`
		Name string      `json:"name" binding:"required"`
		Args interface{} `json:"args"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	plugin, ok := hh.ruleEngine.AllPlugins().Load(form.UUID)
	if ok {
		result := plugin.(typex.XPlugin).Service(typex.ServiceArg{
			Name: form.Name,
			UUID: form.UUID,
			Args: form.Args,
		})
		c.JSON(HTTP_OK, OkWithData(result.Out))
		return
	}
	c.JSON(HTTP_OK, Error("plugin not exists"))
}

package apis

import (
	"fmt"

	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"

	"github.com/gin-gonic/gin"
)

/*
*
* 插件的服务接口
*
 */

func PluginService(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUID string      `json:"uuid" binding:"required"`
		Name string      `json:"name" binding:"required"`
		Args interface{} `json:"args"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	plugin, ok := ruleEngine.AllPlugins().Load(form.UUID)
	if ok {
		result := plugin.(typex.XPlugin).Service(typex.ServiceArg{
			Name: form.Name,
			UUID: form.UUID,
			Args: form.Args,
		})
		c.JSON(common.HTTP_OK, common.OkWithData(result.Out))
		return
	}
	c.JSON(common.HTTP_OK, common.Error("plugin not exists:"+form.UUID))
}

/*
*
* 插件详情
*
 */
func PluginDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	plugin, ok := ruleEngine.AllPlugins().Load(uuid)
	if ok {
		result := plugin.(typex.XPlugin)
		c.JSON(common.HTTP_OK, common.OkWithData(result.PluginMetaInfo()))
		return
	}
	c.JSON(common.HTTP_OK, common.Error400EmptyObj(fmt.Errorf("no such plugin:%s", uuid)))
}

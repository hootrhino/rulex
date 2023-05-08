package httpserver

import (
	"github.com/hootrhino/rulex/typex"

	"github.com/gin-gonic/gin"
)

/*
*
* AiBase
*
 */
func AiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		c.JSON(200, OkWithData(e.GetAiBase().ListAi()))
		return
	}
	if ai := e.GetAiBase().GetAi(uuid); ai != nil {
		c.JSON(200, OkWithData(ai))
		return
	}
	c.JSON(200, Error("ai base not found:"+uuid))
}

/*
*
* 删除
*
 */
func DeleteAiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if ai := e.GetAiBase().GetAi(uuid); ai != nil {
		err := e.GetAiBase().RemoveAi(uuid)
		if err != nil {
			c.JSON(200, Error400(err))
			return
		}
		c.JSON(200, Ok())
		return
	}
	c.JSON(200, Error("ai base not found:"+uuid))
}

/*
*
* Create
*
 */

func CreateAiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(200, Ok())
}

/*
*
* 更新
*
 */
func UpdateAiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(200, Ok())
}

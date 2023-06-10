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
		c.JSON(HTTP_OK, OkWithData(e.GetAiBase().ListAi()))
		return
	}
	if ai := e.GetAiBase().GetAi(uuid); ai != nil {
		c.JSON(HTTP_OK, OkWithData(ai))
		return
	}
	c.JSON(HTTP_OK, Error("ai base not found:"+uuid))
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
			c.JSON(HTTP_OK, Error400(err))
			return
		}
		c.JSON(HTTP_OK, Ok())
		return
	}
	c.JSON(HTTP_OK, Error("ai base not found:"+uuid))
}

/*
*
* Create
*
 */

func CreateAiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(HTTP_OK, Ok())
}

/*
*
* 更新
*
 */
func UpdateAiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(HTTP_OK, Ok())
}

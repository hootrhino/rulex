package httpserver

import (
	"github.com/gin-gonic/gin"
)

/*
*
* AiBase
*
 */
func AiBase(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		c.JSON(HTTP_OK, OkWithData(hh.ruleEngine.GetAiBase().ListAi()))
		return
	}
	if ai := hh.ruleEngine.GetAiBase().GetAi(uuid); ai != nil {
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
func DeleteAiBase(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	if ai := hh.ruleEngine.GetAiBase().GetAi(uuid); ai != nil {
		err := hh.ruleEngine.GetAiBase().RemoveAi(uuid)
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

func CreateAiBase(c *gin.Context, hh *HttpApiServer) {
	c.JSON(HTTP_OK, Ok())
}

/*
*
* 更新
*
 */
func UpdateAiBase(c *gin.Context, hh *HttpApiServer) {
	c.JSON(HTTP_OK, Ok())
}

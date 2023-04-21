package httpserver

import (
	"github.com/i4de/rulex/typex"

	"github.com/gin-gonic/gin"
)

/*
*
* AiBase
*
 */
func AiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	c.JSON(200, OkWithData(uuid))
}

/*
*
* 删除
*
 */
func DeleteAiBase(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	c.JSON(200, OkWithData(uuid))
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

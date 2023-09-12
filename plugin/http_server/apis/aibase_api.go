package apis

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* AiBase
*
 */
func AiBase(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		c.JSON(common.HTTP_OK, common.OkWithData(ruleEngine.GetAiBase().ListAi()))
		return
	}
	if ai := ruleEngine.GetAiBase().GetAi(uuid); ai != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(ai))
		return
	}
	c.JSON(common.HTTP_OK, common.Error("ai base not found:"+uuid))
}

/*
*
* 删除
*
 */
func DeleteAiBase(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if ai := ruleEngine.GetAiBase().GetAi(uuid); ai != nil {
		err := ruleEngine.GetAiBase().RemoveAi(uuid)
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("ai base not found:"+uuid))
}

/*
*
* Create
*
 */

func CreateAiBase(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 更新
*
 */
func UpdateAiBase(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
}

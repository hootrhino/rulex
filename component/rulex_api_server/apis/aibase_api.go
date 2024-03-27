package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/aibase"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* AiBase
*
 */
func AiBaseList(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData(aibase.ListAlgorithm()))
}
func AiBaseDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if ai := aibase.GetAlgorithm(uuid); ai != nil {
		c.JSON(common.HTTP_OK, common.OkWithData(ai))
		return
	}
	c.JSON(common.HTTP_OK, common.Error("not found:"+uuid))
}

/*
*
* 删除
*
 */
func DeleteAiBase(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.Ok())
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

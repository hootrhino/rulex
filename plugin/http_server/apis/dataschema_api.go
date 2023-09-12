package apis

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/core"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 新建模型
*
 */

func CreateDataSchema(c *gin.Context, ruleEngine typex.RuleX) {
	vo := typex.DataSchema{}
	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	bytes, err := json.Marshal(vo.Schema)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if len(vo.Schema) < 1 {
		c.JSON(common.HTTP_OK, common.Error("Must contain less 1 Filed"))
		return
	}
	/*
	*
	* Type: SIMPLE_LINE(简单一线),COMPLEX_LINE(复杂多线)
	*
	 */
	if !utils.SContains([]string{
		"SIMPLE_DATE_LINE",
		"COMPLEX_DATE_LINE",
		"COMPLEX_CROSS_LINE",
	}, vo.Type) {
		c.JSON(common.HTTP_OK, common.Error("'Type' Must one of [SIMPLE_LINE, COMPLEX_LINE]"))
		return
	}
	if vo.Type == "SIMPLE_LINE" {
		if len(vo.Schema) > 1 {
			c.JSON(common.HTTP_OK, common.Error("'SIMPLE_LINE' Type Only Can Have 1 Filed"))
			return
		}
	}

	MDataSchema := model.MDataSchema{
		UUID:   utils.DataSchemaUuid(),
		Name:   vo.Name,
		Type:   vo.Type,
		Schema: string(bytes),
	}
	if err := service.InsertDataSchema(MDataSchema); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 缓冲进内存
	core.SchemaSet(MDataSchema.UUID, typex.DataSchema{
		UUID:   MDataSchema.UUID,
		Name:   MDataSchema.Name,
		Type:   MDataSchema.Type,
		Schema: vo.Schema,
	})
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 更新模型
*
 */
func UpdateDataSchema(c *gin.Context, ruleEngine typex.RuleX) {
	vo := typex.DataSchema{}

	if err := c.ShouldBindJSON(&vo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	bytes, err := json.Marshal(vo.Schema)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MDataSchema := model.MDataSchema{
		UUID:   vo.UUID,
		Name:   vo.Name,
		Type:   vo.Type,
		Schema: string(bytes),
	}
	if err := service.UpdateDataSchema(MDataSchema); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 缓冲进内存
	core.SchemaSet(MDataSchema.UUID, typex.DataSchema{
		UUID:   MDataSchema.UUID,
		Name:   MDataSchema.Name,
		Type:   MDataSchema.Type,
		Schema: vo.Schema,
	})
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 删除模型
*
 */
func DeleteDataSchema(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if err := service.DeleteDataSchema(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 同步删除内存
	core.SchemaDelete(uuid)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 模型列表
*
 */
func ListDataSchema(c *gin.Context, ruleEngine typex.RuleX) {
	DataSchemas := []typex.DataSchema{}
	for _, vv := range service.AllDataSchema() {
		dataDefine := []typex.DataDefine{}
		err := json.Unmarshal([]byte(vv.Schema), &dataDefine)
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		DataSchemas = append(DataSchemas, typex.DataSchema{
			UUID:   vv.UUID,
			Name:   vv.Name,
			Type:   vv.Type,
			Schema: dataDefine,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(DataSchemas))

}

/*
*
* 模型详情
*
 */
func DataSchemaDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mDataSchema, err := service.GetDataSchemaWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	dataDefine := []typex.DataDefine{}
	err1 := json.Unmarshal([]byte(mDataSchema.Schema), &dataDefine)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(
		typex.DataSchema{
			UUID:   mDataSchema.UUID,
			Name:   mDataSchema.Name,
			Type:   mDataSchema.Type,
			Schema: dataDefine,
		},
	),
	)
}

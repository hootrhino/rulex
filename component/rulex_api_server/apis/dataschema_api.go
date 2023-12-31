package apis

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/iotschema"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 新建模型
*
 */
type IoTSchemaVo struct {
	UUID   string              `json:"uuid,omitempty"`
	Name   string              `json:"name"`
	Schema iotschema.IoTSchema `json:"schema"`
}

/*
*
* 新建一个物模型表
*
 */
func CreateDataSchema(c *gin.Context, ruleEngine typex.RuleX) {

	IoTSchema := IoTSchemaVo{}
	if err := c.ShouldBindJSON(&IoTSchema); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, IoTProperty := range IoTSchema.Schema.IoTProperties {
		if err := IoTProperty.ValidateType(); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	bytes, err := json.Marshal(IoTSchema)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	MIotSchema := model.MIotSchema{
		UUID:   utils.DataSchemaUuid(),
		Name:   IoTSchema.Name,
		Schema: string(bytes),
	}
	if err := service.InsertDataSchema(MIotSchema); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 更新模型
*
 */
func UpdateDataSchema(c *gin.Context, ruleEngine typex.RuleX) {
	IoTSchema := IoTSchemaVo{}
	if err := c.ShouldBindJSON(&IoTSchema); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	for _, IoTProperty := range IoTSchema.Schema.IoTProperties {
		if err := IoTProperty.ValidateType(); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	bytes, err := json.Marshal(IoTSchema)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MIotSchema := model.MIotSchema{
		UUID:   IoTSchema.UUID,
		Name:   IoTSchema.Name,
		Schema: string(bytes),
	}
	if err := service.UpdateDataSchema(MIotSchema); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
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
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 模型列表
*
 */
func ListDataSchema(c *gin.Context, ruleEngine typex.RuleX) {
	DataSchemas := []IoTSchemaVo{}
	for _, vv := range service.AllDataSchema() {
		IoTSchemaVo := IoTSchemaVo{UUID: vv.UUID, Name: vv.Name}
		err := json.Unmarshal([]byte(vv.Schema), &IoTSchemaVo)
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		DataSchemas = append(DataSchemas, IoTSchemaVo)
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
	MIotSchema, err := service.GetDataSchemaWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	IoTSchemaVo := IoTSchemaVo{UUID: MIotSchema.UUID, Name: MIotSchema.Name}
	err1 := json.Unmarshal([]byte(MIotSchema.Schema), &IoTSchemaVo)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(IoTSchemaVo))
}

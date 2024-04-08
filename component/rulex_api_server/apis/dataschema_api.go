package apis

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/interdb"
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
	UUID        string `json:"uuid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
*
* 属性 @ component/iotschema/iot_schema_define
*
 */
type IotPropertyVo struct {
	UUID        string            `json:"uuid"`        // UUID
	SchemaId    string            `json:"schemaId"`    //模型ID
	Label       string            `json:"label"`       // UI显示的那个文本
	Name        string            `json:"name"`        // 变量关联名
	Description string            `json:"description"` // 额外信息
	Type        string            `json:"type"`        // 类型, 只能是上面几种
	Rw          string            `json:"rw"`          // R读 W写 RW读写
	Unit        string            `json:"unit"`        // 单位 例如：摄氏度、米、牛等等
	Rule        IoTPropertyRuleVo `json:"rule"`        // 规则,IoTPropertyRule
}
type IoTPropertyRuleVo struct {
	DefaultValue any    `json:"defaultValue"` // 默认值
	Max          *int   `json:"max"`          // 最大值
	Min          *int   `json:"min"`          // 最小值
	TrueLabel    string `json:"trueLabel"`    // 真值label
	FalseLabel   string `json:"falseLabel"`   // 假值label
	Round        *int   `json:"round"`        // 小数点位
}

/*
*
* 属性
*
 */
func (O IoTPropertyRuleVo) String() string {
	if O.Max == nil {
		O.Max = new(int)
	}
	if O.Min == nil {
		O.Min = new(int)
	}
	if O.Round == nil {
		O.Round = new(int)
	}
	if O.DefaultValue == nil {
		O.DefaultValue = ""
	}
	if bytes, err := json.Marshal(O); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}

/*
*
* 从数据库反解析
*
 */
func (O IoTPropertyRuleVo) IoTPropertyRuleFromString(s string) error {
	if err := json.Unmarshal([]byte(s), &O); err != nil {
		return err
	}
	return nil
}

/*
*
* 从数据库保存的String字符串反解析规则
*
 */
func (P *IoTPropertyRuleVo) ParseRuleFromModel(s string) error {
	if P.Max == nil {
		P.Max = new(int)
	}
	if P.Min == nil {
		P.Min = new(int)
	}
	if P.Round == nil {
		P.Round = new(int)
	}

	if P.DefaultValue == nil {
		P.DefaultValue = ""
	}

	if err := json.Unmarshal([]byte(s), &P); err != nil {
		return err
	} else {
		return nil
	}
}

/*
*
* 新建一个物模型表
*
 */
func CreateDataSchema(c *gin.Context, ruleEngine typex.RuleX) {

	IoTSchemaVo := IoTSchemaVo{}
	if err := c.ShouldBindJSON(&IoTSchemaVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MIotSchema := model.MIotSchema{
		UUID:        utils.DataSchemaUuid(),
		Name:        IoTSchemaVo.Name,
		Description: IoTSchemaVo.Description,
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
	IoTSchemaVo := IoTSchemaVo{}
	if err := c.ShouldBindJSON(&IoTSchemaVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MIotSchema := model.MIotSchema{
		UUID:        IoTSchemaVo.UUID,
		Name:        IoTSchemaVo.Name,
		Description: IoTSchemaVo.Description,
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
	if err := service.DeleteDataSchemaAndProperty(uuid); err != nil {
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
		IoTSchemaVo := IoTSchemaVo{
			UUID:        vv.UUID,
			Name:        vv.Name,
			Description: vv.Description,
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
	IoTSchemaVo := IoTSchemaVo{
		UUID:        MIotSchema.UUID,
		Name:        MIotSchema.Name,
		Description: MIotSchema.Description,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(IoTSchemaVo))
}

// 分页获取
func CreateIotSchemaProperty(c *gin.Context, ruleEngine typex.RuleX) {
	IotPropertyVo := IotPropertyVo{}
	if err := c.ShouldBindJSON(&IotPropertyVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Schema, err := service.GetDataSchemaWithUUID(IotPropertyVo.SchemaId)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 不允许重复name
	count := service.CountIotSchemaProperty(IotPropertyVo.Name, Schema.UUID)
	if count > 0 {
		c.JSON(common.HTTP_OK, common.Error("Already Exists Property:"+IotPropertyVo.Name))
		return
	}
	err2 := service.InsertIotSchemaProperty(model.MIotProperty{
		SchemaId:    IotPropertyVo.SchemaId,
		UUID:        utils.MakeUUID("PROPER"),
		Label:       IotPropertyVo.Label,
		Name:        IotPropertyVo.Name,
		Description: IotPropertyVo.Description,
		Type:        IotPropertyVo.Type,
		Rw:          IotPropertyVo.Rw,
		Unit:        IotPropertyVo.Unit,
		Rule:        IotPropertyVo.Rule.String(), // 规则
	})
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 更新属性
func UpdateIotSchemaProperty(c *gin.Context, ruleEngine typex.RuleX) {
	IotPropertyVo := IotPropertyVo{}
	if err := c.ShouldBindJSON(&IotPropertyVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	_, err := service.GetDataSchemaWithUUID(IotPropertyVo.SchemaId)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	err2 := service.UpdateIotSchemaProperty(model.MIotProperty{
		SchemaId:    IotPropertyVo.SchemaId,
		UUID:        IotPropertyVo.UUID,
		Label:       IotPropertyVo.Label,
		Name:        IotPropertyVo.Name,
		Description: IotPropertyVo.Description,
		Type:        IotPropertyVo.Type,
		Rw:          IotPropertyVo.Rw,
		Unit:        IotPropertyVo.Unit,
		Rule:        IotPropertyVo.Rule.String(), // 规则
	})
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 删除属性
func DeleteIotSchemaProperty(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	err := service.DeleteIotSchemaProperty(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 分页查找数据
*
 */
func IotSchemaPropertyDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	record, err := service.FindIotSchemaProperty(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	IotPropertyVo := IotPropertyVo{
		SchemaId:    record.SchemaId,
		UUID:        record.UUID,
		Label:       record.Label,
		Name:        record.Name,
		Description: record.Description,
		Type:        record.Type,
		Rw:          record.Rw,
		Unit:        record.Unit,
	}
	IoTPropertyRuleVo := IoTPropertyRuleVo{}
	if err0 := IoTPropertyRuleVo.ParseRuleFromModel(record.Rule); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	IotPropertyVo.Rule = IoTPropertyRuleVo
	c.JSON(common.HTTP_OK, common.OkWithData(IotPropertyVo))
}

/*
*
* 列表
*
 */
func IotSchemaPropertyPageList(c *gin.Context, ruleEngine typex.RuleX) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	schemaUuid, _ := c.GetQuery("schema_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err1 := interdb.DB().Model(&model.MIotProperty{}).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	var records []model.MIotProperty
	result := tx.Order("created_at DESC").Find(&records,
		&model.MIotProperty{SchemaId: schemaUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	recordsVoList := []IotPropertyVo{}
	for _, record := range records {
		IoTPropertyRuleVo := IoTPropertyRuleVo{}
		if err0 := IoTPropertyRuleVo.ParseRuleFromModel(record.Rule); err0 != nil {
			c.JSON(common.HTTP_OK, common.Error400(err0))
			return
		}
		IotPropertyVo := IotPropertyVo{
			SchemaId:    record.SchemaId,
			UUID:        record.UUID,
			Label:       record.Label,
			Name:        record.Name,
			Description: record.Description,
			Type:        record.Type,
			Rw:          record.Rw,
			Unit:        record.Unit,
		}
		IotPropertyVo.Rule = IoTPropertyRuleVo
		recordsVoList = append(recordsVoList, IotPropertyVo)
	}
	Result := service.WrapPageResult(*pager, recordsVoList, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

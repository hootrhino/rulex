package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/server"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/square/go-jose.v2/json"
)

func InEndDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Model, err := service.GetMInEndWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	inEnd := ruleEngine.GetInEnd(Model.UUID)
	if inEnd == nil {
		tmpInEnd := typex.InEnd{
			UUID:        Model.UUID,
			Type:        typex.InEndType(Model.Type),
			Name:        Model.Name,
			Description: Model.Description,
			BindRules:   map[string]typex.Rule{},
			Config:      Model.GetConfig(),
			State:       typex.SOURCE_STOP,
		}
		c.JSON(common.HTTP_OK, common.OkWithData(tmpInEnd))
		return
	}
	inEnd.State = inEnd.Source.Status()
	c.JSON(common.HTTP_OK, common.OkWithData(inEnd))
}

// Get all inends
func InEnds(c *gin.Context, ruleEngine typex.RuleX) {

	inEnds := []typex.InEnd{}
	for _, v := range service.AllMInEnd() {
		var inEnd *typex.InEnd
		if inEnd = ruleEngine.GetInEnd(v.UUID); inEnd == nil {
			tmpInEnd := typex.InEnd{
				UUID:        v.UUID,
				Type:        typex.InEndType(v.Type),
				Name:        v.Name,
				Description: v.Description,
				BindRules:   map[string]typex.Rule{},
				Config:      v.GetConfig(),
				State:       typex.SOURCE_STOP,
			}
			inEnds = append(inEnds, tmpInEnd)
		}
		if inEnd != nil {
			inEnd.State = inEnd.Source.Status()
			inEnds = append(inEnds, *inEnd)
		}
	}
	c.JSON(common.HTTP_OK, common.OkWithData(inEnds))
}

// Create or Update InEnd
func CreateInend(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建, 非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err0 := c.ShouldBindJSON(&form); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	configJson, err1 := json.Marshal(form.Config)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	isSingle := false
	// 内部消息总线是单例模式
	if form.Type == typex.INTERNAL_EVENT.String() {
		ruleEngine.AllInEnd().Range(func(key, value any) bool {
			In := value.(*typex.InEnd)
			if In.Type.String() == form.Type {
				isSingle = true
				return false
			}
			return true
		})
	}
	if isSingle {
		msg := fmt.Errorf("the %s is singleton Source, can not create again", form.Name)
		c.JSON(common.HTTP_OK, common.Error400(msg))
		return
	}
	newUUID := utils.InUuid()

	if err := service.InsertMInEnd(&model.MInEnd{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
		XDataModels: "[]",
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := server.LoadNewestInEnd(newUUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg(err.Error()))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}
func RestartInEnd(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	err := ruleEngine.RestartInEnd(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 更新输入资源
*
 */
func UpdateInend(c *gin.Context, ruleEngine typex.RuleX) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建, 非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if form.UUID == "" {
		c.JSON(common.HTTP_OK, common.Error("missing 'uuid' fields"))
		return
	}
	// 更新的时候从数据库往外面拿
	InEnd, err := service.GetMInEndWithUUID(form.UUID)
	if err != nil {
		c.JSON(common.HTTP_OK, err)
		return
	}

	if err := service.UpdateMInEnd(InEnd.UUID, &model.MInEnd{
		UUID:        form.UUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := server.LoadNewestInEnd(form.UUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

// Delete inend by UUID
func DeleteInEnd(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := service.GetMInEndWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.DeleteMInEnd(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		old := ruleEngine.GetInEnd(uuid)
		if old != nil {
			old.Source.Stop()
			old.Source.Details().State = typex.SOURCE_STOP
		}
		ruleEngine.RemoveInEnd(uuid)
		c.JSON(common.HTTP_OK, common.Ok())
	}
}

/*
*
* UI配置表
*
 */
func GetInEndConfig(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithEmpty())
}

/*
*
* 属性表
*
 */
func GetInEndModels(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	inend := ruleEngine.GetInEnd(uuid)
	if inend != nil {
		modelsMap := inend.Source.DataModels()
		c.JSON(common.HTTP_OK, common.OkWithData(modelsMap))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithEmpty())
	}

}

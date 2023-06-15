package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gopkg.in/square/go-jose.v2/json"
)

func InEndDetail(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Model, err := hs.GetMInEndWithUUID(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	inEnd := e.GetInEnd(Model.UUID)
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
		c.JSON(HTTP_OK, OkWithData(tmpInEnd))
		return
	}
	c.JSON(HTTP_OK, OkWithData(inEnd))
}

// Get all inends
func InEnds(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		inEnds := []typex.InEnd{}
		for _, v := range hs.AllMInEnd() {
			var device *typex.InEnd
			if device = e.GetInEnd(v.UUID); device == nil {
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
			if device != nil {
				inEnds = append(inEnds, *device)
			}
		}
		c.JSON(HTTP_OK, OkWithData(inEnds))
		return
	}
	Model, err := hs.GetMInEndWithUUID(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	inEnd := e.GetInEnd(Model.UUID)
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
		c.JSON(HTTP_OK, OkWithData(tmpInEnd))
		return
	}
	c.JSON(HTTP_OK, OkWithData(inEnd))

}

// Create or Update InEnd
func CreateInend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建, 非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err0 := c.ShouldBindJSON(&form); err0 != nil {
		c.JSON(HTTP_OK, Error400(err0))
		return
	}
	configJson, err1 := json.Marshal(form.Config)
	if err1 != nil {
		c.JSON(HTTP_OK, Error400(err1))
		return
	}

	newUUID := utils.InUuid()

	if err := hh.InsertMInEnd(&MInEnd{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
		XDataModels: "[]",
	}); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if err := hh.LoadNewestInEnd(newUUID); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	c.JSON(HTTP_OK, Ok())

}

/*
*
* 更新输入资源
*
 */
func UpdateInend(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建, 非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if form.UUID == "" {
		c.JSON(HTTP_OK, Error("missing 'uuid' fields"))
		return
	}
	// 更新的时候从数据库往外面拿
	InEnd, err := hs.GetMInEndWithUUID(form.UUID)
	if err != nil {
		c.JSON(HTTP_OK, err)
		return
	}

	if err := hs.UpdateMInEnd(InEnd.UUID, &MInEnd{
		UUID:        form.UUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if err := hs.LoadNewestInEnd(form.UUID); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	c.JSON(HTTP_OK, Ok())
}

// Delete inend by UUID
func DeleteInEnd(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hh.GetMInEnd(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if err := hh.DeleteMInEnd(uuid); err != nil {
		c.JSON(HTTP_OK, Error400(err))
	} else {
		e.RemoveInEnd(uuid)
		c.JSON(HTTP_OK, Ok())
	}

}

/*
*
* UI配置表
*
 */
func GetInEndConfig(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	inend := e.GetInEnd(uuid)
	if inend != nil {
		c.JSON(HTTP_OK, OkWithData(inend.Source.Configs()))
	} else {
		c.JSON(HTTP_OK, OkWithEmpty())
	}

}

/*
*
* 属性表
*
 */
func GetInEndModels(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	inend := e.GetInEnd(uuid)
	if inend != nil {
		modelsMap := inend.Source.DataModels()
		c.JSON(HTTP_OK, OkWithData(modelsMap))
	} else {
		c.JSON(HTTP_OK, OkWithEmpty())
	}

}

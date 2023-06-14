package httpserver

import (
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

// Get all outends
func OutEnds(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		outends := []*typex.OutEnd{}
		for _, v := range hs.AllMOutEnd() {
			var outEnd *typex.OutEnd
			if outEnd = e.GetOutEnd(v.UUID); outEnd == nil {
				outEnd.State = typex.SOURCE_STOP
			}
			if outEnd != nil {
				outends = append(outends, outEnd)
			}
		}
		c.JSON(HTTP_OK, OkWithData(outends))
	} else {
		Model, err := hs.GetMOutEndWithUUID(uuid)
		if err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
		var outend *typex.OutEnd
		if outend = e.GetOutEnd(Model.UUID); outend == nil {
			outend.State = typex.SOURCE_STOP
		}
		c.JSON(HTTP_OK, OkWithData(outend))
	}
}

// Delete outEnd by UUID
func DeleteOutEnd(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hh.GetMOutEnd(uuid)
	if err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if err := hh.DeleteMOutEnd(uuid); err != nil {
		c.JSON(HTTP_OK, Error400(err))
	} else {
		e.RemoveOutEnd(uuid)
		c.JSON(HTTP_OK, Ok())
	}

}

// Create or Update OutEnd
func CreateOutEnd(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
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
	newUUID := utils.OutUuid()
	if err := hh.InsertMOutEnd(&MOutEnd{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if err := hh.LoadNewestOutEnd(newUUID); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	c.JSON(HTTP_OK, Ok())

}

// 更新
func UpdateOutEnd(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
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
	OutEnd, err := hs.GetMOutEndWithUUID(form.UUID)
	if err != nil {
		c.JSON(HTTP_OK, err)
		return
	}

	if err := hs.UpdateMOutEnd(OutEnd.UUID, &MOutEnd{
		UUID:        form.UUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	if err := hs.LoadNewestOutEnd(form.UUID); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	c.JSON(HTTP_OK, Ok())
}

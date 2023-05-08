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
			var device *typex.OutEnd
			if device = e.GetOutEnd(v.UUID); device == nil {
				device.State = typex.SOURCE_STOP
			}
			if device != nil {
				outends = append(outends, device)
			}
		}
		c.JSON(200, OkWithData(outends))
	} else {
		Model, err := hs.GetMOutEndWithUUID(uuid)
		if err != nil {
			c.JSON(200, Error400(err))
			return
		}
		var outend *typex.OutEnd
		if outend = e.GetOutEnd(Model.UUID); outend == nil {
			outend.State = typex.SOURCE_STOP
		}
		c.JSON(200, OkWithData(outend))
	}
}

// Delete outEnd by UUID
func DeleteOutEnd(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hh.GetMOutEnd(uuid)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	if err := hh.DeleteMOutEnd(uuid); err != nil {
		c.JSON(200, Error400(err))
	} else {
		e.RemoveOutEnd(uuid)
		c.JSON(200, Ok())
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
		c.JSON(200, Error400(err0))
		return
	}
	configJson, err1 := json.Marshal(form.Config)
	if err1 != nil {
		c.JSON(200, Error400(err1))
		return
	}
	var uuid *string = new(string)

	if form.UUID == "" {
		newUUID := utils.OutUuid()
		if err := hh.InsertMOutEnd(&MOutEnd{
			UUID:        newUUID,
			Type:        form.Type,
			Name:        form.Name,
			Description: form.Description,
			Config:      string(configJson),
		}); err != nil {
			c.JSON(200, Error400(err))
			return
		} else {
			uuid = &newUUID
		}
	} else {
		outEnd := e.GetOutEnd(form.UUID)
		if outEnd != nil {
			outEnd.SetState(typex.SOURCE_DOWN)
			outEnd.Target.Reload() // 重启
			hh.DeleteMOutEnd(outEnd.UUID)
			if err := hh.InsertMOutEnd(&MOutEnd{
				UUID:        form.UUID,
				Type:        form.Type,
				Name:        form.Name,
				Description: form.Description,
				Config:      string(configJson),
			}); err != nil {
				c.JSON(200, Error400(err))
				return
			}
			uuid = &form.UUID
		}
	}

	if err := hh.LoadNewestOutEnd(*uuid); err != nil {
		c.JSON(200, Error400(err))
		return
	} else {
		c.JSON(200, Ok())
		return
	}

}

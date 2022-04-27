package httpserver

import (
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

//
// Get all outends
//
func OutEnds(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []interface{}{}
		outEnds := e.AllOutEnd()
		outEnds.Range(func(key, value interface{}) bool {
			data = append(data, value)
			return true
		})
		c.JSON(200, Result{
			Code: 200,
			Msg:  SUCCESS,
			Data: data,
		})
	} else {

		c.JSON(200, Result{
			Code: 200,
			Msg:  SUCCESS,
			Data: e.GetOutEnd(uuid),
		})
	}

}

//
// Delete outend by UUID
//
func DeleteOutend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
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

//
// Create or Update OutEnd
//
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
		outend := e.GetOutEnd(form.UUID)
		if outend != nil {
			outend.SetState(typex.DOWN)
			outend.Target.Reload() // 重启
			hh.DeleteMOutEnd(outend.UUID)
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

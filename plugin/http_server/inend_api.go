package httpserver

import (
	"net/http"
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
	"github.com/ngaut/log"
	"gopkg.in/square/go-jose.v2/json"
)

//
// Get all inends
//
func InEnds(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []interface{}{}
		allInEnds := e.AllInEnd()
		allInEnds.Range(func(key, value interface{}) bool {
			data = append(data, value)
			return true
		})
		c.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  SUCCESS,
			Data: data,
		})
	} else {
		c.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  SUCCESS,
			Data: e.GetInEnd(uuid),
		})
	}
}

//
// Create or Update InEnd
//
func CreateInend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string                 `json:"uuid"` // 如果空串就是新建，非空就是更新
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	configJson, err1 := json.Marshal(form.Config)
	if err1 != nil {
		c.JSON(200, Error400(err1))
		return
	}
	var uuid *string = new(string)
	if form.UUID == "" {
		newUUID := utils.InUuid()
		if err := hh.InsertMInEnd(&MInEnd{
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
		inend := e.GetInEnd(form.UUID)
		if inend != nil {
			inend.Source.Reload() //重启接口
			inend.SetState(typex.DOWN)
			hh.DeleteMInEnd(inend.UUID)
			if err := hh.InsertMInEnd(&MInEnd{
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

	if err := hh.LoadNewestInEnd(*uuid); err != nil {
		log.Error(err)
		c.JSON(200, Error400(err))
		return
	} else {
		c.JSON(200, Ok())
		return
	}

}

//
// Delete inend by UUID
//
func DeleteInend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hh.GetMInEnd(uuid)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	if err := hh.DeleteMInEnd(uuid); err != nil {
		c.JSON(200, Error400(err))
	} else {
		e.RemoveInEnd(uuid)
		c.JSON(200, Ok())
	}

}

/*
*
* GetInEndConfig
*
 */
func GetInEndConfig(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	inend, ok := e.AllInEnd().Load(uuid)
	if ok {
		c.JSON(200, (inend.(*typex.InEnd)).Source.Configs())
	} else {
		c.JSON(400, []interface{}{})
	}

}

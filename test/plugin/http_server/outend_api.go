package httpserver

import (
	"net/http"
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

//
// Get all outends
//
func OutEnds(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	for _, v := range e.AllOutEnd() {
		data = append(data, v)
	}
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}

//
// Delete outend by UUID
//
func DeleteOutend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hh.GetMOutEnd(uuid)
	if err != nil {
		c.JSON(200, gin.H{"msg": err.Error()})
	} else {
		if err := hh.DeleteMOutEnd(uuid); err != nil {
			e.RemoveOutEnd(uuid)
			c.JSON(200, gin.H{"msg": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"msg": "remove success"})
		}
	}
}

//
// Create OutEnd
//
func CreateOutEnd(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Type        string                 `json:"type" binding:"required"`
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Config      map[string]interface{} `json:"config" binding:"required"`
	}
	form := Form{}
	err0 := c.ShouldBindJSON(&form)
	if err0 != nil {
		c.JSON(200, gin.H{"msg": err0.Error()})
		return
	} else {
		configJson, err1 := json.Marshal(form.Config)
		if err1 != nil {
			c.JSON(200, gin.H{"msg": err1.Error()})
			return
		} else {
			uuid := utils.MakeUUID("OUTEND")
			hh.InsertMOutEnd(&MOutEnd{
				UUID:        uuid,
				Type:        form.Type,
				Name:        form.Name,
				Description: form.Description,
				Config:      string(configJson),
			})
			err := hh.LoadNewestOutEnd(uuid)
			if err != nil {
				c.JSON(200, Error400(err))
				return
			} else {
				c.JSON(200, Ok())
				return
			}
		}
	}
}

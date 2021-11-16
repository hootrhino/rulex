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
	data := []interface{}{}
	allInEnds := e.AllInEnd()
	allInEnds.Range(func(key, value interface{}) bool {
		data = append(data, value)
		return true
	})

	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}

//
// Create InEnd
//
func CreateInend(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
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
	uuid := utils.MakeUUID("INEND")
	hh.InsertMInEnd(&MInEnd{
		UUID:        uuid,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	})
	if err := hh.LoadNewestInEnd(uuid); err != nil {
		log.Error(err)
		c.JSON(200, Error400(err))
		return
	} else {
		c.JSON(http.StatusOK, Ok())
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

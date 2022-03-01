package httpserver

import (
	"net/http"
	"rulex/typex"
	"rulex/utils"

	"github.com/emirpasic/gods/maps/linkedhashmap"
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
		c.JSON(http.StatusOK, OkWithData(data))
	} else {
		c.JSON(http.StatusOK, OkWithData(e.GetInEnd(uuid)))

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
		DataModels  []typex.XDataModel     `json:"dataModels"`
	}
	form := Form{}

	if err0 := c.ShouldBindJSON(&form); err0 != nil {
		c.JSON(http.StatusOK, Error400(err0))
		return
	}
	configJson, err1 := json.Marshal(form.Config)
	if err1 != nil {
		c.JSON(http.StatusOK, Error400(err1))
		return
	}
	//
	// 把数据模型表加工成MAP结构来缩短查询时间
	//
	// [{k1:v1}, {k2:v2}] --> {k1 :{k1:v1}, k2 :{k2:v2}}
	//
	var dataModelsMap = linkedhashmap.New()
	for _, v := range form.DataModels {
		dataModelsMap.Put(v.Name, v)
	}
	dataModelsJson, err2 := dataModelsMap.ToJSON()
	if err1 != nil {
		c.JSON(http.StatusOK, Error400(err2))
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
			DataModels:  string(dataModelsJson),
		}); err != nil {
			c.JSON(http.StatusOK, Error400(err))
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
				DataModels:  string(dataModelsJson),
			}); err != nil {
				c.JSON(http.StatusOK, Error400(err))
				return
			}
			uuid = &form.UUID
		}
	}

	if err := hh.LoadNewestInEnd(*uuid); err != nil {
		c.JSON(http.StatusOK, Error400(err))
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
		c.JSON(http.StatusOK, Error400(err))
		return
	}
	if err := hh.DeleteMInEnd(uuid); err != nil {
		c.JSON(http.StatusOK, Error400(err))
	} else {
		e.RemoveInEnd(uuid)
		c.JSON(http.StatusOK, Ok())
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
		c.JSON(http.StatusOK, OkWithData(inend.Source.Configs()))
	} else {
		c.JSON(http.StatusOK, OkWithEmpty())
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
		modelsMap := inend.DataModelsMap
		models := make([]typex.XDataModel, 0)
		for _, v := range modelsMap {
			models = append(models, v)
		}
		c.JSON(http.StatusOK, OkWithData(models))
	} else {
		c.JSON(http.StatusOK, OkWithEmpty())
	}

}

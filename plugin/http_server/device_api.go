package httpserver

import (
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

//
// 获取设备列表
//
func Devices(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []interface{}{}
		Devices := e.AllDevices()
		Devices.Range(func(key, value interface{}) bool {
			data = append(data, value)
			return true
		})
		c.JSON(200, OkWithData(data))
	} else {
		c.JSON(200, OkWithData(e.GetDevice(uuid)))
	}

}

//
// 删除设备
//
func DeleteDevice(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err := hs.GetDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	if err := hs.DeleteDevice(uuid); err != nil {
		c.JSON(200, Error400(err))
	} else {
		e.RemoveDevice(uuid)
		c.JSON(200, Ok())
	}

}

//
// 创建设备
//
func CreateDevice(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID         string                 `json:"uuid"`
		Name         string                 `json:"name"`
		Type         string                 `json:"type"`
		ActionScript string                 `json:"actionScript"`
		Config       map[string]interface{} `json:"config"`
		Description  string                 `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	newUUID := utils.DeviceUuid()
	if err := hs.InsertDevice(&MDevice{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	if err := hs.LoadNewestDevice(newUUID); err != nil {
		c.JSON(200, Error400(err))
		return
	} else {
		c.JSON(200, Ok())
		return
	}

}

//
// 更新设备
//
func UpdateDevice(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID         string                 `json:"uuid"`
		Name         string                 `json:"name"`
		Type         string                 `json:"type"`
		ActionScript string                 `json:"actionScript"`
		Config       map[string]interface{} `json:"config"`
		Description  string                 `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	if form.UUID == "" {
		c.JSON(200, Error("'uuid'参数缺失"))
		return
	}
	Device := e.GetDevice(form.UUID)
	if Device == nil {
		c.JSON(200, Error("设备不存在"))
		return
	}

	Device.Device.Stop()
	if err := hs.UpdateDevice(Device.UUID, &MDevice{
		UUID:        form.UUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}); err != nil {
		c.JSON(200, Error400(err))
		return
	}

	if err := hs.LoadNewestDevice(form.UUID); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	c.JSON(200, Ok())
}

package httpserver

import (
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

/*
*
* 列表先读数据库，然后读内存，合并状态后输出
*
 */
func Devices(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		devices := []*typex.Device{}
		for _, mdev := range hs.AllDevices() {
			device := e.GetDevice(mdev.UUID)
			if device == nil {
				tDevice := new(typex.Device)
				tDevice.UUID = mdev.UUID
				tDevice.Name = mdev.Name
				tDevice.Type = typex.DeviceType(mdev.Type)
				tDevice.Description = mdev.Description
				tDevice.BindRules = map[string]typex.Rule{}
				tDevice.Config = map[string]interface{}{}
				tDevice.State = typex.DEV_STOP
				devices = append(devices, tDevice)
			}
			if device != nil {
				devices = append(devices, device)
			}
		}
		c.JSON(200, OkWithData(devices))
	} else {
		mdev, err := hs.GetDeviceWithUUID(uuid)
		if err != nil {
			c.JSON(200, Error400(err))
			return
		}
		device := e.GetDevice(mdev.UUID)
		if device == nil {
			// 如果内存里面没有就给安排一个死设备
			tDevice := new(typex.Device)
			tDevice.UUID = mdev.UUID
			tDevice.Name = mdev.Name
			tDevice.Type = typex.DeviceType(mdev.Type)
			tDevice.Description = mdev.Description
			tDevice.BindRules = map[string]typex.Rule{}
			tDevice.Config = mdev.GetConfig()
			tDevice.State = typex.DEV_STOP
			c.JSON(200, OkWithData(tDevice))
			return
		}
		c.JSON(200, OkWithData(device))
	}
}

// 删除设备
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

// 创建设备
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

// 更新设备
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
		c.JSON(200, Error("missing 'uuid' fields"))
		return
	}
	Device := e.GetDevice(form.UUID)
	if Device == nil {
		c.JSON(200, Error("device not exists:"+form.UUID))
		return
	}

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

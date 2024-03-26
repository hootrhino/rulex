package apis

import (
	"fmt"

	"github.com/hootrhino/rulex/component/interdb"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/server"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"gorm.io/gorm"

	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2/json"
)

type DeviceVo struct {
	UUID        string                 `json:"uuid"`
	Gid         string                 `json:"gid"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	SchemaId    string                 `json:"schemaId"`
	State       int                    `json:"state"`
	ErrMsg      string                 `json:"errMsg"`
	Config      map[string]interface{} `json:"config"`
	Description string                 `json:"description"`
}

/*
*
* 列表先读数据库，然后读内存，合并状态后输出
*
 */
func DeviceDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mdev, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	DeviceVo := DeviceVo{}
	DeviceVo.UUID = mdev.UUID
	DeviceVo.Name = mdev.Name
	DeviceVo.Type = mdev.Type
	DeviceVo.SchemaId = mdev.SchemaId
	DeviceVo.Description = mdev.Description
	DeviceVo.Config = mdev.GetConfig()
	//
	device := ruleEngine.GetDevice(mdev.UUID)
	if device == nil {
		DeviceVo.State = int(typex.DEV_STOP)
	} else {
		DeviceVo.State = int(device.Device.Status())
	}
	Group := service.GetVisualGroup(mdev.UUID)
	DeviceVo.Gid = Group.UUID
	c.JSON(common.HTTP_OK, common.OkWithData(DeviceVo))
}

/*
*
* 分组查看
*
 */
func ListDeviceByGroup(c *gin.Context, ruleEngine typex.RuleX) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Gid, _ := c.GetQuery("uuid")
	count, MDevices := service.PageDeviceByGroup(pager.Current, pager.Size, Gid)
	err1 := interdb.DB().Model(&model.MDevice{}).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	devices := []DeviceVo{}
	for _, mdev := range MDevices {
		DeviceVo := DeviceVo{}
		DeviceVo.UUID = mdev.UUID
		DeviceVo.Name = mdev.Name
		DeviceVo.Type = mdev.Type
		DeviceVo.SchemaId = mdev.SchemaId
		DeviceVo.Description = mdev.Description
		DeviceVo.Config = mdev.GetConfig()
		//
		device := ruleEngine.GetDevice(mdev.UUID)
		if device == nil {
			DeviceVo.State = int(typex.DEV_STOP)
		} else {
			DeviceVo.State = int(device.Device.Status())
		}
		Group := service.GetVisualGroup(mdev.UUID)
		DeviceVo.Gid = Group.UUID

		devices = append(devices, DeviceVo)
	}

	Result := service.WrapPageResult(*pager, devices, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 重启
func RestartDevice(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	err := ruleEngine.RestartDevice(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 删除设备
func DeleteDevice(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	Mdev, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 检查是否有规则被绑定了
	for _, ruleId := range Mdev.BindRules {
		if ruleId != "" {
			_, err0 := service.GetMRuleWithUUID(ruleId)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
			c.JSON(common.HTTP_OK, common.Error("Can't remove, Already have rule bind:"+Mdev.BindRules.String()))
			return
		}

	}

	// 检查是否通用Modbus设备.需要同步删除点位表记录
	if Mdev.Type == typex.GENERIC_MODBUS.String() {
		if err := service.DeleteAllModbusPointByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	// 西门子需要同步删除点位表记录
	if Mdev.Type == typex.SIEMENS_PLC.String() {
		if err := service.DeleteAllSiemensPointByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	// 华中数控需要同步删除点位表记录
	if Mdev.Type == typex.HNC8.String() {
		if err := service.DeleteAllHnc8PointByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	old := ruleEngine.GetDevice(uuid)
	if old != nil {
		if old.Device.Status() == typex.DEV_UP {
			old.Device.Stop()
		}
	}
	// 事务
	txErr := interdb.DB().Transaction(func(tx *gorm.DB) error {
		Group := service.GetVisualGroup(uuid)
		err3 := service.DeleteDevice(uuid)
		if err3 != nil {
			return err3
		}
		// 解除关联
		err2 := interdb.DB().Where("gid=? and rid =?", Group.UUID, uuid).
			Delete(&model.MGenericGroupRelation{}).Error
		if err2 != nil {
			return err2
		}
		ruleEngine.RemoveDevice(uuid)
		return nil
	})
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

// 创建设备
func CreateDevice(c *gin.Context, ruleEngine typex.RuleX) {
	form := DeviceVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if service.CheckNameDuplicate(form.Name) {
		c.JSON(common.HTTP_OK, common.Error("Device Name Duplicated"))
		return
	}
	if utils.IsValidName(form.Name) {
		c.JSON(common.HTTP_OK, common.Error("Device Name Invalid, Must Between 6-12 characters"))
		return
	}
	isSingle := false
	// 红外线是单例模式
	if form.Type == typex.INTERNAL_EVENT.String() {
		ruleEngine.AllDevices().Range(func(key, value any) bool {
			In := value.(*typex.Device)
			if In.Type.String() == form.Type {
				isSingle = true
				return false
			}
			return true
		})
	}
	if isSingle {
		msg := fmt.Errorf("the %s is singleton Device, can not create again", form.Name)
		c.JSON(common.HTTP_OK, common.Error400(msg))
		return
	}
	newUUID := utils.DeviceUuid()
	MDevice := model.MDevice{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
		SchemaId:    form.SchemaId,
		BindRules:   []string{},
	}
	if err := service.InsertDevice(&MDevice); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 新建大屏的时候必须给一个分组
	if err := service.BindResource(form.Gid, MDevice.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error("Group not found"))
		return
	}
	if err := server.LoadNewestDevice(newUUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg(err.Error()))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 更新设备
func UpdateDevice(c *gin.Context, ruleEngine typex.RuleX) {

	form := DeviceVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if utils.IsValidName(form.Name) {
		c.JSON(common.HTTP_OK, common.Error("Device Name Invalid, Must Between 6-12 characters"))
		return
	}
	//
	// 取消绑定分组,删除原来旧的分组
	txErr := service.ReBindResource(func(tx *gorm.DB) error {
		MDevice := model.MDevice{
			Type:        form.Type,
			Name:        form.Name,
			SchemaId:    form.SchemaId,
			Description: form.Description,
			Config:      string(configJson),
		}
		return tx.Model(MDevice).
			Where("uuid=?", form.UUID).
			Updates(&MDevice).Error
	}, form.UUID, form.Gid)
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	if err := server.LoadNewestDevice(form.UUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取设备挂了的异常信息
*
 */
func GetDeviceErrorMsg(c *gin.Context, ruleEngine typex.RuleX) {

	c.JSON(common.HTTP_OK, common.OkWithData("Error Msg Not Found"))
}

/*
*
* 获取设备点位表挂了的异常信息
*
 */
func GetDevicePointErrorMsg(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData("Error Msg Not Found"))
}

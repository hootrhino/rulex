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
	State       int                    `json:"state"`
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
	Gid, _ := c.GetQuery("uuid")
	devices := []DeviceVo{}
	MDevices := service.FindDeviceByGroup(Gid)
	for _, mdev := range MDevices {
		DeviceVo := DeviceVo{}
		DeviceVo.UUID = mdev.UUID
		DeviceVo.Name = mdev.Name
		DeviceVo.Type = mdev.Type
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
	c.JSON(common.HTTP_OK, common.OkWithData(devices))
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
	// 要处理这个空字符串 ""
	if Mdev.BindRules.Len() == 1 && len(Mdev.BindRules[0]) != 0 {
		c.JSON(common.HTTP_OK, common.Error("Can't remove, Already have rule bind:"+Mdev.BindRules.String()))
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
		}

	}

	// 检查是否通用Modbus设备.需要同步删除点位表记录
	if Mdev.Type == "GENERIC_MODBUS" {
		if err := service.DeleteAllModbusPointByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	// 西门子的
	if Mdev.Type == "SIEMENS_PLC" {
		if err := service.DeleteAllSiemensPointByDevice(uuid); err != nil {
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
	err1 := service.DeleteDevice(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}

	ruleEngine.RemoveDevice(uuid)
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
	if form.UUID == "" {
		c.JSON(common.HTTP_OK, common.Error("missing 'uuid' fields"))
		return
	}
	// 更新的时候从数据库往外面拿
	Device, err := service.GetMDeviceWithUUID(form.UUID)
	if err != nil {
		c.JSON(common.HTTP_OK, err)
		return
	}
	MDevice := model.MDevice{
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}
	if err := service.UpdateDevice(Device.UUID, &MDevice); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	// 只有检查到分组变了才更改
	Group := service.GetVisualGroup(Device.UUID)
	if Group.UUID != form.Gid {
		// 取消绑定分组,删除原来旧的分组
		txErr := interdb.DB().Transaction(func(tx *gorm.DB) error {
			err1 := tx.Where("gid=? and rid =?", Group.UUID, Device.UUID).
				Delete(&model.MGenericGroupRelation{}).Error
			if err1 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err))
				return err1
			}
			// 重新绑定分组,首先确定分组是否存在
			MGroup := model.MGenericGroup{}
			if err2 := interdb.DB().Where("uuid=?", Group.UUID).First(&MGroup).Error; err2 != nil {
				return err2
			}
			Relation := model.MGenericGroupRelation{
				Gid: MGroup.UUID,
				Rid: Device.UUID,
			}
			err3 := tx.Save(&Relation).Error
			if err3 != nil {
				return err3
			}
			return nil
		})
		if txErr != nil {
			c.JSON(common.HTTP_OK, common.Error400(txErr))
			return
		}
	}
	if err := server.LoadNewestDevice(form.UUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.Ok())
}

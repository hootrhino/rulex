package service

import (
	"errors"

	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"

	"gorm.io/gorm"
)

// -----------------------------------------------------------------------------------
func GetMRule(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetAllMRule() ([]model.MRule, error) {
	m := []model.MRule{}
	if err := sqlitedao.Sqlite.DB().Find(&m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func GetMRuleWithUUID(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMRule(r *model.MRule) error {
	return sqlitedao.Sqlite.DB().Table("m_rules").Create(r).Error
}

func DeleteMRule(uuid string) error {
	if sqlitedao.Sqlite.DB().Table("m_rules").Where("uuid=?", uuid).Delete(&model.MRule{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func UpdateMRule(uuid string, r *model.MRule) error {
	m := model.MRule{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*r)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func GetMInEnd(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := sqlitedao.Sqlite.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetMInEndWithUUID(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := sqlitedao.Sqlite.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMInEnd(i *model.MInEnd) error {
	return sqlitedao.Sqlite.DB().Table("m_in_ends").Create(i).Error
}

func DeleteMInEnd(uuid string) error {
	if sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MInEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func UpdateMInEnd(uuid string, i *model.MInEnd) error {
	m := model.MInEnd{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*i)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func GetMOutEnd(id string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := sqlitedao.Sqlite.DB().First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetMOutEndWithUUID(uuid string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMOutEnd(o *model.MOutEnd) error {
	return sqlitedao.Sqlite.DB().Table("m_out_ends").Create(o).Error
}

func DeleteMOutEnd(uuid string) error {
	if sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MOutEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func UpdateMOutEnd(uuid string, o *model.MOutEnd) error {
	m := model.MOutEnd{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
// USER
// -----------------------------------------------------------------------------------
func GetMUser(username string, password string) (*model.MUser, error) {
	m := new(model.MUser)
	if err := sqlitedao.Sqlite.DB().Where("Username=?", username).Where("Password=?",
		password).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMUser(o *model.MUser) {
	sqlitedao.Sqlite.DB().Table("m_users").Create(o)
}

func UpdateMUser(uuid string, o *model.MUser) error {
	m := model.MUser{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func AllMRules() []model.MRule {
	rules := []model.MRule{}
	sqlitedao.Sqlite.DB().Table("m_rules").Find(&rules)
	return rules
}

func AllMInEnd() []model.MInEnd {
	inends := []model.MInEnd{}
	sqlitedao.Sqlite.DB().Table("m_in_ends").Find(&inends)
	return inends
}

func AllMOutEnd() []model.MOutEnd {
	outends := []model.MOutEnd{}
	sqlitedao.Sqlite.DB().Table("m_out_ends").Find(&outends)
	return outends
}

func AllMUser() []model.MUser {
	users := []model.MUser{}
	sqlitedao.Sqlite.DB().Find(&users)
	return users
}

func AllDevices() []model.MDevice {
	devices := []model.MDevice{}
	sqlitedao.Sqlite.DB().Find(&devices)
	return devices
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func GetMDeviceWithUUID(uuid string) (*model.MDevice, error) {
	m := new(model.MDevice)
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

// 删除设备
func DeleteDevice(uuid string) error {
	if sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MDevice{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建设备
func InsertDevice(o *model.MDevice) error {
	return sqlitedao.Sqlite.DB().Table("m_devices").Create(o).Error
}

// 更新设备信息
func UpdateDevice(uuid string, o *model.MDevice) error {
	m := model.MDevice{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*o)
		return nil
	}
}

// -------------------------------------------------------------------------------------
// ModbusPointPositions
// -------------------------------------------------------------------------------------

// InsertModbusPointPosition 插入modbus点位表
func InsertModbusPointPosition(list []model.MModbusPointPosition) error {
	m := model.MModbusPointPosition{}
	return sqlitedao.Sqlite.DB().Model(m).Create(list).Error
}

// DeleteModbusPointAndDevice 删除modbus点位与设备
func DeleteModbusPointAndDevice(deviceUuid string) error {
	return sqlitedao.Sqlite.DB().Transaction(func(tx *gorm.DB) (err error) {

		err = tx.Where("device_uuid = ?", deviceUuid).Delete(&model.MModbusPointPosition{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("uuid = ?", deviceUuid).Delete(&model.MDevice{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// UpdateModbusPoint 更新modbus点位
func UpdateModbusPoint(mm model.MModbusPointPosition) error {
	m := model.MDevice{}
	if err := sqlitedao.Sqlite.DB().Where("id = ?", mm.ID).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(&m)
		return nil
	}
}

// AllModbusPointByDeviceUuid 根据设备UUID查询设备点位
func AllModbusPointByDeviceUuid(deviceUuid string) (list []model.MModbusPointPosition, err error) {

	err = sqlitedao.Sqlite.DB().Where("device_uuid = ?", deviceUuid).Find(&list).Error
	return
}

// -------------------------------------------------------------------------------------
// Goods
// -------------------------------------------------------------------------------------

// 获取Goods列表
func AllGoods() []model.MGoods {
	m := []model.MGoods{}
	sqlitedao.Sqlite.DB().Find(&m)
	return m

}
func GetGoodsWithUUID(uuid string) (*model.MGoods, error) {
	m := model.MGoods{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Goods
func DeleteGoods(uuid string) error {
	if sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MGoods{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建Goods
func InsertGoods(goods *model.MGoods) error {
	return sqlitedao.Sqlite.DB().Table("m_goods").Create(goods).Error
}

// 更新Goods
func UpdateGoods(uuid string, goods *model.MGoods) error {
	m := model.MGoods{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*goods)
		return nil
	}
}

// -------------------------------------------------------------------------------------
// App Dao
// -------------------------------------------------------------------------------------

// 获取App列表
func AllApp() []model.MApp {
	m := []model.MApp{}
	sqlitedao.Sqlite.DB().Find(&m)
	return m

}
func GetMAppWithUUID(uuid string) (*model.MApp, error) {
	m := model.MApp{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func DeleteApp(uuid string) error {
	return sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MApp{}).Error
}

// 创建App
func InsertApp(app *model.MApp) error {
	return sqlitedao.Sqlite.DB().Create(app).Error
}

// 更新App
func UpdateApp(app *model.MApp) error {
	m := model.MApp{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", app.UUID).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*app)
		return nil
	}
}

// 获取AiBase列表
func AllAiBase() []model.MAiBase {
	m := []model.MAiBase{}
	sqlitedao.Sqlite.DB().Find(&m)
	return m

}
func GetAiBaseWithUUID(uuid string) (*model.MAiBase, error) {
	m := model.MAiBase{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除AiBase
func DeleteAiBase(uuid string) error {
	return sqlitedao.Sqlite.DB().Where("uuid=?", uuid).Delete(&model.MAiBase{}).Error
}

// 创建AiBase
func InsertAiBase(AiBase *model.MAiBase) error {
	return sqlitedao.Sqlite.DB().Create(AiBase).Error
}

// 更新AiBase
func UpdateAiBase(AiBase *model.MAiBase) error {
	m := model.MAiBase{}
	if err := sqlitedao.Sqlite.DB().Where("uuid=?", AiBase.UUID).First(&m).Error; err != nil {
		return err
	} else {
		sqlitedao.Sqlite.DB().Model(m).Updates(*AiBase)
		return nil
	}
}

// -------------------------------------------------------------------------------------
// Cron Task
// -------------------------------------------------------------------------------------

// AllEnabledCronTask
func AllEnabledCronTask() []model.MCronTask {
	tasks := make([]model.MCronTask, 0)
	sqlitedao.Sqlite.DB().Where("enable = ?", "1").Find(&tasks)
	return tasks
}

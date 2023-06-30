package httpserver

import (
	"errors"

	model "github.com/hootrhino/rulex/plugin/http_server/model"

	"gorm.io/gorm"
)

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMRule(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	if err := s.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetAllMRule() ([]model.MRule, error) {
	m := []model.MRule{}
	if err := s.DB().Find(&m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) GetMRuleWithUUID(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	if err := s.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMRule(r *model.MRule) error {
	return s.DB().Table("m_rules").Create(r).Error
}

func (s *HttpApiServer) DeleteMRule(uuid string) error {
	if s.DB().Table("m_rules").Where("uuid=?", uuid).Delete(&model.MRule{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMRule(uuid string, r *model.MRule) error {
	m := model.MRule{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*r)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMInEnd(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := s.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetMInEndWithUUID(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := s.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMInEnd(i *model.MInEnd) error {
	return s.DB().Table("m_in_ends").Create(i).Error
}

func (s *HttpApiServer) DeleteMInEnd(uuid string) error {
	if s.DB().Where("uuid=?", uuid).Delete(&model.MInEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMInEnd(uuid string, i *model.MInEnd) error {
	m := model.MInEnd{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*i)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMOutEnd(id string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := s.DB().First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetMOutEndWithUUID(uuid string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := s.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMOutEnd(o *model.MOutEnd) error {
	return s.DB().Table("m_out_ends").Create(o).Error
}

func (s *HttpApiServer) DeleteMOutEnd(uuid string) error {
	if s.DB().Where("uuid=?", uuid).Delete(&model.MOutEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMOutEnd(uuid string, o *model.MOutEnd) error {
	m := model.MOutEnd{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
// USER
// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMUser(username string, password string) (*model.MUser, error) {
	m := new(model.MUser)
	if err := s.DB().Where("Username=?", username).Where("Password=?",
		password).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMUser(o *model.MUser) {
	s.DB().Table("m_users").Create(o)
}

func (s *HttpApiServer) UpdateMUser(uuid string, o *model.MUser) error {
	m := model.MUser{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) AllMRules() []model.MRule {
	rules := []model.MRule{}
	s.DB().Table("m_rules").Find(&rules)
	return rules
}

func (s *HttpApiServer) AllMInEnd() []model.MInEnd {
	inends := []model.MInEnd{}
	s.DB().Table("m_in_ends").Find(&inends)
	return inends
}

func (s *HttpApiServer) AllMOutEnd() []model.MOutEnd {
	outends := []model.MOutEnd{}
	s.DB().Table("m_out_ends").Find(&outends)
	return outends
}

func (s *HttpApiServer) AllMUser() []model.MUser {
	users := []model.MUser{}
	s.DB().Find(&users)
	return users
}

func (s *HttpApiServer) AllDevices() []model.MDevice {
	devices := []model.MDevice{}
	s.DB().Find(&devices)
	return devices
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func (s *HttpApiServer) GetMDeviceWithUUID(uuid string) (*model.MDevice, error) {
	m := new(model.MDevice)
	if err := s.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

// 删除设备
func (s *HttpApiServer) DeleteDevice(uuid string) error {
	if s.DB().Where("uuid=?", uuid).Delete(&model.MDevice{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建设备
func (s *HttpApiServer) InsertDevice(o *model.MDevice) error {
	return s.DB().Table("m_devices").Create(o).Error
}

// 更新设备信息
func (s *HttpApiServer) UpdateDevice(uuid string, o *model.MDevice) error {
	m := model.MDevice{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*o)
		return nil
	}
}

// InsertModbusPointPosition 插入modbus点位表
func (s *HttpApiServer) InsertModbusPointPosition(list []model.MModbusPointPosition) error {
	m := model.MModbusPointPosition{}
	return s.DB().Model(m).Create(list).Error
}

// DeleteModbusPointAndDevice 删除modbus点位与设备
func (s *HttpApiServer) DeleteModbusPointAndDevice(deviceUuid string) error {
	return s.DB().Transaction(func(tx *gorm.DB) (err error) {

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

// -------------------------------------------------------------------------------------
// Goods
// -------------------------------------------------------------------------------------

// 获取Goods列表
func (s *HttpApiServer) AllGoods() []model.MGoods {
	m := []model.MGoods{}
	s.DB().Find(&m)
	return m

}
func (s *HttpApiServer) GetGoodsWithUUID(uuid string) (*model.MGoods, error) {
	m := model.MGoods{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Goods
func (s *HttpApiServer) DeleteGoods(uuid string) error {
	if s.DB().Where("uuid=?", uuid).Delete(&model.MGoods{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建Goods
func (s *HttpApiServer) InsertGoods(goods *model.MGoods) error {
	return s.DB().Table("m_goods").Create(goods).Error
}

// 更新Goods
func (s *HttpApiServer) UpdateGoods(uuid string, goods *model.MGoods) error {
	m := model.MGoods{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*goods)
		return nil
	}
}

// -------------------------------------------------------------------------------------
// App Dao
// -------------------------------------------------------------------------------------

// 获取App列表
func (s *HttpApiServer) AllApp() []model.MApp {
	m := []model.MApp{}
	s.DB().Find(&m)
	return m

}
func (s *HttpApiServer) GetMAppWithUUID(uuid string) (*model.MApp, error) {
	m := model.MApp{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func (s *HttpApiServer) DeleteApp(uuid string) error {
	return s.DB().Where("uuid=?", uuid).Delete(&model.MApp{}).Error
}

// 创建App
func (s *HttpApiServer) InsertApp(app *model.MApp) error {
	return s.DB().Create(app).Error
}

// 更新App
func (s *HttpApiServer) UpdateApp(app *model.MApp) error {
	m := model.MApp{}
	if err := s.DB().Where("uuid=?", app.UUID).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*app)
		return nil
	}
}

// 获取AiBase列表
func (s *HttpApiServer) AllAiBase() []model.MAiBase {
	m := []model.MAiBase{}
	s.DB().Find(&m)
	return m

}
func (s *HttpApiServer) GetAiBaseWithUUID(uuid string) (*model.MAiBase, error) {
	m := model.MAiBase{}
	if err := s.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除AiBase
func (s *HttpApiServer) DeleteAiBase(uuid string) error {
	return s.DB().Where("uuid=?", uuid).Delete(&model.MAiBase{}).Error
}

// 创建AiBase
func (s *HttpApiServer) InsertAiBase(AiBase *model.MAiBase) error {
	return s.DB().Create(AiBase).Error
}

// 更新AiBase
func (s *HttpApiServer) UpdateAiBase(AiBase *model.MAiBase) error {
	m := model.MAiBase{}
	if err := s.DB().Where("uuid=?", AiBase.UUID).First(&m).Error; err != nil {
		return err
	} else {
		s.DB().Model(m).Updates(*AiBase)
		return nil
	}
}

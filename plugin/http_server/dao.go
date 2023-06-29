package httpserver

import (
	"errors"
	"os"

	model "github.com/hootrhino/rulex/plugin/http_server/model"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
)

/*
*
* 初始化数据库
*
 */
func (s *HttpApiServer) InitDb(dbPath string) {
	var err error
	if core.GlobalConfig.AppDebugMode {
		s.sqliteDb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		s.sqliteDb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}

	if err != nil {
		glogger.GLogger.Error(err)
		// Sqlite 创建失败应该是致命错误了, 多半是环境出问题，直接给panic了, 不尝试救活
		panic(err)
	}
	// 注册数据库配置表
	// 这么写看起来是很难受, 但是这玩意就是go的哲学啊(大道至简？？？)
	if err := s.sqliteDb.AutoMigrate(
		&model.MInEnd{},
		&model.MOutEnd{},
		&model.MRule{},
		&model.MUser{},
		&model.MDevice{},
		&model.MGoods{},
		&model.MApp{},
		&model.MAiBase{},
		&model.MModbusPointPosition{},
	); err != nil {
		glogger.GLogger.Fatal(err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMRule(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetAllMRule() ([]model.MRule, error) {
	m := []model.MRule{}
	if err := s.sqliteDb.Find(&m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) GetMRuleWithUUID(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMRule(r *model.MRule) error {
	return s.sqliteDb.Table("m_rules").Create(r).Error
}

func (s *HttpApiServer) DeleteMRule(uuid string) error {
	if s.sqliteDb.Table("m_rules").Where("uuid=?", uuid).Delete(&model.MRule{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMRule(uuid string, r *model.MRule) error {
	m := model.MRule{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*r)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMInEnd(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := s.sqliteDb.Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetMInEndWithUUID(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := s.sqliteDb.Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMInEnd(i *model.MInEnd) error {
	return s.sqliteDb.Table("m_in_ends").Create(i).Error
}

func (s *HttpApiServer) DeleteMInEnd(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&model.MInEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMInEnd(uuid string, i *model.MInEnd) error {
	m := model.MInEnd{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*i)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMOutEnd(id string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := s.sqliteDb.First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetMOutEndWithUUID(uuid string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMOutEnd(o *model.MOutEnd) error {
	return s.sqliteDb.Table("m_out_ends").Create(o).Error
}

func (s *HttpApiServer) DeleteMOutEnd(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&model.MOutEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMOutEnd(uuid string, o *model.MOutEnd) error {
	m := model.MOutEnd{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
// USER
// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMUser(username string, password string) (*model.MUser, error) {
	m := new(model.MUser)
	if err := s.sqliteDb.Where("Username=?", username).Where("Password=?",
		password).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMUser(o *model.MUser) {
	s.sqliteDb.Table("m_users").Create(o)
}

func (s *HttpApiServer) UpdateMUser(uuid string, o *model.MUser) error {
	m := model.MUser{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) AllMRules() []model.MRule {
	rules := []model.MRule{}
	s.sqliteDb.Table("m_rules").Find(&rules)
	return rules
}

func (s *HttpApiServer) AllMInEnd() []model.MInEnd {
	inends := []model.MInEnd{}
	s.sqliteDb.Table("m_in_ends").Find(&inends)
	return inends
}

func (s *HttpApiServer) AllMOutEnd() []model.MOutEnd {
	outends := []model.MOutEnd{}
	s.sqliteDb.Table("m_out_ends").Find(&outends)
	return outends
}

func (s *HttpApiServer) AllMUser() []model.MUser {
	users := []model.MUser{}
	s.sqliteDb.Find(&users)
	return users
}

func (s *HttpApiServer) AllDevices() []model.MDevice {
	devices := []model.MDevice{}
	s.sqliteDb.Find(&devices)
	return devices
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func (s *HttpApiServer) GetMDeviceWithUUID(uuid string) (*model.MDevice, error) {
	m := new(model.MDevice)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

// 删除设备
func (s *HttpApiServer) DeleteDevice(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&model.MDevice{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建设备
func (s *HttpApiServer) InsertDevice(o *model.MDevice) error {
	return s.sqliteDb.Table("m_devices").Create(o).Error
}

// 更新设备信息
func (s *HttpApiServer) UpdateDevice(uuid string, o *model.MDevice) error {
	m := model.MDevice{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

// InsertModbusPointPosition 插入modbus点位表
func (s *HttpApiServer) InsertModbusPointPosition(list []model.MModbusPointPosition) error {
	m := model.MModbusPointPosition{}
	return s.sqliteDb.Model(m).Create(list).Error
}

// DeleteModbusPointAndDevice 删除modbus点位与设备
func (s *HttpApiServer) DeleteModbusPointAndDevice(deviceUuid string) error {
	return s.sqliteDb.Transaction(func(tx *gorm.DB) (err error) {

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
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetGoodsWithUUID(uuid string) (*model.MGoods, error) {
	m := model.MGoods{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Goods
func (s *HttpApiServer) DeleteGoods(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&model.MGoods{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建Goods
func (s *HttpApiServer) InsertGoods(goods *model.MGoods) error {
	return s.sqliteDb.Table("m_goods").Create(goods).Error
}

// 更新Goods
func (s *HttpApiServer) UpdateGoods(uuid string, goods *model.MGoods) error {
	m := model.MGoods{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*goods)
		return nil
	}
}

// -------------------------------------------------------------------------------------
// App Dao
// -------------------------------------------------------------------------------------

// 获取App列表
func (s *HttpApiServer) AllApp() []model.MApp {
	m := []model.MApp{}
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetAppWithUUID(uuid string) (*model.MApp, error) {
	m := model.MApp{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func (s *HttpApiServer) DeleteApp(uuid string) error {
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&model.MApp{}).Error
}

// 创建App
func (s *HttpApiServer) InsertApp(app *model.MApp) error {
	return s.sqliteDb.Create(app).Error
}

// 更新App
func (s *HttpApiServer) UpdateApp(app *model.MApp) error {
	m := model.MApp{}
	if err := s.sqliteDb.Where("uuid=?", app.UUID).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*app)
		return nil
	}
}

// 获取AiBase列表
func (s *HttpApiServer) AllAiBase() []model.MAiBase {
	m := []model.MAiBase{}
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetAiBaseWithUUID(uuid string) (*model.MAiBase, error) {
	m := model.MAiBase{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除AiBase
func (s *HttpApiServer) DeleteAiBase(uuid string) error {
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&model.MAiBase{}).Error
}

// 创建AiBase
func (s *HttpApiServer) InsertAiBase(AiBase *model.MAiBase) error {
	return s.sqliteDb.Create(AiBase).Error
}

// 更新AiBase
func (s *HttpApiServer) UpdateAiBase(AiBase *model.MAiBase) error {
	m := model.MAiBase{}
	if err := s.sqliteDb.Where("uuid=?", AiBase.UUID).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*AiBase)
		return nil
	}
}

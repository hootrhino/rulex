package httpserver

import (
	"errors"
	"os"

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
		&MInEnd{},
		&MOutEnd{},
		&MRule{},
		&MUser{},
		&MDevice{},
		&MGoods{},
		&MApp{},
		&MAiBase{},
		&MModbusPointPosition{},
	); err != nil {
		glogger.GLogger.Fatal(err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMRule(uuid string) (*MRule, error) {
	m := new(MRule)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetAllMRule() ([]MRule, error) {
	m := []MRule{}
	if err := s.sqliteDb.Find(&m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) GetMRuleWithUUID(uuid string) (*MRule, error) {
	m := new(MRule)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMRule(r *MRule) error {
	return s.sqliteDb.Table("m_rules").Create(r).Error
}

func (s *HttpApiServer) DeleteMRule(uuid string) error {
	if s.sqliteDb.Table("m_rules").Where("uuid=?", uuid).Delete(&MRule{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMRule(uuid string, r *MRule) error {
	m := MRule{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*r)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMInEnd(uuid string) (*MInEnd, error) {
	m := new(MInEnd)
	if err := s.sqliteDb.Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetMInEndWithUUID(uuid string) (*MInEnd, error) {
	m := new(MInEnd)
	if err := s.sqliteDb.Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMInEnd(i *MInEnd) error {
	return s.sqliteDb.Table("m_in_ends").Create(i).Error
}

func (s *HttpApiServer) DeleteMInEnd(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&MInEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMInEnd(uuid string, i *MInEnd) error {
	m := MInEnd{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*i)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMOutEnd(id string) (*MOutEnd, error) {
	m := new(MOutEnd)
	if err := s.sqliteDb.First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (s *HttpApiServer) GetMOutEndWithUUID(uuid string) (*MOutEnd, error) {
	m := new(MOutEnd)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMOutEnd(o *MOutEnd) error {
	return s.sqliteDb.Table("m_out_ends").Create(o).Error
}

func (s *HttpApiServer) DeleteMOutEnd(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&MOutEnd{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

func (s *HttpApiServer) UpdateMOutEnd(uuid string, o *MOutEnd) error {
	m := MOutEnd{}
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
func (s *HttpApiServer) GetMUser(username string, password string) (*MUser, error) {
	m := new(MUser)
	if err := s.sqliteDb.Where("Username=?", username).Where("Password=?",
		password).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (s *HttpApiServer) InsertMUser(o *MUser) {
	s.sqliteDb.Table("m_users").Create(o)
}

func (s *HttpApiServer) UpdateMUser(uuid string, o *MUser) error {
	m := MUser{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func (s *HttpApiServer) AllMRules() []MRule {
	rules := []MRule{}
	s.sqliteDb.Table("m_rules").Find(&rules)
	return rules
}

func (s *HttpApiServer) AllMInEnd() []MInEnd {
	inends := []MInEnd{}
	s.sqliteDb.Table("m_in_ends").Find(&inends)
	return inends
}

func (s *HttpApiServer) AllMOutEnd() []MOutEnd {
	outends := []MOutEnd{}
	s.sqliteDb.Table("m_out_ends").Find(&outends)
	return outends
}

func (s *HttpApiServer) AllMUser() []MUser {
	users := []MUser{}
	s.sqliteDb.Find(&users)
	return users
}

func (s *HttpApiServer) AllDevices() []MDevice {
	devices := []MDevice{}
	s.sqliteDb.Find(&devices)
	return devices
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func (s *HttpApiServer) GetMDeviceWithUUID(uuid string) (*MDevice, error) {
	m := new(MDevice)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

// 删除设备
func (s *HttpApiServer) DeleteDevice(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&MDevice{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建设备
func (s *HttpApiServer) InsertDevice(o *MDevice) error {
	return s.sqliteDb.Table("m_devices").Create(o).Error
}

// 更新设备信息
func (s *HttpApiServer) UpdateDevice(uuid string, o *MDevice) error {
	m := MDevice{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

// InsertModbusPointPosition 插入modbus点位表
func (s *HttpApiServer) InsertModbusPointPosition(list []MModbusPointPosition) error {
	m := MModbusPointPosition{}
	return s.sqliteDb.Model(m).Create(list).Error
}

// DeleteModbusPointAndDevice 删除modbus点位与设备
func (s *HttpApiServer) DeleteModbusPointAndDevice(deviceUuid string) error {
	return s.sqliteDb.Transaction(func(tx *gorm.DB) (err error) {

		err = tx.Where("device_uuid = ?", deviceUuid).Delete(&MModbusPointPosition{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("uuid = ?", deviceUuid).Delete(&MDevice{}).Error
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
func (s *HttpApiServer) AllGoods() []MGoods {
	m := []MGoods{}
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetGoodsWithUUID(uuid string) (*MGoods, error) {
	m := MGoods{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Goods
func (s *HttpApiServer) DeleteGoods(uuid string) error {
	if s.sqliteDb.Where("uuid=?", uuid).Delete(&MGoods{}).RowsAffected == 0 {
		return errors.New("not found:" + uuid)
	}
	return nil
}

// 创建Goods
func (s *HttpApiServer) InsertGoods(goods *MGoods) error {
	return s.sqliteDb.Table("m_goods").Create(goods).Error
}

// 更新Goods
func (s *HttpApiServer) UpdateGoods(uuid string, goods *MGoods) error {
	m := MGoods{}
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
func (s *HttpApiServer) AllApp() []MApp {
	m := []MApp{}
	s.sqliteDb.Find(&m)
	return m

}
func (s *HttpApiServer) GetAppWithUUID(uuid string) (*MApp, error) {
	m := MApp{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func (s *HttpApiServer) DeleteApp(uuid string) error {
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&MApp{}).Error
}

// 创建App
func (s *HttpApiServer) InsertApp(app *MApp) error {
	return s.sqliteDb.Create(app).Error
}

// 更新App
func (s *HttpApiServer) UpdateApp(app *MApp) error {
	m := MApp{}
	if err := s.sqliteDb.Where("uuid=?", app.UUID).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*app)
		return nil
	}
}

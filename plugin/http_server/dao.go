package httpserver

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
*
* 初始化数据库
*
 */
func (s *HttpApiServer) InitDb(dbPath string) {
	var err error
	s.sqliteDb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	var err1 error = nil
	err1 = s.sqliteDb.AutoMigrate(&MInEnd{})
	err1 = s.sqliteDb.AutoMigrate(&MOutEnd{})
	err1 = s.sqliteDb.AutoMigrate(&MRule{})
	err1 = s.sqliteDb.AutoMigrate(&MUser{})
	err1 = s.sqliteDb.AutoMigrate(&MDevice{})
	if err1 != nil {
		log.Error(err1)
		os.Exit(1)
	}
}

//-----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMRule(uuid string) (*MRule, error) {
	m := new(MRule)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
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
	return s.sqliteDb.Table("m_rules").Where("uuid=?", uuid).Delete(&MRule{}).Error
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

//-----------------------------------------------------------------------------------
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
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&MInEnd{}).Error
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

//-----------------------------------------------------------------------------------
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
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&MOutEnd{}).Error
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

//-----------------------------------------------------------------------------------
// USER
//-----------------------------------------------------------------------------------
func (s *HttpApiServer) GetMUser(username string, password string) (*MUser, error) {
	m := new(MUser)
	if err := s.sqliteDb.Where("Username=?", username).Where("Password=?", password).First(m).Error; err != nil {
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

//-----------------------------------------------------------------------------------
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
	s.sqliteDb.Find(&outends)
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

//-------------------------------------------------------------------------------------

//
// 获取设备列表
//
func (s *HttpApiServer) GetDeviceWithUUID(uuid string) (*MDevice, error) {
	m := new(MDevice)
	if err := s.sqliteDb.Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

//
// 删除设备
//
func (s *HttpApiServer) DeleteDevice(uuid string) error {
	return s.sqliteDb.Where("uuid=?", uuid).Delete(&MDevice{}).Error
}

//
// 创建设备
//
func (s *HttpApiServer) InsertDevice(o *MDevice) error {
	return s.sqliteDb.Table("m_devices").Create(o).Error
}

//
// 更新设备信息
//
func (s *HttpApiServer) UpdateDevice(uuid string, o *MDevice) error {
	m := MDevice{}
	if err := s.sqliteDb.Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		s.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

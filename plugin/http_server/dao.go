package httpserver

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (hh *HttpApiServer) InitDb() {
	var err error
	hh.sqliteDb, err = gorm.Open(sqlite.Open("./reluxdb.sqlite3"),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	hh.sqliteDb.AutoMigrate(&MInEnd{})
	hh.sqliteDb.AutoMigrate(&MRule{})
	hh.sqliteDb.AutoMigrate(&MOutEnd{})
	hh.sqliteDb.AutoMigrate(&MUser{})
}

//-----------------------------------------------------------------------------------
func (hh *HttpApiServer) GetMRule(id int) (*MRule, error) {
	m := new(MRule)
	if err := hh.sqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (hh *HttpApiServer) InsertMRule(r *MRule) {
	hh.sqliteDb.Table("m_rules").Create(r)
}

func (hh *HttpApiServer) UpdateMRule(id int, r *MRule) error {
	m := MRule{}
	if err := hh.sqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		hh.sqliteDb.Model(m).Updates(*r)
		return nil
	}
}

//-----------------------------------------------------------------------------------
func (hh *HttpApiServer) GetMInEnd(id int) (*MInEnd, error) {
	m := new(MInEnd)
	if err := hh.sqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (hh *HttpApiServer) InsertMInEnd(i *MInEnd) {
	hh.sqliteDb.Table("m_in_ends").Create(i)
}
func (hh *HttpApiServer) UpdateMInEnd(id int, i *MInEnd) error {
	m := MInEnd{}
	if err := hh.sqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		hh.sqliteDb.Model(m).Updates(*i)
		return nil
	}
}

//-----------------------------------------------------------------------------------
func (hh *HttpApiServer) GetMOutEnd(id int) (*MOutEnd, error) {
	m := new(MOutEnd)
	if err := hh.sqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (hh *HttpApiServer) InsertMOutEnd(o *MOutEnd) {
	hh.sqliteDb.Table("m_out_ends").Create(o)
	log.Debug("Create outend success")

}
func (hh *HttpApiServer) UpdateMOutEnd(id int, o *MOutEnd) error {
	m := MOutEnd{}
	if err := hh.sqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		hh.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

//-----------------------------------------------------------------------------------

//-----------------------------------------------------------------------------------
func (hh *HttpApiServer) GetMUser(id int) (*MUser, error) {
	m := new(MUser)
	if err := hh.sqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (hh *HttpApiServer) InsertMUser(o *MUser) {
	hh.sqliteDb.Table("m_users").Create(o)
	log.Debug("Create outend success")

}
func (hh *HttpApiServer) UpdateMUser(id int, o *MUser) error {
	m := MUser{}
	if err := hh.sqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		hh.sqliteDb.Model(m).Updates(*o)
		return nil
	}
}

//-----------------------------------------------------------------------------------
func (hh *HttpApiServer) AllMRules() []MRule {
	rules := []MRule{}
	hh.sqliteDb.Find(&rules)
	return rules
}
func (hh *HttpApiServer) AllMInEnd() []MInEnd {
	inends := []MInEnd{}
	hh.sqliteDb.Table("m_in_ends").Find(&inends)
	return inends
}
func (hh *HttpApiServer) AllMOutEnd() []MOutEnd {
	outends := []MOutEnd{}
	hh.sqliteDb.Find(&outends)
	return outends
}
func (hh *HttpApiServer) AllMUser() []MUser {
	users := []MUser{}
	hh.sqliteDb.Find(&users)
	return users
}

//-----------------------------------------------------------------------------------
func (hh *HttpApiServer) Truncate() error {
	log.Warn("This operation will truncate table. So ONLY for test or debug!!")
	hh.sqliteDb.Unscoped().Where("1 = 1").Delete(&MUser{})
	hh.sqliteDb.Unscoped().Where("1 = 1").Delete(&MInEnd{})
	hh.sqliteDb.Unscoped().Where("1 = 1").Delete(&MOutEnd{})
	hh.sqliteDb.Unscoped().Where("1 = 1").Delete(&MRule{})
	return nil
}

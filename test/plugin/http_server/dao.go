package httpserver

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (hh *HttpApiServer) InitDb() {
	var err error
	hh.sqliteDb, err = gorm.Open(sqlite.Open("./rulex.db"))
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

func (hh *HttpApiServer) GetMRuleWithUUID(uuid string) (*MRule, error) {
	m := new(MRule)
	if err := hh.sqliteDb.Where("UUID=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (hh *HttpApiServer) InsertMRule(r *MRule) {
	hh.sqliteDb.Table("m_rules").Create(r)
}

func (hh *HttpApiServer) DeleteMRule(uuid string) {
	hh.sqliteDb.Where("UUID=?", uuid).Unscoped().Delete(&MRule{})
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
	if err := hh.sqliteDb.Table("m_in_ends").Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func (hh *HttpApiServer) GetMInEndWithUUID(uuid string) (*MInEnd, error) {
	m := new(MInEnd)
	if err := hh.sqliteDb.Table("m_in_ends").Where("UUID=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (hh *HttpApiServer) InsertMInEnd(i *MInEnd) {
	hh.sqliteDb.Table("m_in_ends").Create(i)
}

func (hh *HttpApiServer) DeleteMInEnd(uuid string) {
	hh.sqliteDb.Where("UUID=?", uuid).Unscoped().Delete(&MInEnd{})
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
func (hh *HttpApiServer) GetMOutEndWithUUID(uuid string) (*MOutEnd, error) {
	m := new(MOutEnd)
	if err := hh.sqliteDb.Where("UUID=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (hh *HttpApiServer) InsertMOutEnd(o *MOutEnd) {
	hh.sqliteDb.Table("m_out_ends").Create(o)
}

func (hh *HttpApiServer) DeleteMOutEnd(uuid string) {
	hh.sqliteDb.Where("UUID=?", uuid).Unscoped().Delete(&MOutEnd{})
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
func (hh *HttpApiServer) GetMUser(username string, password string) (*MUser, error) {
	m := new(MUser)
	if err := hh.sqliteDb.Where("Username=?", username).Where("Password=?", password).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (hh *HttpApiServer) InsertMUser(o *MUser) {
	hh.sqliteDb.Table("m_users").Create(o)
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

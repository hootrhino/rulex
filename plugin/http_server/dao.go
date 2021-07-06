package httpserver

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/ngaut/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SqliteDb *gorm.DB

func init() {
	InitDb()
}

func InitDb() {
	var err error
	SqliteDb, err = gorm.Open(sqlite.Open("./reluxdb.sqlite3"),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	SqliteDb.AutoMigrate(&MInEnd{})
	SqliteDb.AutoMigrate(&MRule{})
	SqliteDb.AutoMigrate(&MOutEnd{})
	SqliteDb.AutoMigrate(&MUser{})
}

//-----------------------------------------------------------------------------------
func GetMRule(id int) (*MRule, error) {
	m := new(MRule)
	if err := SqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func InsertMRule(r *MRule) {
	SqliteDb.Create(r)
}

func UpdateMRule(id int, r *MRule) error {
	m := MRule{}
	if err := SqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		SqliteDb.Model(m).Updates(*r)
		return nil
	}
}

//-----------------------------------------------------------------------------------
func GetMInEnd(id int) (*MInEnd, error) {
	m := new(MInEnd)
	if err := SqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func InsertMInEnd(i *MInEnd) {
	SqliteDb.Create(i)
}
func UpdateMInEnd(id int, i *MInEnd) error {
	m := MInEnd{}
	if err := SqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		SqliteDb.Model(m).Updates(*i)
		return nil
	}
}

//-----------------------------------------------------------------------------------
func GetMOutEnd(id int) (*MOutEnd, error) {
	m := new(MOutEnd)
	if err := SqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func InsertMOutEnd(o *MOutEnd) {
	SqliteDb.Create(o)
	log.Debug("Create outend success")

}
func UpdateMOutEnd(id int, o *MOutEnd) error {
	m := MOutEnd{}
	if err := SqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		SqliteDb.Model(m).Updates(*o)
		return nil
	}
}

//-----------------------------------------------------------------------------------

//-----------------------------------------------------------------------------------
func GetMUser(id int) (*MUser, error) {
	m := new(MUser)
	if err := SqliteDb.Where("Id=?", id).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func InsertMUser(o *MUser) {
	SqliteDb.Create(o)
	log.Debug("Create outend success")

}
func UpdateMUser(id int, o *MUser) error {
	m := MUser{}
	if err := SqliteDb.Where("Id=?", id).First(&m).Error; err != nil {
		return err
	} else {
		SqliteDb.Model(m).Updates(*o)
		return nil
	}
}

//-----------------------------------------------------------------------------------
func AllMRules() []MRule {
	rules := []MRule{}
	SqliteDb.Find(&rules)
	return rules
}
func AllMInEnd() []MInEnd {
	inends := []MInEnd{}
	SqliteDb.Find(&inends)
	return inends
}
func AllMOutEnd() []MOutEnd {
	outends := []MOutEnd{}
	SqliteDb.Find(&outends)
	return outends
}
func AllMUser() []MUser {
	users := []MUser{}
	SqliteDb.Find(&users)
	return users
}

//-----------------------------------------------------------------------------------

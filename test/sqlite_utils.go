package test

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
*
* SQLite 查询工具 用来辅助单元测试
*
 */
var unitTestDB *gorm.DB

func LoadUnitTestDB() {
	var err error
	unitTestDB, err = gorm.Open(sqlite.Open("unitest.db"), &gorm.Config{})
	if err != nil {
		panic("failed to load unitest database")
	}
}

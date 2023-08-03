package initialize

import (
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/plugin/api_server/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

// Gorm 暂时写死sqlite,可根据实际扩展其他数据库
func Gorm(dbPath string) *gorm.DB {
	return GormSqlite(dbPath)
}

func GormSqlite(dbPath string) *gorm.DB {
	if core.GlobalConfig.AppDebugMode {
		if db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		}); err != nil {
			return nil
		} else {
			return db
		}
	} else {
		if db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		}); err != nil {
			return nil
		} else {
			return db
		}
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.MRule{},
	)
	if err != nil {
		glogger.GLogger.Fatal(err)
		os.Exit(0)
	}
}

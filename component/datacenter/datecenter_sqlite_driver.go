package datacenter

import (
	"runtime"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"

	"github.com/hootrhino/rulex/glogger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const __DEFAULT_DB_PATH string = "./INTERNAL_DATACENTER.db"

var __Sqlite *SqliteDAO

/*
*
* Sqlite 数据持久层
*
 */
type SqliteDAO struct {
	engine typex.RuleX
	name   string   // 框架可以根据名称来选择不同的数据库驱动,为以后扩展准备
	db     *gorm.DB // Sqlite 驱动
}

/*
*
* 初始化DAO
*
 */
func InitSqliteDAO(engine typex.RuleX) *SqliteDAO {
	__Sqlite = &SqliteDAO{name: "Sqlite3", engine: engine}

	var err error
	if core.GlobalConfig.AppDebugMode {
		__Sqlite.db, err = gorm.Open(sqlite.Open(__DEFAULT_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: false,
		})
	} else {
		__Sqlite.db, err = gorm.Open(sqlite.Open(__DEFAULT_DB_PATH), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Error),
			SkipDefaultTransaction: false,
		})
	}
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	return __Sqlite
}

/*
*
* 停止
*
 */
func Stop() {
	__Sqlite.db = nil
	runtime.GC()
}

/*
*
* 返回数据库查询句柄
*
 */
func DB() *gorm.DB {
	return __Sqlite.db
}

/*
*
* 返回名称
*
 */
func Name() string {
	return __Sqlite.name
}

/*
*
* 注册数据模型
*
 */
func RegisterModel(dist ...interface{}) {
	__Sqlite.db.AutoMigrate(dist...)
}

package glogger

import (
	"time"

	"github.com/sirupsen/logrus"
)

var LUA_LOGGER *LogWriter
var GLOBAL_LOGGER *LogWriter

/*
*
* 配置全局logging记录器
*
 */

var GLogger *logrus.Logger

func StartGLogger(path string) {
	GLOBAL_LOGGER = NewLogWriter("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
	GLogger = logrus.New()
	GLogger.Formatter = new(logrus.JSONFormatter)
	GLogger.Formatter.(*logrus.JSONFormatter).PrettyPrint = true
	GLogger.SetOutput(GLOBAL_LOGGER)
}

/*
*
* StartLuaLogger
*
 */
func StartLuaLogger(path string) {
	LUA_LOGGER = NewLogWriter("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
}

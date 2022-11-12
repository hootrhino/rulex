package glogger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var private_lua_logger *LogWriter
var private_local_logger *LogWriter
var private_remote_logger *UdpLogger

type LogMsg struct {
	Sn      string `json:"sn"`
	Uid     string `json:"uid"`
	Content string `json:"content"`
}

/*
*
* 配置全局logging记录器
*
 */

var GLogger *logrus.Logger = logrus.New()

func StartGLogger(EnableConsole bool, path string) {
	private_local_logger = NewLogWriter("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
	GLogger.Formatter = new(logrus.JSONFormatter)
	GLogger.SetReportCaller(true)
	GLogger.Formatter.(*logrus.JSONFormatter).PrettyPrint = true
	if EnableConsole {
		GLogger.SetOutput(os.Stdout)
	} else {
		GLogger.SetOutput(private_local_logger)
	}
}

/*
*
* StartLuaLogger
*
 */
func StartLuaLogger(path string) {
	private_lua_logger = NewLogWriter("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
}

/*
*
* LUA 脚本的日志接口
*
 */
func Log(b []byte) {
	private_lua_logger.Write(b)
}

/*
*
* 关闭日志
*
 */
func Close() error {
	if err := private_local_logger.Close(); err != nil {
		return err
	}
	if err := private_lua_logger.Close(); err != nil {
		return err
	}
	return nil
}

package glogger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var private_local_logger *LogWriter

/*
*
* 配置全局logging记录器
*
 */

var GLogger *logrus.Logger = logrus.New()

func StartGLogger(EnableConsole bool, path string) {
	private_local_logger = NewLogWriter("./" + time.Now().Format("2006-01-02_15-04-05-") + path)
	GLogger.Formatter = new(logrus.JSONFormatter)
	GLogger.SetReportCaller(true)
	// GLogger.Formatter.(*logrus.JSONFormatter).PrettyPrint = true
	if EnableConsole {
		GLogger.SetOutput(os.Stdout)
	} else {
		GLogger.SetOutput(private_local_logger)
	}
}

/*
*
* 关闭日志记录器
*
 */
func Close() error {
	if err := private_local_logger.Close(); err != nil {
		return err
	}
	if err := private_lua_logger.Close(); err != nil {
		return err
	}
	if err := private_remote_logger.Close(); err != nil {
		return err
	}
	return nil
}

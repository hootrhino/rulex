package glogger

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var private_local_logger *LogWriter

/*
*
* 配置全局logging记录器
*
 */

var Logrus *logrus.Logger = logrus.New()
var GLogger *logrus.Entry

func StartGLogger(LogLevel string,
	EnableConsole bool,
	AppDebugMode bool,
	path string,
	key string, value interface{}) {
	GLogger = Logrus.WithField("appId", value)
	private_local_logger = NewLogWriter("./" + time.Now().Format("2006-01-02-") + path)
	Logrus.Formatter = new(logrus.JSONFormatter)
	if AppDebugMode {
		Logrus.SetReportCaller(true)
	}
	if EnableConsole {
		Logrus.SetOutput(os.Stdout)
	} else {
		Logrus.SetOutput(private_local_logger)
	}
	setLogLevel(LogLevel)
}
func setLogLevel(LogLevel string) {
	switch LogLevel {
	case "fatal":
		Logrus.SetLevel(logrus.FatalLevel)
	case "error":
		Logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		Logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		Logrus.SetLevel(logrus.DebugLevel)
	case "info":
		Logrus.SetLevel(logrus.InfoLevel)
	case "all", "trace":
		Logrus.SetLevel(logrus.TraceLevel)
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
	return nil
}

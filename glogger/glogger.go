package glogger

import (
	"os"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

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
	Logrus.Formatter = new(logrus.JSONFormatter)
	if AppDebugMode {
		Logrus.SetReportCaller(true)
	}
	if EnableConsole {
		Logrus.SetOutput(os.Stdout)
	} else {
		Logrus.SetOutput(&lumberjack.Logger{
			Filename:   path + ".txt",
			MaxSize:    5,    // 超过5Mb备份
			MaxBackups: 2,    // 最多备份2次
			MaxAge:     7,    // 最大保留天数
			Compress:   true, // 压缩备份
		})
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
	return nil
}

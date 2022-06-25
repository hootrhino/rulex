package core

import (
	"net/http"
	"os"
	"runtime"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/sirupsen/logrus"

	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RulexConfig
var INIPath string

//
// Init config
//
func InitGlobalConfig(path string) typex.RulexConfig {
	glogger.GLogger.Info("Init rulex config")
	cfg, err := ini.Load(path)
	if err != nil {
		glogger.GLogger.Fatalf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	INIPath = path
	//---------------------------------------
	if err := cfg.Section("app").MapTo(&GlobalConfig); err != nil {
		glogger.GLogger.Fatalf("Fail to map config file: %v", err)
		os.Exit(1)
	}
	glogger.GLogger.Info("Rulex config init successfully")
	return GlobalConfig
}

func SetLogLevel() {
	switch GlobalConfig.LogLevel {
	case "fatal":
		glogger.GLogger.SetLevel(logrus.FatalLevel)
	case "error":
		glogger.GLogger.SetLevel(logrus.ErrorLevel)
	case "warn":
		glogger.GLogger.SetLevel(logrus.WarnLevel)
	case "debug":
		glogger.GLogger.SetLevel(logrus.DebugLevel)
	case "info":
		glogger.GLogger.SetLevel(logrus.InfoLevel)
	}

}

/*
*
* 设置性能，通常用来Debug用，生产环境建议关闭
*
 */
func SetPerformance() {
	if GlobalConfig.GomaxProcs > 0 {
		if GlobalConfig.GomaxProcs < runtime.NumCPU() {
			runtime.GOMAXPROCS(GlobalConfig.GomaxProcs)
		} else {
			glogger.GLogger.Warnf("GomaxProcs is %v, but current CPU number is:%v", GlobalConfig.GomaxProcs, runtime.NumCPU())
		}
	}
	//------------------------------------------------------
	// pprof: https://segmentfault.com/a/1190000016412013
	//------------------------------------------------------
	if GlobalConfig.EnablePProf {
		glogger.GLogger.Debug("Start PProf debug at: 0.0.0.0:6060")
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		runtime.SetCPUProfileRate(1)
		go http.ListenAndServe("0.0.0.0:6060", nil)
	}
}

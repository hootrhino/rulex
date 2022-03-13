package core

import (
	"encoding/json"
	"net/http"
	"os"
	"rulex/typex"
	"runtime"

	"github.com/ngaut/log"
	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RulexConfig
var INIPath string

//
// Init config
//
func InitGlobalConfig(path string) typex.RulexConfig {
	log.Info("Init rulex config")
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	INIPath = path
	//---------------------------------------
	if err := cfg.Section("app").MapTo(&GlobalConfig); err != nil {
		log.Fatalf("Fail to map config file: %v", err)
		os.Exit(1)
	}
	log.Info("Rulex config init successfully")
	bytes, err := json.MarshalIndent(GlobalConfig, " ", "  ")
	if err != nil {
		log.Fatalf("Fail to marshal config file: %v", err)
		os.Exit(1)
	} else {
		log.Info(string(bytes))
	}
	return GlobalConfig
}

func SetLogLevel() {
	log.SetHighlighting(false)
	switch GlobalConfig.LogLevel {
	case "fatal":
		log.SetLevel(log.LogLevel(log.LOG_FATAL))
	case "error":
		log.SetLevel(log.LogLevel(log.LOG_ERROR))
	case "warn":
		log.SetLevel(log.LogLevel(log.LOG_LEVEL_WARN))
	case "warning":
		log.SetLevel(log.LogLevel(log.LOG_WARNING))
	case "debug":
		log.SetLevel(log.LogLevel(log.LOG_DEBUG))
	case "info":
		log.SetLevel(log.LogLevel(log.LOG_INFO))
	case "all":
		log.SetLevel(log.LogLevel(log.LOG_LEVEL_ALL))
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
			log.Warnf("GomaxProcs is %v, but current CPU number is:%v", GlobalConfig.GomaxProcs, runtime.NumCPU())
		}
	}
	//------------------------------------------------------
	// pprof: https://segmentfault.com/a/1190000016412013
	//------------------------------------------------------
	if GlobalConfig.EnablePProf {
		log.Debug("Start PProf debug at: 0.0.0.0:6060")
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		runtime.SetCPUProfileRate(1)
		go http.ListenAndServe("0.0.0.0:6060", nil)
	}
}

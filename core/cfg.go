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

func SetPerformance() {
	log.Info("Go max process is:", GlobalConfig.GomaxProcs)
	runtime.GOMAXPROCS(GlobalConfig.GomaxProcs)
	//------------------------------------------------------
	// pprof: https://segmentfault.com/a/1190000016412013
	//------------------------------------------------------
	if GlobalConfig.EnablePProf {
		log.Debug("Start PProf debug")
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		runtime.SetCPUProfileRate(1)
		go http.ListenAndServe("0.0.0.0:6060", nil)
	}
}

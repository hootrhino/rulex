package core

import (
	"net/http"
	"os"
	"rulex/typex"
	"runtime"

	"github.com/ngaut/log"
	"gopkg.in/ini.v1"
)

var TM typex.TargetRegistry
var RM typex.ResourceRegistry

//
// Global config
//
type RulexConfig struct {
	MaxQueueSize            int    `json:"maxQueueSize"`
	ResourceRestartInterval int    `json:"resourceRestartInterval"`
	GomaxProcs              int    `json:"gomaxProcs"`
	EnablePProf             bool   `json:"enablePProf"`
	LogLevel                string `json:"logLevel"`
	LogPath                 string `json:"logPath"`
	LuaLogPath              string `json:"luaLogPath"`
}

var GlobalConfig RulexConfig

//
// Init config
//
func InitGlobalConfig() {
	log.Info("Init rulex config")
	cfg, err := ini.Load("conf/rulex.ini")
	if err != nil {
		log.Fatalf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	TM = NewTargetTypeManager()
	RM = NewResourceTypeManager()
	//---------------------------------------
	GlobalConfig.MaxQueueSize = cfg.Section("app").Key("max_queue_size").MustInt(5000)
	GlobalConfig.ResourceRestartInterval = cfg.Section("app").Key("resource_restart_interval").MustInt(204800)
	GlobalConfig.GomaxProcs = cfg.Section("app").Key("gomax_procs").MustInt(2)
	GlobalConfig.EnablePProf = cfg.Section("app").Key("enable_pprof").MustBool(false)
	GlobalConfig.LogLevel = cfg.Section("app").Key("log_level").MustString("info")
	GlobalConfig.LogPath = cfg.Section("app").Key("log_path").MustString("./rulex-log.txt")
	GlobalConfig.LuaLogPath = cfg.Section("app").Key("lua_log_path").MustString("./rulex-lua-log.txt")

	log.Info("Rulex config init successfully")

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

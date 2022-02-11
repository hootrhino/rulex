package core

import (
	"net/http"
	"os"
	"rulex/typex"
	"runtime"

	"github.com/ngaut/log"
	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RulexConfig

//
// Init config
//
func InitGlobalConfig() typex.RulexConfig {
	log.Info("Init rulex config")
	cfg, err := ini.Load("conf/rulex.ini")
	if err != nil {
		log.Fatalf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	//---------------------------------------
	GlobalConfig.MaxQueueSize = cfg.Section("app").Key("max_queue_size").MustInt(5000)
	log.Info("| MaxQueueSize is:", GlobalConfig.MaxQueueSize)
	GlobalConfig.SourceRestartInterval = cfg.Section("app").Key("source_restart_interval").MustInt(204800)
	log.Info("| SourceRestartInterval is:", GlobalConfig.SourceRestartInterval)
	GlobalConfig.GomaxProcs = cfg.Section("app").Key("gomax_procs").MustInt(2)
	log.Info("| GomaxProcs is:", GlobalConfig.GomaxProcs)
	GlobalConfig.EnablePProf = cfg.Section("app").Key("enable_pprof").MustBool(false)
	log.Info("| EnablePProf is:", GlobalConfig.EnablePProf)
	GlobalConfig.LogLevel = cfg.Section("app").Key("log_level").MustString("info")
	log.Info("| LogLevel is:", GlobalConfig.LogLevel)
	GlobalConfig.LogPath = cfg.Section("app").Key("log_path").MustString("./rulex-log.txt")
	log.Info("| LogPath is:", GlobalConfig.LogPath)
	GlobalConfig.LuaLogPath = cfg.Section("app").Key("lua_log_path").MustString("./rulex-lua-log.txt")
	log.Info("| LuaLogPath is:", GlobalConfig.LuaLogPath)

	log.Info("Rulex config init successfully")

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

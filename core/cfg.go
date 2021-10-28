package core

import (
	"net/http"
	"os"
	"runtime"

	"github.com/ngaut/log"
	"gopkg.in/ini.v1"
)

//
// Global config
//
type RulexConfig struct {
	Name                    string
	Path                    string
	Token                   string
	Secret                  string
	MaxQueueSize            int
	ResourceRestartInterval int
	GomaxProcs              int
	EnablePProf             bool
	LogLevel                string
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
	//---------------------------------------
	GlobalConfig.Name = cfg.Section("app").Key("name").MustString("rulex")
	GlobalConfig.MaxQueueSize = cfg.Section("app").Key("max_queue_size").MustInt(5000)
	GlobalConfig.ResourceRestartInterval = cfg.Section("app").Key("resource_restart_interval").MustInt(204800)
	GlobalConfig.GomaxProcs = cfg.Section("app").Key("gomax_procs").MustInt(2)
	GlobalConfig.EnablePProf = cfg.Section("app").Key("enable_pprof").MustBool(false)
	GlobalConfig.LogLevel = cfg.Section("app").Key("log_level").MustString("info")
	//---------------------------------------
	GlobalConfig.Path = cfg.Section("cloud").Key("path").MustString("")
	GlobalConfig.Token = cfg.Section("cloud").Key("token").MustString("")
	GlobalConfig.Secret = cfg.Section("cloud").Key("secret").MustString("")
	log.Info("Rulex config init successfully")

}
func SetLogLevel() {
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

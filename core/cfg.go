package core

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/i4de/rulex/typex"

	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RulexConfig
var INIPath string

// Init config
func InitGlobalConfig(path string) typex.RulexConfig {
	log.Println("Init rulex config")
	cfg, err := ini.ShadowLoad(path)
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
	if err := cfg.Section("extlibs").MapTo(&GlobalConfig.Extlibs); err != nil {
		log.Fatalf("Fail to map config file: %v", err)
		os.Exit(1)
	}
	log.Println("Rulex config init successfully")
	return GlobalConfig
}

/*
*
* 设置go的线程，通常=0 不需要配置
*
 */
func SetGomaxProcs(GomaxProcs int) {
	if GomaxProcs > 0 {
		if GlobalConfig.GomaxProcs < runtime.NumCPU() {
			runtime.GOMAXPROCS(GlobalConfig.GomaxProcs)
		}
	}
}

/*
*
* 设置性能，通常用来Debug用，生产环境建议关闭
*
 */
func SetDebugMode(EnablePProf bool) {

	//------------------------------------------------------
	// pprof: https://segmentfault.com/a/1190000016412013
	//------------------------------------------------------
	if EnablePProf {
		log.Println("Start PProf debug at: 0.0.0.0:6060")
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		runtime.SetCPUProfileRate(1)
		go http.ListenAndServe("0.0.0.0:6060", nil)
	}
}

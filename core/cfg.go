package core

import "github.com/ngaut/log"

import "gopkg.in/ini.v1"
import "os"

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
	GlobalConfig.Name = cfg.Section("app").Key("name").MustString("rulex")
	GlobalConfig.MaxQueueSize = cfg.Section("app").Key("max_queue_size").MustInt(5000)
	GlobalConfig.ResourceRestartInterval = cfg.Section("app").Key("resource_restart_interval").MustInt(204800)
	GlobalConfig.Path = cfg.Section("cloud").Key("path").MustString("")
	GlobalConfig.Token = cfg.Section("cloud").Key("token").MustString("")
	GlobalConfig.Secret = cfg.Section("cloud").Key("secret").MustString("")
	log.Info("Rulex config init successfully")

}

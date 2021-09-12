package core

import "github.com/ngaut/log"

import "gopkg.in/ini.v1"
import "os"

//
// Global config
//
type RulexConfig struct {
	Name   string
	Path   string
	Token  string
	Secret string
}

var GlobalConfig RulexConfig

//
// Init config
//
func init() {
	log.Info("Init rulex config")
	cfg, err := ini.Load("conf/rulex.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
		os.Exit(1)
	}
	GlobalConfig.Name = cfg.Section("app").Key("name").MustString("rulex")
	GlobalConfig.Path = cfg.Section("cloud").Key("path").MustString("")
	GlobalConfig.Token = cfg.Section("cloud").Key("token").MustString("")
	GlobalConfig.Secret = cfg.Section("cloud").Key("secret").MustString("")
	log.Info("Rulex config init success.")

}

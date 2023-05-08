package utils

import (
	"errors"
	"reflect"

	"github.com/hootrhino/rulex/glogger"
	"gopkg.in/ini.v1"
)

/*
*
* 把ini配置映射成结构体
*
* type s struct {
*     Name string`ini:"name"`
* }
 */

func INIToStruct(iniPath string, s string, v interface{}) error {
	cfg, err := ini.Load(iniPath)
	if err != nil {
		glogger.GLogger.Fatalf("Fail to read config file: %v", err)
	}
	return cfg.Section(s).MapTo(v)
}

/*
*
* GetINI
*
 */
func GetINISection(iniPath string, s string) *ini.Section {
	cfg, err := ini.Load(iniPath)
	if err != nil {
		glogger.GLogger.Fatalf("Fail to read config file: %v", err)
	}
	return cfg.Section(s)
}

// INI转结构体
func InIMapToStruct(section *ini.Section, s interface{}) error {
	if reflect.ValueOf(s).Kind() != reflect.Ptr {
		return errors.New("config must be a pointer")
	}
	return section.MapTo(s)
}

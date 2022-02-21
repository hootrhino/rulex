package utils

import (
	"log"
	"os"

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
		log.Fatalf("Fail to read config file: %v", err)
		os.Exit(1)
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
		log.Fatalf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	return cfg.Section(s)
}

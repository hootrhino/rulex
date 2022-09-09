package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/engine"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

func HttpPost(data map[string]interface{}, url string) string {
	p, errs1 := json.Marshal(data)
	if errs1 != nil {
		glogger.GLogger.Fatal(errs1)
	}
	r, errs2 := http.Post(url, "application/json", bytes.NewBuffer(p))
	if errs2 != nil {
		glogger.GLogger.Fatal(errs2)
	}
	defer r.Body.Close()

	body, errs5 := io.ReadAll(r.Body)
	if errs5 != nil {
		glogger.GLogger.Fatal(errs5)
	}
	return string(body)
}

func HttpGet(api string) string {
	var err error
	request, err := http.NewRequest("GET", api, nil)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	return string(body)
}

/*
*
* 起一个测试服务
*
 */
func RunTestEngine() typex.RuleX {
	mainConfig := core.InitGlobalConfig("conf/rulex.ini")
	glogger.StartNewRealTimeLogger(core.GlobalConfig.LogLevel)
	glogger.StartGLogger(mainConfig.EnableConsole, core.GlobalConfig.LogPath)
	glogger.StartLuaLogger(core.GlobalConfig.LuaLogPath)
	//
	core.StartStore(core.GlobalConfig.MaxQueueSize)
	core.SetLogLevel()
	core.SetPerformance()
	// engine
	engine := engine.NewRuleEngine(mainConfig)
	return engine
}

/*
*
* 生成测试数据库的文件名
*
 */
func GenDate() string {
	return "rulex-test_" + time.Now().Format("2006-01-02-15_04_05")
}

/*
*
* 创建文件夹
*
 */
func MKDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

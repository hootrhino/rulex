// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"gopkg.in/ini.v1"
)

var GlobalConfig typex.RulexConfig
var INIPath string

// Init config
func InitGlobalConfig(path string) typex.RulexConfig {
	log.Println("Init rulex config:", path)
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
	if EnablePProf {
		go func() {
			readyDebug := false
			for {
				select {
				case <-context.Background().Done():
					{
						glogger.GLogger.Info("PProf exited")
						return
					}
				default:
					{
						time.Sleep(utils.GiveMeSeconds(3))
						if !readyDebug {
							fmt.Printf("HeapObjects,\tHeapAlloc,\tTotalAlloc,\tHeapSys")
							fmt.Printf(",\tHeapIdle,\tHeapReleased,\tHeapIdle-HeapReleased")
							fmt.Println()
						}
						readyDebug = true
						utils.TraceMemStats()
					}
				}
			}

		}()

	}
}

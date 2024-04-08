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

package aibase

import (
	"fmt"
	"sync"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultAIRuntime *AlgorithmRuntime

/*
*
* 管理器
*
 */

type AlgorithmRuntime struct {
	RuleEngine  typex.RuleX
	locker      sync.Mutex
	XAlgorithms map[string]XAlgorithm
}

/*
*
* 初始化
*
 */
func InitAlgorithmRuntime(re typex.RuleX) *AlgorithmRuntime {
	__DefaultAIRuntime = new(AlgorithmRuntime)
	__DefaultAIRuntime.RuleEngine = re
	__DefaultAIRuntime.XAlgorithms = make(map[string]XAlgorithm)
	__DefaultAIRuntime.locker = sync.Mutex{}
	// Yolo8
	// err1 := LoadAlgorithm(NewYolo8ObjectDetectionCpu(), map[string]interface{}{})
	// if err1 != nil {
	// 	glogger.GLogger.Error(err1)
	// }
	// Tensorflow
	// err2 :=LoadAlgorithm(NewTfObjectDetectionCpu(), map[string]interface{}{})
	// if err2 != nil {
	// 	glogger.GLogger.Error(err1)
	// }
	return __DefaultAIRuntime
}

/*
*
* 停止运行时
*
 */
func Stop() {
	for _, v := range __DefaultAIRuntime.XAlgorithms {
		glogger.GLogger.Info("Try to Stop Algorithm:", v.AlgorithmDetail().Name)
		v.Unload()
		glogger.GLogger.Info("Algorithm Stop:", v.AlgorithmDetail().Name, " Success")
	}
	glogger.GLogger.Info("Algorithm Runtime stopped")
}

/*
*
* 列表
*
 */
func ListAlgorithm() []XAlgorithm {
	ll := []XAlgorithm{}
	for _, v := range __DefaultAIRuntime.XAlgorithms {
		ll = append(ll, v)
	}
	return ll
}

/*
*
* 加载算法
*
 */
func LoadAlgorithm(Algorithm XAlgorithm, Config map[string]interface{}) error {
	__DefaultAIRuntime.locker.Lock()
	defer __DefaultAIRuntime.locker.Unlock()
	if err := Algorithm.Init(Config); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if err := Algorithm.Load(); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	Info := Algorithm.AlgorithmDetail()
	__DefaultAIRuntime.XAlgorithms[Info.UUID] = Algorithm
	return fmt.Errorf("not support:%s", Algorithm.AlgorithmDetail().UUID)
}

/*
*
* 获取一个算法
*
 */
func GetAlgorithm(uuid string) XAlgorithm {
	return __DefaultAIRuntime.XAlgorithms[uuid]
}

/*
*
* 更新算法
*
 */
func UpdateAlgorithm(Algorithm XAlgorithm) error {
	__DefaultAIRuntime.locker.Lock()
	defer __DefaultAIRuntime.locker.Unlock()
	if _, ok := __DefaultAIRuntime.XAlgorithms[Algorithm.AlgorithmDetail().UUID]; ok {
		__DefaultAIRuntime.XAlgorithms[Algorithm.AlgorithmDetail().UUID] = Algorithm
		glogger.GLogger.Error("XAI.Start updated")
		return nil
	}
	return fmt.Errorf("Algorithm not exists:" + Algorithm.AlgorithmDetail().UUID)
}

/*
*
* 卸载算法
*
 */
func UnloadAlgorithm(uuid string) error {
	__DefaultAIRuntime.locker.Lock()
	defer __DefaultAIRuntime.locker.Unlock()
	if Algorithm, ok := __DefaultAIRuntime.XAlgorithms[uuid]; ok {
		Algorithm.Unload()
		glogger.GLogger.Error("XAI.Start stopped")
		return nil
	}
	return nil
}

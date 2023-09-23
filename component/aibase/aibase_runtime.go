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

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultAIRuntime *AIRuntime

func AIBaseRuntime() *AIRuntime {
	return __DefaultAIRuntime
}
func InitAIRuntime(re typex.RuleX) *AIRuntime {
	__DefaultAIRuntime = new(AIRuntime)
	__DefaultAIRuntime.RuleEngine = re
	__DefaultAIRuntime.AiBases = make(map[string]*typex.AI)
	// 预加载内置模型
	LoadAi(&typex.AI{
		UUID:      "BODY_POSE_RECOGNITION",
		Name:      "人体姿态识别",
		Type:      typex.BODY_POSE_RECOGNITION,
		IsBuildIn: true,
		Filepath:  "...",
		Config: map[string]interface{}{
			"algorithm":   "ANN",
			"inputs":      10,
			"layout":      []int{5, 3, 3},
			"output":      3,
			"activation:": "Sigmoid",
			"bias":        true,
		},
		Description: "一个轻量级人体姿态识别模型",
	})
	return __DefaultAIRuntime
}

/*
*
* 管理器
*
 */

type AIRuntime struct {
	RuleEngine typex.RuleX
	AiBases    map[string]*typex.AI
}

func GetRuleX() typex.RuleX {
	return __DefaultAIRuntime.RuleEngine
}
func ListAi() []*typex.AI {
	ll := []*typex.AI{}
	for _, v := range __DefaultAIRuntime.AiBases {
		ll = append(ll, v)
	}
	return ll
}
func LoadAi(Ai *typex.AI) error {
	if Ai.Type == typex.BODY_POSE_RECOGNITION {
		Ai.XAI = NewBodyPoseRecognition(__DefaultAIRuntime.RuleEngine)
		__DefaultAIRuntime.AiBases[Ai.UUID] = Ai
	}
	return fmt.Errorf("not support:%s", Ai.Type)
}
func GetAi(uuid string) *typex.AI {
	return __DefaultAIRuntime.AiBases[uuid]
}
func RemoveAi(uuid string) error {
	if v, ok := __DefaultAIRuntime.AiBases[uuid]; ok {
		// 内建类型不可删除
		if v.IsBuildIn {
			return fmt.Errorf("can not remove build-in aibase")
		}
		delete(__DefaultAIRuntime.AiBases, uuid)
		glogger.GLogger.Error("XAI.Start deleted")
		return nil
	}
	return fmt.Errorf("aibase not exists:" + uuid)

}
func UpdateAi(Ai *typex.AI) error {
	if v, ok := __DefaultAIRuntime.AiBases[Ai.UUID]; ok {
		// 内建类型不可修改
		if v.IsBuildIn {
			return fmt.Errorf("can not change 'BUILDIN' aibase")
		}
		__DefaultAIRuntime.AiBases[Ai.UUID] = Ai
		glogger.GLogger.Error("XAI.Start updated")

		return nil
	}
	return fmt.Errorf("aibase not exists:" + Ai.UUID)
}
func StartAi(uuid string) error {
	if ai, ok := __DefaultAIRuntime.AiBases[uuid]; ok {
		// 内建类型不可修改
		if ai.IsBuildIn {
			return fmt.Errorf("can not change 'BUILDIN' aibase")
		}
		err := ai.XAI.Start(map[string]interface{}{})
		if err != nil {
			glogger.GLogger.Error("XAI.Start error:", err)
		}
		return err
	}
	return nil
}
func StopAi(uuid string) error {
	if ai, ok := __DefaultAIRuntime.AiBases[uuid]; ok {
		// 内建类型不可修改
		if ai.IsBuildIn {
			return fmt.Errorf("can not change 'BUILDIN' aibase")
		}
		ai.XAI.Stop()
		glogger.GLogger.Error("XAI.Start stopped")
		return nil
	}
	return nil
}
func Stop() {
	glogger.GLogger.Info("AIRuntime stopped")
}

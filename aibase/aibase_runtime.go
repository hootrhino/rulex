package aibase

import (
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

func NewAIRuntime(re typex.RuleX) typex.XAiRuntime {
	ai := new(AIRuntime)
	ai.re = re
	ai.aiBases = make(map[string]*typex.AI)
	// 预加载内置模型
	ai.LoadAi(&typex.AI{
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
	return ai
}

/*
*
* 管理器
*
 */

type AIRuntime struct {
	re      typex.RuleX
	aiBases map[string]*typex.AI
}

func (airt *AIRuntime) GetRuleX() typex.RuleX {
	return airt.re
}
func (airt *AIRuntime) ListAi() []*typex.AI {
	ll := []*typex.AI{}
	for _, v := range airt.aiBases {
		ll = append(ll, v)
	}
	return ll
}
func (airt *AIRuntime) LoadAi(Ai *typex.AI) error {
	if Ai.Type == typex.BODY_POSE_RECOGNITION {
		Ai.XAI = NewBodyPoseRecognition(airt.re)
		airt.aiBases[Ai.UUID] = Ai
	}
	return fmt.Errorf("not support:%s", Ai.Type)
}
func (airt *AIRuntime) GetAi(uuid string) *typex.AI {
	return airt.aiBases[uuid]
}
func (airt *AIRuntime) RemoveAi(uuid string) error {
	if v, ok := airt.aiBases[uuid]; ok {
		// 内建类型不可删除
		if v.IsBuildIn {
			return fmt.Errorf("can not remove build-in aibase")
		}
		delete(airt.aiBases, uuid)
		glogger.GLogger.Error("XAI.Start deleted")
		return nil
	}
	return fmt.Errorf("aibase not exists:" + uuid)

}
func (airt *AIRuntime) UpdateAi(Ai *typex.AI) error {
	if v, ok := airt.aiBases[Ai.UUID]; ok {
		// 内建类型不可修改
		if v.IsBuildIn {
			return fmt.Errorf("can not change 'BUILDIN' aibase")
		}
		airt.aiBases[Ai.UUID] = Ai
		glogger.GLogger.Error("XAI.Start updated")

		return nil
	}
	return fmt.Errorf("aibase not exists:" + Ai.UUID)
}
func (airt *AIRuntime) StartAi(uuid string) error {
	if ai, ok := airt.aiBases[uuid]; ok {
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
func (airt *AIRuntime) StopAi(uuid string) error {
	if ai, ok := airt.aiBases[uuid]; ok {
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
func (airt *AIRuntime) Stop() {
	glogger.GLogger.Info("AIRuntime stopped")
}

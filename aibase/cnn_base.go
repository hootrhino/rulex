package aibase

import "github.com/i4de/rulex/typex"

func NewCNNBaseRuntime(re typex.RuleX) typex.XAiBase {
	ai := new(CNNBaseRuntime)
	ai.re = re
	ai.aiBases = make(map[string]*typex.AiBase)
	return ai
}

/*
*
* 管理器
*
 */
type CNNBaseRuntime struct {
	re      typex.RuleX
	aiBases map[string]*typex.AiBase
}

func (cnn *CNNBaseRuntime) GetRuleX() typex.RuleX {

	return cnn.re
}
func (cnn *CNNBaseRuntime) ListAi() []*typex.AiBase {
	result := []*typex.AiBase{}
	for _, v := range cnn.aiBases {
		result = append(result, v)
	}
	return result
}
func (cnn *CNNBaseRuntime) LoadAi(Ai *typex.AiBase) error {
	return nil
}
func (cnn *CNNBaseRuntime) GetAi(uuid string) *typex.AiBase {
	return cnn.aiBases[uuid]
}
func (cnn *CNNBaseRuntime) RemoveAi(uuid string) error {
	return nil
}
func (cnn *CNNBaseRuntime) UpdateAi(Ai typex.AiBase) error {
	return nil
}
func (cnn *CNNBaseRuntime) StartAi(uuid string) error {
	return nil
}
func (cnn *CNNBaseRuntime) StopAi(uuid string) error {
	return nil
}
func (cnn *CNNBaseRuntime) Stop() {

}

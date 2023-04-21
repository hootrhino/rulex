package aibase

import "github.com/i4de/rulex/typex"

func NewANNBaseRuntime(re typex.RuleX) typex.XAiBase {
	ai := new(ANNBaseRuntime)
	ai.re = re
	ai.aiBases = make(map[string]*typex.AiBase)
	return ai
}

/*
*
* 管理器
*
 */
type ANNBaseRuntime struct {
	re      typex.RuleX
	aiBases map[string]*typex.AiBase
}

func (ann *ANNBaseRuntime) GetRuleX() typex.RuleX {

	return ann.re
}
func (ann *ANNBaseRuntime) ListAi() []*typex.AiBase {
	result := []*typex.AiBase{}
	for _, v := range ann.aiBases {
		result = append(result, v)
	}
	return result

}
func (ann *ANNBaseRuntime) LoadAi(Ai *typex.AiBase) error {
	return nil
}
func (ann *ANNBaseRuntime) GetAi(uuid string) *typex.AiBase {
	return ann.aiBases[uuid]
}
func (ann *ANNBaseRuntime) RemoveAi(uuid string) error {
	return nil
}
func (ann *ANNBaseRuntime) UpdateAi(Ai typex.AiBase) error {
	return nil
}
func (ann *ANNBaseRuntime) StartAi(uuid string) error {
	return nil
}
func (ann *ANNBaseRuntime) StopAi(uuid string) error {
	return nil
}
func (ann *ANNBaseRuntime) Stop() {

}

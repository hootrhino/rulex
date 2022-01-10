package core

import "rulex/typex"

type targetTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[string]func(typex.RuleX) typex.XTarget
}

func NewTargetTypeManager() typex.TargetRegistry {
	return &targetTypeManager{
		registry: map[string]func(typex.RuleX) typex.XTarget{},
	}

}
func (rm *targetTypeManager) Register(name string, f func(typex.RuleX) typex.XTarget) {
	rm.registry[name] = f
}

func (rm *targetTypeManager) Find(name string) func(typex.RuleX) typex.XTarget {

	return rm.registry[name]
}

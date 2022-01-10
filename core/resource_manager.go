package core

import "rulex/typex"

type resourceTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[string]func(typex.RuleX) typex.XResource
}

func NewResourceTypeManager() typex.ResourceRegistry {
	return &resourceTypeManager{
		registry: map[string]func(typex.RuleX) typex.XResource{},
	}

}
func (rm *resourceTypeManager) Register(name string, f func(typex.RuleX) typex.XResource) {
	rm.registry[name] = f
}

func (rm *resourceTypeManager) Find(name string) func(typex.RuleX) typex.XResource {

	return rm.registry[name]
}

package core

import "rulex/typex"

type resourceTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.InEndType]*typex.XConfig
}

func NewResourceTypeManager() typex.ResourceRegistry {
	return &resourceTypeManager{
		registry: map[typex.InEndType]*typex.XConfig{},
	}

}
func (rm *resourceTypeManager) Register(name typex.InEndType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *resourceTypeManager) Find(name typex.InEndType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *resourceTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

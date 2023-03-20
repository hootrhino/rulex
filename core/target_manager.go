package core

import "github.com/i4de/rulex/typex"

type TargetTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.TargetType]*typex.XConfig
}

func NewTargetTypeManager() typex.TargetRegistry {
	return &TargetTypeManager{
		registry: map[typex.TargetType]*typex.XConfig{},
	}

}
func (rm *TargetTypeManager) Register(name typex.TargetType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *TargetTypeManager) Find(name typex.TargetType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *TargetTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

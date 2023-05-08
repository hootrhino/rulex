package core

import "github.com/hootrhino/rulex/typex"

type SourceTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.InEndType]*typex.XConfig
}

func NewSourceTypeManager() typex.SourceRegistry {
	return &SourceTypeManager{
		registry: map[typex.InEndType]*typex.XConfig{},
	}

}
func (rm *SourceTypeManager) Register(name typex.InEndType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *SourceTypeManager) Find(name typex.InEndType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *SourceTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

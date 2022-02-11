package core

import "rulex/typex"

type sourceTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.InEndType]*typex.XConfig
}

func NewSourceTypeManager() typex.SourceRegistry {
	return &sourceTypeManager{
		registry: map[typex.InEndType]*typex.XConfig{},
	}

}
func (rm *sourceTypeManager) Register(name typex.InEndType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *sourceTypeManager) Find(name typex.InEndType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *sourceTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

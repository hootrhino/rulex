package core

import "github.com/i4de/rulex/typex"

type deviceTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.DeviceType]*typex.XConfig
}

func NewDeviceTypeManager() typex.DeviceRegistry {
	return &deviceTypeManager{
		registry: map[typex.DeviceType]*typex.XConfig{},
	}

}
func (rm *deviceTypeManager) Register(name typex.DeviceType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *deviceTypeManager) Find(name typex.DeviceType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *deviceTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

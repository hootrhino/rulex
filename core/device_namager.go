package core

import "github.com/i4de/rulex/typex"

type DeviceTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.DeviceType]*typex.XConfig
}

func NewDeviceTypeManager() typex.DeviceRegistry {
	return &DeviceTypeManager{
		registry: map[typex.DeviceType]*typex.XConfig{},
	}

}
func (rm *DeviceTypeManager) Register(name typex.DeviceType, f *typex.XConfig) {
	rm.registry[name] = f
	f.Type = string(name)
}

func (rm *DeviceTypeManager) Find(name typex.DeviceType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *DeviceTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

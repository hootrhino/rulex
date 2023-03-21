package typex

type DeviceRegistry interface {
	Register(DeviceType, *XConfig)
	Find(DeviceType) *XConfig
	All() []*XConfig
}

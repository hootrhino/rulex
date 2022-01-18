package typex

type ResourceRegistry interface {
	Register(InEndType, *XConfig)
	Find(InEndType) *XConfig
	All() []*XConfig
}

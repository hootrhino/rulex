package typex

type SourceRegistry interface {
	Register(InEndType, *XConfig)
	Find(InEndType) *XConfig
	All() []*XConfig
}

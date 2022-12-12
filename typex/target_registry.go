package typex

type TargetRegistry interface {
	Register(TargetType, *XConfig)
	Find(TargetType) *XConfig
	All() []*XConfig
}

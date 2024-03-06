package typex

import "context"

// Global context
var GCTX, GCancel = context.WithCancel(context.Background())

// child context
type CCTX struct {
	Ctx       context.Context
	CancelCTX context.CancelFunc
}

func NewCCTX() (context.Context, context.CancelFunc) {
	ctx, cancelCTX := context.WithCancel(GCTX)
	return ctx, cancelCTX
}

/*
*
* 资源管理器
*
 */
type ResourceRegistry interface {
	Register(DeviceType, *XConfig)
	Find(DeviceType) *XConfig
	All() []*XConfig
}
type RulexResourceRegistry struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[InEndType]*XConfig
}

func NewRulexResourceRegistry() *RulexResourceRegistry {
	return &RulexResourceRegistry{
		registry: map[InEndType]*XConfig{},
	}
}
func (rm *RulexResourceRegistry) Register(name InEndType, f *XConfig) {
	rm.registry[name] = f
}

func (rm *RulexResourceRegistry) Find(name InEndType) *XConfig {

	return rm.registry[name]
}
func (rm *RulexResourceRegistry) All() []*XConfig {
	data := make([]*XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

package typex

import "context"

// Global context
var GCTX = context.Background()

// child context
type CCTX struct {
	Ctx       context.Context
	CancelCTX context.CancelFunc
}

func NewCCTX() (context.Context, context.CancelFunc) {
	ctx, cancelCTX := context.WithCancel(GCTX)
	return ctx, cancelCTX
}

package typex

import "context"

// Global context
var GCTX = context.Background()

// child context
type CCTX struct {
	Ctx       context.Context
	CancelCTX context.CancelFunc
}

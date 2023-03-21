package typex

//------------------------------------------------
// 					Remote Stream
//------------------------------------------------
// ┌───────────────┐          ┌────────────────┐
// │   RULEX       │ <─────── │   SERVER       │
// │   RULEX       │  ───────>│   SERVER       │
// └───────────────┘          └────────────────┘
//------------------------------------------------
type XStream interface {
	Start() error
	State() XStatus
	Close()
}

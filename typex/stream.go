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
	OnStreamApproached(data string) error
	State() XStatus
	Close()
}

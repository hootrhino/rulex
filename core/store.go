package core

import (
	"github.com/i4de/rulex/store"
	"github.com/i4de/rulex/typex"
)

var GlobalStore typex.XStore

func StartStore(maxSize int) {
	GlobalStore = store.NewRulexStore(maxSize)

}

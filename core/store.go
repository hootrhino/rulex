package core

import (
	"rulex/store"
	"rulex/typex"
)

var GlobalStore typex.XStore

func StartStore(maxSize int) {
	GlobalStore = store.NewRulexStore(maxSize)

}

package test

import (
	"testing"

	"github.com/i4de/rulex/core"
)

func Test_get_set(t *testing.T) {
	core.StartStore(1024)
	core.GlobalStore.Set("k", "v")
	t.Log(core.GlobalStore.Get("k"))
}

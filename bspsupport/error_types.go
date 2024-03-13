package archsupport

import (
	"errors"
	"runtime"
)

var errArchNotSupport = errors.New("not support current OS:" + runtime.GOOS)

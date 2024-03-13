package archsupport

import (
	"errors"
	"runtime"
)

var errInvalidLen = errors.New("invalid length")
var errInvalidValue = errors.New("invalid value")
var errArchNotSupport = errors.New("not support current OS:" + runtime.GOOS)

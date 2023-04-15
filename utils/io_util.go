package utils

import (
	"context"
	"errors"
	"io"
)

// 直接把io包下面的同名函数抄过来了，加上了Context支持，主要解决读取超时问题
var errShortBuffer = errors.New("short buffer")
var errEOF = errors.New("EOF")
var errUnexpectedEOF = errors.New("unexpected EOF")
var errTimeout = errors.New("read Timeout")

/*
* 读取字节，核心原理是一个一个读，这样就不会出问题.
*
 */

func ReadAtLeast(ctx context.Context, r io.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		n = 0
		err = errShortBuffer
		return
	}
	for n < min && err == nil {
		select {
		case <-ctx.Done():
			err = errTimeout
			return
		default:
		}
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == errEOF {
		err = errUnexpectedEOF
	}
	return
}

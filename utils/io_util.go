package utils

import (
	"context"
	"errors"
	"io"
	"time"
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

/*
*
* 时间片读写请求
*
 */
func SliceRequest(ctx context.Context,
	iio io.ReadWriter, writeBytes []byte,
	resultBuffer []byte,
	showError bool,
	td time.Duration) (int, error) {
	_, errW := iio.Write(writeBytes)
	if errW != nil {
		return 0, errW
	}
	return SliceReceive(ctx, iio, resultBuffer, showError, td)
}

/*
*
* 读取数据的时候，如果出现错误就返回
*
 */
func SliceReceiveWithError(ctx context.Context,
	iio io.Reader, resultBuffer []byte, td time.Duration) (int, error) {
	return SliceReceive(ctx, iio, resultBuffer, true, td)
}

/*
*
* 读取数据的时候，如果出现错误则判断是否是串口引起的超时，如果是就忽略
*
 */
func SliceReceiveWithoutError(ctx context.Context,
	iio io.Reader, resultBuffer []byte, td time.Duration) (int, error) {
	return SliceReceive(ctx, iio, resultBuffer, false, td)
}

/*
*
* 通过一个定时时间片读取
*
 */
func SliceReceive(ctx context.Context,
	iio io.Reader, resultBuffer []byte,
	showError bool,
	td time.Duration) (int, error) {
	var peerCount int
	sliceTimer := time.NewTimer(td)
	sliceTimer.Stop()
	for {
		select {
		case <-ctx.Done():
			return peerCount, nil
		case <-sliceTimer.C:
			return peerCount, nil
		default:
			readCount, errR := iio.Read(resultBuffer[peerCount:])
			// Note:当串口设置超时以后会输出：
			//    This operation returned because the timeout period expired.
			// 因此这里需要过滤这个异常, 当出现这个异常信息的时候认为其是正常的
			//
			if errR != nil {
				if showError {
					return peerCount, errR
				}
			}
			if readCount != 0 {
				sliceTimer.Reset(td)
				peerCount += readCount
			}
		}
	}
}

/*
*
* 某个时间片期望最少收到字节数
*
 */
func SliceReceiveAtLeast(ctx context.Context,
	iio io.Reader, resultBuffer []byte, td time.Duration, min int) (int, error) {
	// 后期实现
	return 0, nil
}

/*
*
* Slice分页计算器
*
 */
func Paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}

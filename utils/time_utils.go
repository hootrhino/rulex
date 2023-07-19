package utils

import "time"

/*
*
* 将go的时间单位稍作封装
*
 */

/*
 *
 * 返回秒
 *
 */
func GiveMeSeconds(c int64) time.Duration {
	return time.Second * time.Duration(c)
}

/*
 *
 * 返回毫秒
 *
 */
func GiveMeMilliseconds(c int64) time.Duration {
	return time.Millisecond * time.Duration(c)
}

/*
 *
 * 返回微秒
 *
 */
func GiveMeMicroseconds(c int64) time.Duration {
	return time.Microsecond * time.Duration(c)
}

// NowTs 返回当前Unix时间戳
func NowTs() int64 {
	return time.Now().Unix()
}

// NowTime 返回当前时间串，eg: 2023-07-10 13:14:52
func NowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Sleep 当前线程阻塞(休眠)指定时间，单位：毫秒
func Sleep(ts int) {
	time.Sleep(time.Millisecond * time.Duration(ts))
}

// Zone 返回当前时区
func Zone() string {
	zone, _ := time.Now().Zone()
	return zone
}

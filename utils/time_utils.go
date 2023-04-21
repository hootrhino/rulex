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

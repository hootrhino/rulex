package driver

import "encoding/hex"

/*
*
* 处理空字符串
*
 */
 func covertEmptyHex(v []byte) string {
	if len(v) < 1 {
		return ""
	}
	return hex.EncodeToString(v)
}

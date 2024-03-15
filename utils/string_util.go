package utils

import (
	"regexp"

	"github.com/plgd-dev/kit/v2/strings"
)

/*
*
* 检查列表元素是不是不重复了
*
 */
func IsListDuplicated(list []string) bool {
	tmpMap := make(map[string]int)
	for _, value := range list {
		tmpMap[value] = 1
	}
	var keys []interface{}
	for k := range tmpMap {
		keys = append(keys, k)
	}
	// 对比列表的Key元素和列表长度是否等长
	return len(keys) != len(list)
}

/*
*
* 列表包含
*
 */
func SContains(s []string, e string) bool {
	return strings.SliceContains(s, e)
}

/*
*
* 合法名称
*
 */
func IsValidName(username string) bool {
	// 检查用户名长度是否在 6 到 32 之间
	if len(username) < 6 || len(username) > 32 {
		return false
	}
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return match
}

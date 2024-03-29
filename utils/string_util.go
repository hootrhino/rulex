package utils

import (
	"regexp"
	"strings"
	"sync"
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
func SContains(slice []string, e string) bool {
	for _, s := range slice {
		if s == e {
			return true
		}
	}
	return false
}

/*
*
* 合法名称
*
 */
var (
	s               = `^[a-zA-Z_][a-zA-Z0-9._\u4e00-\u9fa5]*$`
	usernamePattern = regexp.MustCompile(s)
	once            sync.Once
)

func IsValidName(username string) bool {
	once.Do(func() {
		usernamePattern = regexp.MustCompile(s)
	})
	return usernamePattern.MatchString(username)
}

/*
*
* Tag name 不可出现非法字符
*
 */
func IsValidColumnName(columnName string) bool {
	// 列名不能以数字开头
	if len(columnName) == 0 || (columnName[0] >= '0' && columnName[0] <= '9') {
		return false
	}
	invalidChars := []string{" ", "-", ";"}
	for _, char := range invalidChars {
		if strings.Contains(columnName, char) {
			return false
		}
	}
	return true
}

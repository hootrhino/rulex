package utils

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

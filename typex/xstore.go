package typex

/*
*
* 缓存器接口
*
 */
type XStore interface {
	// 设置值
	Set(k string, v string)
	// 获取值
	Get(k string) string
	// 删除值
	Delete(k string) error
	// 统计数量
	Count() int
	// 模糊查询匹配
	// 支持: *AAA AAA* A*B
	FuzzyGet(k string) string
}

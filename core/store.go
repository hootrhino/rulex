package core

import (
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	cache "github.com/wwhai/tinycache"
)

var GlobalStore typex.XStore

func StartStore(maxSize int) {
	GlobalStore = NewRulexStore(maxSize)

}

type RulexStore struct {
	cache *cache.Cache
}

func NewRulexStore(maxSize int) typex.XStore {
	return &RulexStore{
		cache: cache.New(0, 0),
	}

}

/*
*
* 设置过期时间
*
 */
func (rs *RulexStore) SetDuration(k string, v string, d time.Duration) {
	if (rs.cache.ItemCount() + 1) > 10000 {
		glogger.GLogger.Error("Max store size reached:", rs.cache.ItemCount())
	}
	rs.cache.Set(k, v, d)
}

// 设置值
func (rs *RulexStore) Set(k string, v string) {
	if (rs.cache.ItemCount() + 1) > 10000 {
		glogger.GLogger.Error("Max store size reached:", rs.cache.ItemCount())
	}
	rs.cache.Set(k, v, -1)
}

// 获取值
func (rs *RulexStore) Get(k string) string {
	v, ok := rs.cache.Get(k)
	if ok {
		return v.(string)
	} else {
		return ""
	}
}
func (rs *RulexStore) Delete(k string) error {
	rs.cache.Delete(k)
	return nil
}

// 统计数量
func (rs *RulexStore) Count() int {
	return rs.cache.ItemCount()
}

// 模糊查询匹配
// 支持: *AAA AAA* A*B
func (rs *RulexStore) FuzzyGet(k string) string {
	// TODO 未来实现
	return ""
}

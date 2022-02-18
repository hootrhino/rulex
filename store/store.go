package store

import (
	"rulex/typex"
	"sync"

	"github.com/ngaut/log"
)

type RulexStore struct {
	maxSize int
	len     int
	bucket  sync.Map
}

func NewRulexStore(maxSize int) typex.XStore {
	return &RulexStore{
		maxSize: maxSize,
		len:     0,
		bucket:  sync.Map{},
	}

}

// 设置值
func (rs *RulexStore) Set(k string, v string) {
	rs.bucket.Store(k, v)
	if (rs.len + 1) > rs.len {
		log.Error("Max store size reached:", rs.len)
	} else {
		rs.len += 1
	}
}

// 获取值
func (rs *RulexStore) Get(k string) string {
	v, ok := rs.bucket.Load(k)
	if ok {
		return v.(string)
	} else {
		return ""
	}
}
func (rs *RulexStore) Delete(k string) error {
	rs.bucket.Delete(k)
	if rs.len > 0 {
		rs.len -= 1
	}
	return nil
}

// 统计数量
func (rs *RulexStore) Count() int {
	return rs.len
}

// 模糊查询匹配
// 支持: *AAA AAA* A*B
func (rs *RulexStore) FuzzyGet(k string) string {
	// TODO 未来实现
	return ""
}

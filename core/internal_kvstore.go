// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"strings"
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
func (rs *RulexStore) FuzzyGet(Key string) any {
	for k, v := range rs.cache.Items() {
		if fuzzyMatch(k, Key) {
			return v
		}
	}
	return ""
}

// 将主字符串和子字符串都转换为小写，以进行不区分大小写的匹配
func fuzzyMatch(mainStr, subStr string) bool {
	mainStr = strings.ToLower(mainStr)
	subStr = strings.ToLower(subStr)
	return strings.Contains(mainStr, subStr)
}

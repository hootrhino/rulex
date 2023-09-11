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
	"sync"

	"github.com/hootrhino/rulex/typex"
)

/*
*
* 全局缓冲器，用来保存内部的数据模型用
*
 */
var __InternalSchemaCache map[string]typex.DataSchema
var __lock sync.Mutex

/*
*
* 初始化缓冲器
*
 */
func InitInternalSchemaCache() {
	__InternalSchemaCache = make(map[string]typex.DataSchema)
	__lock = sync.Mutex{}
}

/*
*
* 第一次缓冲
*
 */
func FirstCache(id string, schema typex.DataSchema) typex.DataSchema {
	if DataSchema, ok := __InternalSchemaCache[id]; ok {
		return DataSchema
	}
	SchemaSet(id, schema)
	return schema
}

/*
*
* 增加一条缓存数据
*
 */
func SchemaSet(id string, schema typex.DataSchema) {
	__lock.Lock()
	__InternalSchemaCache[id] = schema
	__lock.Unlock()
}

/*
*
* 获取某个模型
*
 */
func SchemaGet(id string) (typex.DataSchema, bool) {
	v, ok := __InternalSchemaCache[id]
	return v, ok
}

/*
*
* 删除
*
 */
func SchemaDelete(id string) {
	delete(__InternalSchemaCache, id)
}

/*
*
* 数目
*
 */
func SchemaCount() int {
	return len(__InternalSchemaCache)
}

/*
*
* 清空数据
*
 */
func SchemaFlush() {
	for key := range __InternalSchemaCache {
		delete(__InternalSchemaCache, key)
	}
}

/*
*
* 所有
*
 */
func AllSchema() map[string]typex.DataSchema {
	return __InternalSchemaCache
}

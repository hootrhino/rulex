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
package iotschema

import (
	"sync"

	"github.com/hootrhino/rulex/component/iotschema"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 全局缓冲器，用来保存内部的数据模型用
*
 */
var __InternalSchemaCache map[string]iotschema.IoTSchema
var __lock sync.Mutex

/*
*
* 初始化缓冲器
*
 */
func InitInternalSchemaCache(rulex typex.RuleX) {
	__InternalSchemaCache = make(map[string]iotschema.IoTSchema)
	__lock = sync.Mutex{}
}

/*
*
* 第一次缓冲
*
 */
func FirstCache(id string, schema iotschema.IoTSchema) iotschema.IoTSchema {
	__lock.Lock()
	defer __lock.Unlock()
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
func SchemaSet(id string, schema iotschema.IoTSchema) {
	__lock.Lock()
	defer __lock.Unlock()
	__InternalSchemaCache[id] = schema
}

/*
*
* 获取某个模型
*
 */
func SchemaGet(id string) (iotschema.IoTSchema, bool) {
	v, ok := __InternalSchemaCache[id]
	return v, ok
}

/*
*
* 删除
*
 */
func SchemaDelete(id string) {
	__lock.Lock()
	defer __lock.Unlock()
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
	__lock.Lock()
	defer __lock.Unlock()
	for key := range __InternalSchemaCache {
		delete(__InternalSchemaCache, key)
	}
}

/*
*
* 所有
*
 */
func AllSchema() map[string]iotschema.IoTSchema {
	return __InternalSchemaCache
}

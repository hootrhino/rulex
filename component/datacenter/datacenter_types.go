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

package datacenter

/*
*
* 数据仓库的详情
*
 */
type SchemaDetail struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	// SchemaType 数据库类型，如果是"INTERNAL"则表示原始RULEX；如果是 "EXTERNAL" 则表示扩展设备
	// 该字段是为了扩展数据中心, 将数据中心扩展为通用的一个中间层
	SchemaType   string       `json:"schemaType"`
	LocalPath    string       `json:"local_path"`
	NetAddr      string       `json:"net_addr"`
	CreateTs     uint64       `json:"create_ts"`
	Size         float32      `json:"size"`
	StorePath    string       `json:"store_path"`
	SchemaDefine SchemaDefine `json:"schema_define"`
	Description  string       `json:"description"`
}

/*
*
* 列定义
*
 */
type Column struct {
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Description string      `json:"description,omitempty"`
}

/*
*
*  表结构定义
*
 */
type SchemaDefine struct {
	UUID    string   `json:"uuid"`
	Columns []Column `json:"columns"`
}

/*
*
* 数据源接口
*
 */
type DataSource interface {
	Init() error
	GetSchemaDetail(goodsId string) SchemaDetail // 表详情
	Query(goodsId, query string) ([]map[string]any, error)
}

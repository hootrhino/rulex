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

package test

import (
	"encoding/json"
	"testing"

	"github.com/hootrhino/rulex/component/iotschema"
)

/*
*
* 新建模型
*
 */
type IoTSchemaVo struct {
	UUID   string              `json:"uuid,omitempty"`
	Name   string              `json:"name"`
	Schema iotschema.IoTSchema `json:"schema"`
}

// go test -timeout 30s -run ^Test_IoTSchema_gen github.com/hootrhino/rulex/test -v -count=1
func Test_IoTSchema_gen(t *testing.T) {
	IoTSchemaVo := IoTSchemaVo{
		UUID: "AABBCCDDEEFFTT",
		Name: "测试物模型: 机场大厅数据统计",
		Schema: iotschema.IoTSchema{
			IoTProperties: []iotschema.IoTProperty{
				{
					Label:       "整形：今天进出人数统计",
					Description: "今天进出人数统计",
					Name:        "peopleCount",
					Type:        "INTEGER",
					Rule: iotschema.IoTPropertyRule{
						Min:          0,
						Max:          100,
						DefaultValue: 0,
					},
					Unit: "人",
					Rw:   "R",
				},
				{
					Label:       "浮点型：温度",
					Description: "浮点型：温度",
					Name:        "temp",
					Type:        "FLOAT",
					Rule: iotschema.IoTPropertyRule{
						DefaultValue: 0,
						Min:          -100,
						Max:          100,
						Round:        2,
					},
					Unit: "摄氏度",
					Rw:   "R",
				},
				{
					Label:       "浮点型：湿度",
					Description: "浮点型：湿度",
					Name:        "humi",
					Type:        "FLOAT",
					Rule: iotschema.IoTPropertyRule{
						DefaultValue: 0,
						Min:          0,
						Max:          100,
						Round:        2,
					},
					Rw:   "R",
					Unit: "摄氏度",
				},
				{
					Label:       "布尔型：灯开关-1",
					Description: "布尔型：灯开关-1",
					Name:        "switcher1",
					Type:        "BOOL",
					Rule: iotschema.IoTPropertyRule{
						DefaultValue: false,
						TrueLabel:    "开启",
						FalseLabel:   "关闭",
					},
					Rw: "RW",
				},
				{
					Label:       "文本型：滚动显示屏幕文本",
					Description: "文本型：滚动显示屏幕文本",
					Name:        "string1",
					Type:        "STRING",
					Rule: iotschema.IoTPropertyRule{
						DefaultValue: "Hello, 大牛",
					},
					Rw: "R",
				},
			},
		},
	}
	b, _ := json.MarshalIndent(IoTSchemaVo, "", "    ")
	t.Log(string(b))
}

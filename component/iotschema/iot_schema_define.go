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
	"fmt"

	"github.com/hootrhino/rulex/utils"
)

type IoTPropertyType string

const (
	// 目前边缘侧暂时只支持常见类型
	IoTPropertyTypeString  IoTPropertyType = "STRING"
	IoTPropertyTypeInteger IoTPropertyType = "INTEGER"
	IoTPropertyTypeFloat   IoTPropertyType = "FLOAT"
	IoTPropertyTypeBool    IoTPropertyType = "BOOL"
	IoTPropertyTypeGeo     IoTPropertyType = "GEO"
)

// string
type IoTPropertyString string

// int
type IoTPropertyInteger int

// float
type IoTPropertyFloat float32

// bool
type IoTPropertyBool bool

// 地理坐标系统
type IoTPropertyGeo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

/*
* 物模型,边缘端目前暂时只支持属性
*
 */
type IoTSchema struct {
	IoTProperties map[string]IoTProperty `json:"iotProperties"`
}

// 规则
type IoTPropertyRule struct {
	DefaultValue any    `json:"defaultValue"`         // 默认值
	Max          int    `json:"max,omitempty"`        // 最大值
	Min          int    `json:"min,omitempty"`        // 最小值
	TrueLabel    string `json:"trueLabel,omitempty"`  // 真值label
	FalseLabel   string `json:"falseLabel,omitempty"` // 假值label
	Round        int    `json:"round,omitempty"`      // 小数点位
}

// 物模型属性
type IoTProperty struct {
	UUID        string          `json:"uuid"`            // Cache uuid
	Label       string          `json:"label"`           // UI显示的那个文本
	Name        string          `json:"name"`            // 变量关联名
	Description string          `json:"description"`     // 额外信息
	Type        IoTPropertyType `json:"type"`            // 类型, 只能是上面几种
	Rw          string          `json:"rw"`              // R读 W写 RW读写
	Unit        string          `json:"unit"`            // 单位 例如：摄氏度、米、牛等等
	Value       any             `json:"value,omitempty"` // Value 是运行时值, 前端不用填写
	Rule        IoTPropertyRule `json:"rule"`            // 规则
}

func (I *IoTProperty) StringValue() string {
	if I == nil {
		return ""
	}
	switch I.Type {
	case IoTPropertyTypeString:
		{
			return I.Value.(string)
		}
	}
	return ""
}
func (I *IoTProperty) IntValue() int {
	if I == nil {
		return 0
	}
	switch I.Type {
	case IoTPropertyTypeInteger:
		{
			return I.Value.(int)
		}
	}
	return 0
}
func (I *IoTProperty) FloatValue() float64 {
	if I == nil {
		return 0
	}
	switch I.Type {
	case IoTPropertyTypeFloat:
		{
			return I.Value.(float64)
		}
	}
	return 0
}
func (I *IoTProperty) BoolValue() bool {
	if I == nil {
		return false
	}
	switch I.Type {
	case IoTPropertyTypeBool:
		{
			return I.Value.(bool)
		}
	}
	return false
}

/*
*
* 验证类型
*
 */
func (V IoTProperty) ValidateFields() error {
	if utils.SContains([]string{"R", "W", "RW"}, V.Rw) {
		return fmt.Errorf("RW Value Only Support 'R' or 'W' or 'RW'")
	}
	return nil

}

/*
*
* 验证物模型本身是否合法, 包含了 IoTPropertyType，Rule 的对应关系
*
 */
func (V IoTProperty) ValidateRule() error {
	switch V.Type {
	case IoTPropertyTypeString:
		{
			return nil // TODO
		}
	case IoTPropertyTypeInteger:
		{
			return nil // TODO
		}
	case IoTPropertyTypeFloat:
		{
			return nil // TODO
		}
	case IoTPropertyTypeBool:
		{
			return nil // TODO
		}
	case IoTPropertyTypeGeo:
		{
			return nil // TODO
		}
	default:
		return fmt.Errorf("Unknown And Invalid IoT Property Type:%v", V.Type)
	}
}

/*
*
* 验证数据类型
*
 */
func (V IoTProperty) ValidateType() error {
	switch V.Type {
	case IoTPropertyTypeString:
		{
			return nil // TODO
		}
	case IoTPropertyTypeInteger:
		{
			return nil // TODO
		}
	case IoTPropertyTypeFloat:
		{
			return nil // TODO
		}
	case IoTPropertyTypeBool:
		{
			return nil // TODO
		}
	case IoTPropertyTypeGeo:
		{
			return nil // TODO
		}
	default:
		return fmt.Errorf("Unknown And Invalid IoT Property Type:%v", V.Type)
	}
}

/*
*
* 物模型规则 : String|Float|Int|Bool
*
 */
type SchemaRule interface {
	Validate(Value interface{}) error
}

/*
*
* 字符串规则
*
 */
type StringRule struct {
	MaxLength    int    `json:"maxLength"`
	DefaultValue string `json:"defaultValue"`
}

func (S StringRule) Validate(Value interface{}) error {
	switch SV := Value.(type) {
	case string:
		L := len(SV)
		if L >= S.MaxLength {
			return fmt.Errorf("Value exceed Max Length:", L)
		}
	default:
		{
			return fmt.Errorf("Invalid Value type, Expect UTF8 string:", SV)
		}
	}
	return nil
}

/*
*
* 整数规则
*
 */
type IntegerRule struct {
	DefaultValue int `json:"defaultValue"`
	Max          int `json:"max"`
	Min          int `json:"min"`
}

func (V IntegerRule) Validate(Value interface{}) error {

	return nil
}

/*
*
* 浮点数规则
*
 */
type FloatRule struct {
	DefaultValue float64 `json:"defaultValue"`
	Max          int     `json:"max"`
	Min          int     `json:"min"`
	Round        int     `json:"round"`
}

func (V FloatRule) Validate(Value interface{}) error {

	return nil
}

/*
*
* 布尔规则
*
 */
type BoolRule struct {
	DefaultValue bool   `json:"defaultValue"`
	TrueLabel    string `json:"trueLabel"`  // UI界面对应的label
	FalseLabel   string `json:"falseLabel"` // UI界面对应的label
}

func (V BoolRule) Validate(Value interface{}) error {

	return nil
}

/*
*
* 地理坐标规则
*
 */
type GeoRule struct {
	DefaultValue IoTPropertyGeo `json:"defaultValue"`
}

func (V GeoRule) Validate(Value interface{}) error {

	return nil
}

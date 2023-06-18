package typex

import "encoding/json"

// Rule type is for property store,
// XSource implements struct type is actually worker
type ModelType int

// 'T' means Type
const (
	T_INT32  ModelType = iota // int32
	T_FLOAT                   // float
	T_DOUBLE                  // double
	T_TEXT                    // pure text
	T_BOOL                    // boolean
	T_JSON                    // json
	T_BIN                     // byte
)

/*
* 数据模型, 例如某个Modbus电表可以支持读取电流/C 和电压/V参数:
*[
*    {
*        "name":"voltage",
*        "tag":"voltage",
*        "valueType":"float",
*        "value":220
*    }
*]
*
 */
type XDataModel struct {
	Name      string      `json:"name"`      // 字段名
	Tag       string      `json:"tag"`       // 标签
	ValueType ModelType   `json:"valueType"` // 值类型
	Value     interface{} `json:"value"`     // 具体的值
}

func (m XDataModel) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

//
//
// 创建资源的时候需要一个通用配置类
//
//

type XConfig struct {
	Type      string              `json:"type"` // 类型
	Engine    RuleX               `json:"-"`
	NewDevice func(RuleX) XDevice `json:"-"`
	NewSource func(RuleX) XSource `json:"-"`
	NewTarget func(RuleX) XTarget `json:"-"`
}

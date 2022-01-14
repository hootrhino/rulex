package typex

//
// Rule type is for property store,
// XResource implements struct type is actually worker
//
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
*        "name":"current",
*        "valueType":"float",
*        "value":5
*    },
*    {
*        "name":"volgate",
*        "valueType":"float",
*        "value":220
*    }
*]
*
 */
type XDataModel struct {
	Name      string      `json:"name"`
	ValueType ModelType   `json:"valueType"` // 值类型
	Value     interface{} `json:"value"`     // 具体的值
}

//
//
// 创建资源的时候需要一个通用配置类
// XConfig 可认为是接收参数的Form
// 前端可以拿来渲染界面(from v0.0.2)
//
//

type XConfig struct {
	Type    string        `json:"type"`    // 类型
	HelpTip string        `json:"helpTip"` // 关于这个配置的简介和帮助信息
	Views   []interface{} `json:"view"`    // 枚举，一般用来实现Select
}

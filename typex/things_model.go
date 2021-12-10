package typex

//
// Rule type is for property store,
// XResource implements struct type is actually worker
//
type ModelType int

// 'T' means Type
const (
	T_NUMBER  ModelType = 1
	T_STRING  ModelType = 2
	T_BOOLEAN ModelType = 3
	T_JSON    ModelType = 4
	T_BIN     ModelType = 5
	T_RAW     ModelType = 6
)

//
// 驱动的数据模型
//
type XDataModel struct {
	Type      ModelType
	Name      string
	MaxLength int
	MinLength int
}

//
//
// 创建资源的时候需要一个通用配置类
// XConfig 可认为是接收参数的Form
// 前端可以拿来渲染界面(from v0.0.2)
//
//
type ConfigType string

const (
	T_INPUT    ModelType = 1 //HTML input tag
	T_SELECT   ModelType = 2 //HTML select tag
	T_RADIO    ModelType = 3 //HTML radio tag
	T_SWITCH   ModelType = 4 //HTML switch tag
	T_CHECKBOX ModelType = 5 //HTML checkbox tag
)

type XConfig struct {
	UiType    string      `json:"uiType"`    // UI上显示的组件
	Field     string      `json:"field"`     // 字段名
	Title     string      `json:"title"`     // 标题
	Info      string      `json:"info"`      // 提示信息
	Hidden    bool        `json:"hidden"`    // 是否隐藏
	ValueType string      `json:"valueType"` // 值类型
	Value     interface{} `json:"value"`     // 具体的值
}

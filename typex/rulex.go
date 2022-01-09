package typex

import "sync"

//
// RuleX interface
//
type RuleX interface {
	//
	// 启动规则引擎
	//
	Start() *sync.Map
	//
	// 消息推到队列
	//
	PushQueue(QueueData) error
	PushInQueue(in *InEnd, data string) error
	PushOutQueue(out *OutEnd, data string) error
	//
	// 执行任务
	//
	Work(*InEnd, string) (bool, error)
	//
	// 获取配置
	//
	GetConfig(k string) interface{}
	//
	// 加载输入
	//
	LoadInEnd(*InEnd) error
	//
	// 获取输入
	//
	GetInEnd(string) *InEnd
	//
	// 保存输入
	//
	SaveInEnd(*InEnd)
	//
	// 删除输入
	//
	RemoveInEnd(string)
	//
	// 所有输入列表
	//
	AllInEnd() *sync.Map
	//
	// 加载输出
	//
	LoadOutEnd(*OutEnd) error
	//
	// 所有输出
	//
	AllOutEnd() *sync.Map
	//
	// 获取输出
	//
	GetOutEnd(string) *OutEnd
	//
	// 保存输出
	//
	SaveOutEnd(*OutEnd)
	//
	// 删除输出
	//
	RemoveOutEnd(string)
	//
	// 加载Hook
	//
	LoadHook(XHook) error
	//
	// 加载插件
	//
	LoadPlugin(XPlugin) error
	//
	// 所有插件列表
	//
	AllPlugins() *sync.Map
	//
	// 加载规则
	//
	LoadRule(*Rule) error
	//
	// 所有规则列表
	//
	AllRule() *sync.Map
	//
	// 获取规则
	//
	GetRule(id string) *Rule
	//
	// 删除规则
	//
	RemoveRule(uuid string)
	//
	// 运行 lua 回调
	//
	RunLuaCallbacks(*InEnd, string)
	//
	// 运行 hook
	//
	RunHooks(string)
	//
	// 获取版本
	//
	Version() Version

	//
	// 停止规则引擎
	//
	Stop()
}

//
// 拓扑接入点，比如 modbus 检测点等
// UUID: gyh9uo7uh7o67uijh
// Name: ModbusMeter001
// Alive: true
// Tag: modbus
//
type TopologyPoint struct {
	UUID   string `json:"uuid"`
	Parent string `json:"parent"`
	Name   string `json:"name"`
	Alive  bool   `json:"alive"`
	Tag    string `json:"tag"`
}

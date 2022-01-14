package typex

//
// XStatus for resource status
//
type XStatus struct {
	PointId    string // Input: Resource; Output: Target
	Enable     bool
	RuleEngine RuleX
}

//
// XResource: 终端资源, 比如实际上的 MQTT 客户端
//
type XResource interface {
	//
	// 测试资源是否可用
	//
	Test(inEndId string) bool

	//
	// 注册InEndID到资源
	//
	Register(inEndId string) error
	//
	// 启动资源
	//
	Start() error
	//
	// 资源是否被启用
	//
	Enabled() bool
	//
	// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
	//
	DataModels() []XDataModel
	//
	// 获取前端表单定义
	//
	Configs() XConfig
	//
	// 重载: 比如可以在重启的时候把某些数据保存起来
	//
	Reload()
	//
	// 挂起资源, 用来做暂停资源使用
	//
	Pause()
	//
	// 获取资源状态
	//
	Status() ResourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *InEnd
	//
	// 不经过规则引擎处理的直达数据接口
	//
	OnStreamApproached(data string) error
	//
	// 驱动接口, 通常用来和硬件交互
	//
	Driver() XExternalDriver
	//
	//
	//
	Topology() []TopologyPoint
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}

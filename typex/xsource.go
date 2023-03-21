package typex

import "context"

//
// XStatus for source status
//
type XStatus struct {
	PointId     string // Input: Source; Output: Target
	Enable      bool
	Ctx         context.Context
	CancelCTX   context.CancelFunc
	XDataModels []XDataModel `json:"dataModels" title:"数据模型" info:""`
	RuleEngine  RuleX
}

//
// XSource: 终端资源, 比如实际上的 MQTT 客户端
//
type XSource interface {
	//
	// 测试资源是否可用
	//
	Test(inEndId string) bool
	//
	// 用来初始化传递资源配置
	//
	Init(inEndId string, configMap map[string]interface{}) error
	//
	// 启动资源
	//
	Start(CCTX) error
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
	Configs() *XConfig
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
	Status() SourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *InEnd
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
	//
	// 来自外面的数据
	//
	DownStream([]byte) (int, error)
	//
	// 上行数据
	//
	UpStream([]byte) (int, error)
}

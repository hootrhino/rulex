package typex

import "context"

// XStatus for source status
type XStatus struct {
	PointId     string             // Input: Source; Output: Target
	Enable      bool               // 是否开启
	Ctx         context.Context    // context
	CancelCTX   context.CancelFunc // cancel
	XDataModels []XDataModel       // 数据模型
	RuleEngine  RuleX              // rulex
	Busy        bool               // 是否处于忙碌状态, 防止请求拥挤
}

// XSource: 终端资源, 比如实际上的 MQTT 客户端
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
	// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
	//
	DataModels() []XDataModel
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

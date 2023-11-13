package typex

//
// Stream from source and to target
//
type XTarget interface {
	//
	// 用来初始化传递资源配置
	//
	Init(outEndId string, configMap map[string]interface{}) error
	//
	// 启动资源
	//
	Start(CCTX) error
	//
	// 获取资源状态
	//
	Status() SourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *OutEnd
	//
	// 数据出口
	//
	To(data interface{}) (interface{}, error)
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}

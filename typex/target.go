package typex

//
// Stream from source and to target
//
type XTarget interface {
	//
	// 测试资源是否可用
	//
	Test(outEndId string) bool

	//
	// 注册InEndID到资源
	//
	Register(outEndId string) error
	//
	// 启动资源
	//
	Start() error
	//
	// 资源是否被启用
	//
	Enabled() bool
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
	Details() *OutEnd
	//
	//
	//
	Configs() *XConfig
	//
	// 数据出口
	//
	To(data interface{}) error
	//
	// 不经过规则引擎处理的直达数据
	//
	OnStreamApproached(data string) error
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}

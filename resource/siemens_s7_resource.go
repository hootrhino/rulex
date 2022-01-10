package resource

import (
	"rulex/typex"
	"rulex/utils"

	"github.com/robinson/gos7"
)

type siemensS7config struct {
	Host        string `json:"host"`
	Rack        int    `json:"rack"`
	Slot        int    `json:"slot"`
	Timeout     int    `json:"timeout"`
	IdleTimeout int    `json:"idleTimeout"`
}
type siemensS7Resource struct {
	typex.XStatus
	client gos7.Client
}

func NewSiemensS7Resource(e typex.RuleX) typex.XResource {
	s7 := siemensS7Resource{}
	s7.RuleEngine = e
	return &s7
}

//
// 测试资源是否可用
//
func (s7 *siemensS7Resource) Test(inEndId string) bool {
	return true
}

//
// 注册InEndID到资源
//
func (s7 *siemensS7Resource) Register(inEndId string) error {
	s7.PointId = inEndId
	return nil
}

//
// 启动资源
//
func (s7 *siemensS7Resource) Start() error {
	config := s7.RuleEngine.GetInEnd(s7.PointId).Config
	var mainConfig siemensS7config
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	handler := gos7.NewTCPClientHandler(mainConfig.Host, mainConfig.Rack, mainConfig.Slot)
	if err := handler.Connect(); err != nil {
		return err
	}
	client := gos7.NewClient(handler)
	s7.client = client
	return nil
}

//
// 资源是否被启用
//
func (s7 *siemensS7Resource) Enabled() bool {
	return true
}

//
// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
//
func (s7 *siemensS7Resource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

//
// 获取前端表单定义
//
func (s7 *siemensS7Resource) Configs() []typex.XConfig {
	return []typex.XConfig{}
}

//
// 重载: 比如可以在重启的时候把某些数据保存起来
//
func (s7 *siemensS7Resource) Reload() {

}

//
// 挂起资源, 用来做暂停资源使用
//
func (s7 *siemensS7Resource) Pause() {

}

//
// 获取资源状态
//
func (s7 *siemensS7Resource) Status() typex.ResourceState {
	return typex.UP
}

//
// 获取资源绑定的的详情
//
func (s7 *siemensS7Resource) Details() *typex.InEnd {
	return s7.RuleEngine.GetInEnd(s7.PointId)
}

//
// 不经过规则引擎处理的直达数据接口
//
func (s7 *siemensS7Resource) OnStreamApproached(data string) error {
	return nil
}

//
// 驱动接口, 通常用来和硬件交互
//
func (s7 *siemensS7Resource) Driver() typex.XExternalDriver {
	return nil
}

//
//
//
func (s7 *siemensS7Resource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 停止资源, 用来释放资源
//
func (s7 *siemensS7Resource) Stop() {
	if s7.client != nil {

	}
}

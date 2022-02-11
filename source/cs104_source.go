package source

import (
	"fmt"
	"rulex/typex"
	"rulex/utils"

	"github.com/thinkgos/go-iecp5/asdu"
	"github.com/thinkgos/go-iecp5/cs104"
)

type cs104Config struct {
	Host       string             `json:"host" validate:"required" title:"地址" info:""`
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	LogMode    bool               `json:"logMode" validate:"required" title:"日志" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}
type cs104Source struct {
	typex.XStatus

	host    string
	port    uint16
	logMode bool

	client *cs104.Client
}

func NewCs104Source() typex.XSource {
	cs := cs104Source{}
	return &cs
}

type cs104Client struct{}

func (cs104Client) InterrogationHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cs104Client) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (cs104Client) ReadHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cs104Client) TestCommandHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

func (cs104Client) ClockSyncHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (cs104Client) ResetProcessHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (cs104Client) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}
func (cs104Client) ASDUHandler(asdu.Connect, *asdu.ASDU) error {
	return nil
}

//
// 测试资源是否可用
//
func (cs *cs104Source) Test(inEndId string) bool {
	return true

}

//
// 注册InEndID到资源
//
func (cs *cs104Source) Register(inEndId string) error {
	cs.PointId = inEndId
	return nil

}

//
// 启动资源
//
func (cs *cs104Source) Start(cctx typex.CCTX) error {
	option := cs104.NewOption()
	config := cs.RuleEngine.GetInEnd(cs.PointId).Config
	var mainConfig cs104Config
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	cs.host = mainConfig.Host
	if err := option.AddRemoteServer(
		fmt.Sprintf("%s:%d", mainConfig.Host, mainConfig.Port),
	); err != nil {
		return err
	}
	cs.client = cs104.NewClient(cs104Client{}, option)
	cs.client.SetOnConnectHandler(func(c *cs104.Client) {
		c.SendStartDt()
	})
	cs.client.LogMode(cs.logMode)

	if err := cs.client.Start(); err != nil {
		return err
	}

	return nil
}

//
// 资源是否被启用
//
func (cs *cs104Source) Enabled() bool {
	return false
}

//
// 数据模型, 用来描述该资源支持的数据, 对应的是云平台的物模型
//
func (cs *cs104Source) DataModels() []typex.XDataModel {
	return nil

}

//
// 获取前端表单定义
//
func (cs *cs104Source) Configs() *typex.XConfig {
	return nil

}

//
// 重载: 比如可以在重启的时候把某些数据保存起来
//
func (cs *cs104Source) Reload() {
}

//
// 挂起资源, 用来做暂停资源使用
//
func (cs *cs104Source) Pause() {
}

//
// 获取资源状态
//
func (cs *cs104Source) Status() typex.SourceState {
	return typex.UP

}

//
// 获取资源绑定的的详情
//
func (cs *cs104Source) Details() *typex.InEnd {
	return nil

}

//
// 不经过规则引擎处理的直达数据接口
//
func (cs *cs104Source) OnStreamApproached(data string) error {
	return nil

}

//
// 驱动接口, 通常用来和硬件交互
//
func (cs *cs104Source) Driver() typex.XExternalDriver {
	return nil

}

//
//
//
func (cs *cs104Source) Topology() []typex.TopologyPoint {
	return nil

}

//
// 停止资源, 用来释放资源
//
func (cs *cs104Source) Stop() {
	cs.CancelCTX()
}

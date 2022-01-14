package resource

import (
	"fmt"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/nats-io/nats.go"
	"github.com/ngaut/log"
)

type natsConfig struct {
	User     string `json:"user" validate:"required" title:"连接账户" info:""`
	Password string `json:"password" validate:"required" title:"连接密码" info:""`
	Host     string `json:"host" validate:"required" title:"服务地址" info:""`
	Port     int32  `json:"port" validate:"required" title:"服务端口" info:""`
	Topic    string `json:"topic" validate:"required" title:"消息来源" info:""`
}
type natsResource struct {
	typex.XStatus
	user          string
	password      string
	host          string
	port          int32
	topic         string
	natsConnector *nats.Conn
}

func NewNatsResource(e typex.RuleX) typex.XResource {
	nt := &natsResource{}
	nt.RuleEngine = e
	return nt
}

func (nt *natsResource) Start() error {
	config := nt.RuleEngine.GetInEnd(nt.PointId).Config
	var mainConfig natsConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	nc, err := nats.Connect(fmt.Sprintf("%s:%v", mainConfig.Host, mainConfig.Port), func(o *nats.Options) error {
		o.User = mainConfig.User
		o.Password = mainConfig.Password
		return nil
	})
	if err != nil {
		return err
	} else {
		nt.natsConnector = nc
		nt.host = mainConfig.Host
		nt.port = mainConfig.Port
		nt.user = mainConfig.User
		nt.password = mainConfig.Password
		nt.topic = mainConfig.Topic
		//
		_, err := nt.natsConnector.Subscribe(nt.topic, func(msg *nats.Msg) {
			if nt.natsConnector != nil {
				nt.RuleEngine.Work(nt.RuleEngine.GetInEnd(nt.PointId), string(msg.Data))
			}
		})
		if err != nil {
			log.Error("NatsResource PushQueue error: ", err)
		}
		return nil
	}
}

// 测试资源状态
func (nt *natsResource) Test(inendId string) bool {
	if nt.natsConnector != nil {
		return nt.natsConnector.IsConnected()
	} else {
		return false
	}
}

// 先注册资源ID到出口
func (nt *natsResource) Register(inendId string) error {
	nt.PointId = inendId
	return nil
}

func (nt *natsResource) Enabled() bool {
	return true
}

func (nt *natsResource) Reload() {

}

func (nt *natsResource) Pause() {
}

func (nt *natsResource) Status() typex.ResourceState {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			return typex.UP
		}
	}
	return typex.DOWN

}

func (nt *natsResource) Details() *typex.InEnd {
	return nt.RuleEngine.GetInEnd(nt.PointId)
}

//--------------------------------------------------------
// To: 数据出口
//--------------------------------------------------------

func (nt *natsResource) OnStreamApproached(data string) error {
	return nil
}

func (nt *natsResource) Stop() {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			nt.natsConnector.Drain()
			nt.natsConnector.Close()
			nt.natsConnector = nil
		}
	}
}
func (nt *natsResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("NATS_SERVER", "", natsConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

func (nt *natsResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (nt *natsResource) Driver() typex.XExternalDriver {
	return nil
}

//
// 拓扑
//
func (*natsResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

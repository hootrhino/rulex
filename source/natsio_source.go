package source

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
type natsSource struct {
	typex.XStatus
	user          string
	password      string
	host          string
	port          int32
	topic         string
	natsConnector *nats.Conn
}

func NewNatsSource(e typex.RuleX) typex.XSource {
	nt := &natsSource{}
	nt.RuleEngine = e
	return nt
}

func (nt *natsSource) Start(cctx typex.CCTX) error {
	nt.Ctx = cctx.Ctx
	nt.CancelCTX = cctx.CancelCTX

	config := nt.RuleEngine.GetInEnd(nt.PointId).Config
	var mainConfig natsConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
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
			log.Error("NatsSource PushQueue error: ", err)
		}
		return nil
	}
}

// 测试资源状态
func (nt *natsSource) Test(inendId string) bool {
	if nt.natsConnector != nil {
		return nt.natsConnector.IsConnected()
	} else {
		return false
	}
}

// 先注册资源ID到出口
func (nt *natsSource) Register(inendId string) error {
	nt.PointId = inendId
	return nil
}

func (nt *natsSource) Enabled() bool {
	return true
}

func (nt *natsSource) Reload() {

}

func (nt *natsSource) Pause() {
}

func (nt *natsSource) Status() typex.SourceState {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			return typex.UP
		}
	}
	return typex.DOWN

}

func (nt *natsSource) Details() *typex.InEnd {
	return nt.RuleEngine.GetInEnd(nt.PointId)
}

//--------------------------------------------------------
// To: 数据出口
//--------------------------------------------------------

func (nt *natsSource) OnStreamApproached(data string) error {
	return nil
}

func (nt *natsSource) Stop() {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			nt.natsConnector.Drain()
			nt.natsConnector.Close()
			nt.natsConnector = nil
		}
	}
	nt.CancelCTX()
}
func (nt *natsSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.NATS_SERVER, "NATS_SERVER", natsConfig{})
}

func (nt *natsSource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}
func (nt *natsSource) Driver() typex.XExternalDriver {
	return nil
}

//
// 拓扑
//
func (*natsSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

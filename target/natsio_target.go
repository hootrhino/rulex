package target

import (
	"errors"
	"fmt"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/nats-io/nats.go"
)

type natsConfig struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     string `json:"port" validate:"required"`
	Topic    string `json:"topic" validate:"required"`
}
type natsTarget struct {
	typex.XStatus
	user          string
	password      string
	host          string
	port          string
	topic         string
	natsConnector *nats.Conn
}

func NewNatsTarget(e typex.RuleX) typex.XTarget {
	nt := &natsTarget{}
	nt.RuleEngine = e
	return nt
}

func (nt *natsTarget) Start(cctx typex.CCTX) error {
	nt.Ctx = cctx.Ctx
	nt.CancelCTX = cctx.CancelCTX
	outEnd := nt.RuleEngine.GetOutEnd(nt.PointId)
	config := outEnd.Config
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
		return nil
	}
}

// 测试资源状态
func (nt *natsTarget) Test(outEndId string) bool {
	if nt.natsConnector != nil {
		return nt.natsConnector.IsConnected()
	} else {
		return false
	}
}

// 先注册资源ID到出口
func (nt *natsTarget) Register(outEndId string) error {
	nt.PointId = outEndId
	return nil
}
func (nt *natsTarget) Init(outEndId string, cfg map[string]interface{}) error {
	nt.PointId = outEndId
	return nil
}
func (nt *natsTarget) Enabled() bool {
	return true
}

func (nt *natsTarget) Reload() {

}

func (nt *natsTarget) Pause() {
}

func (nt *natsTarget) Status() typex.SourceState {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			return typex.SOURCE_UP
		}
	}
	return typex.SOURCE_DOWN

}

func (nt *natsTarget) Details() *typex.OutEnd {
	return nt.RuleEngine.GetOutEnd(nt.PointId)
}

//--------------------------------------------------------
// To: 数据出口
//--------------------------------------------------------
func (nt *natsTarget) To(data interface{}) (interface{}, error) {
	if nt.natsConnector != nil {
		return nil, nt.natsConnector.Publish(nt.topic, []byte((data.(string))))
	}
	return nil, errors.New("natsConnector is nil")
}

func (nt *natsTarget) Stop() {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			nt.natsConnector.Drain()
			nt.natsConnector.Close()
			nt.natsConnector = nil
		}
	}
	nt.CancelCTX()
}

/*
*
* 配置
*
 */
func (*natsTarget) Configs() *typex.XConfig {
	return core.GenOutConfig(typex.NATS_TARGET, "NATS_TARGET", httpConfig{})
}

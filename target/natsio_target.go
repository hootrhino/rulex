package target

import (
	"errors"
	"fmt"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/nats-io/nats.go"
)

type natsTarget struct {
	typex.XStatus
	natsConnector *nats.Conn
	mainConfig    common.NatsConfig
	status        typex.SourceState
}

func NewNatsTarget(e typex.RuleX) typex.XTarget {
	nt := &natsTarget{}
	nt.RuleEngine = e
	nt.mainConfig=common.NatsConfig{}
	return nt
}
func (nt *natsTarget) Init(outEndId string, configMap map[string]interface{}) error {
	nt.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &nt.mainConfig); err != nil {
		return err
	}
	return nil
}
func (nt *natsTarget) Start(cctx typex.CCTX) error {
	nt.Ctx = cctx.Ctx
	nt.CancelCTX = cctx.CancelCTX

	nc, err := nats.Connect(fmt.Sprintf("%s:%v", nt.mainConfig.Host, nt.mainConfig.Port), func(o *nats.Options) error {
		o.User = nt.mainConfig.Username
		o.Password = nt.mainConfig.Password
		return nil
	})
	if err != nil {
		return err
	} else {
		nt.natsConnector = nc
		nt.status = typex.SOURCE_UP
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

func (nt *natsTarget) Enabled() bool {
	return true
}

func (nt *natsTarget) Reload() {

}

func (nt *natsTarget) Pause() {
}

func (nt *natsTarget) Status() typex.SourceState {
	return nt.status
}

func (nt *natsTarget) Details() *typex.OutEnd {
	return nt.RuleEngine.GetOutEnd(nt.PointId)
}

// --------------------------------------------------------
// To: 数据出口
// --------------------------------------------------------
func (nt *natsTarget) To(data interface{}) (interface{}, error) {
	if nt.natsConnector != nil {
		return nil, nt.natsConnector.Publish(nt.mainConfig.Topic, []byte((data.(string))))
	}
	return nil, errors.New("nats Connector is nil")
}

func (nt *natsTarget) Stop() {
	nt.status = typex.SOURCE_STOP
	nt.CancelCTX()
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			nt.natsConnector.Drain()
			nt.natsConnector.Close()
			nt.natsConnector = nil
		}
	}

}

/*
*
* 配置
*
 */
func (*natsTarget) Configs() *typex.XConfig {
	return typex.GenOutConfig(typex.NATS_TARGET, "NATS_TARGET", common.NatsConfig{})
}

package source

import (
	"fmt"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/nats-io/nats.go"
)

type natsSource struct {
	typex.XStatus
	natsConnector *nats.Conn
	mainConfig    common.NatsConfig
}

func NewNatsSource(e typex.RuleX) typex.XSource {
	nt := &natsSource{}
	nt.RuleEngine = e
	return nt
}

func (nt *natsSource) Start(cctx typex.CCTX) error {
	nt.Ctx = cctx.Ctx
	nt.CancelCTX = cctx.CancelCTX

	nc, err := nats.Connect(fmt.Sprintf("%s:%v", nt.mainConfig.Host, nt.mainConfig.Port), func(o *nats.Options) error {
		o.User = nt.mainConfig.User
		o.Password = nt.mainConfig.Password
		return nil
	})
	if err != nil {
		return err
	} else {
		nt.natsConnector = nc
		//
		_, err := nt.natsConnector.Subscribe(nt.mainConfig.Topic, func(msg *nats.Msg) {
			if nt.natsConnector != nil {
				work, err1 := nt.RuleEngine.WorkInEnd(nt.RuleEngine.GetInEnd(nt.PointId), string(msg.Data))
				if !work {
					glogger.GLogger.Error(err1)
				}
			}
		})
		if err != nil {
			glogger.GLogger.Error("NatsSource PushQueue error: ", err)
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

func (nt *natsSource) Init(inEndId string, configMap map[string]interface{}) error {
	nt.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &nt.mainConfig); err != nil {
		return err
	}
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
			return typex.SOURCE_UP
		}
	}
	return typex.SOURCE_DOWN

}

func (nt *natsSource) Details() *typex.InEnd {
	return nt.RuleEngine.GetInEnd(nt.PointId)
}

//--------------------------------------------------------
// To: 数据出口
//--------------------------------------------------------

func (nt *natsSource) Stop() {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			err := nt.natsConnector.Drain()
			if err != nil {
				glogger.GLogger.Error(err)
				return
			}
			nt.natsConnector.Close()
			nt.natsConnector = nil
		}
	}
	nt.CancelCTX()
}
func (nt *natsSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.NATS_SERVER, "NATS_SERVER", common.NatsConfig{})
}

func (nt *natsSource) DataModels() []typex.XDataModel {
	return nt.XDataModels
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

//
// 来自外面的数据
//
func (*natsSource) DownStream([]byte) {}

//
// 上行数据
//
func (*natsSource) UpStream() {}

package source

import (
	"fmt"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/nats-io/nats.go"
)

type natsSource struct {
	typex.XStatus
	natsConnector *nats.Conn
	mainConfig    common.NatsConfig
	status        typex.SourceState
}

func NewNatsSource(e typex.RuleX) typex.XSource {
	nt := &natsSource{}
	nt.RuleEngine = e
	return nt
}

func (nt *natsSource) Start(cctx typex.CCTX) error {
	nt.Ctx = cctx.Ctx
	nt.CancelCTX = cctx.CancelCTX

	nc, err := nats.Connect(fmt.Sprintf("%s:%v", nt.mainConfig.Host, nt.mainConfig.Port),
		func(o *nats.Options) error {
			o.User = nt.mainConfig.Username
			o.Password = nt.mainConfig.Password
			o.Name = "rulex-nats-source"
			o.AllowReconnect = true
			o.ReconnectWait = 5 * time.Second
			return nil
		},
	)

	nc.SetDisconnectHandler(func(c *nats.Conn) {
		glogger.GLogger.Warn("connection disconnected")
	})
	nc.SetReconnectHandler(func(c *nats.Conn) {
		glogger.GLogger.Warn("connection reconnect")
		time.Sleep(2 * time.Second)
		nt.subscribeNats()
	})

	if err != nil {
		return err
	} else {
		nt.natsConnector = nc
		//
		nt.subscribeNats()
		nt.status = typex.SOURCE_UP
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

func (nt *natsSource) Status() typex.SourceState {
	return nt.status

}

func (nt *natsSource) Details() *typex.InEnd {
	return nt.RuleEngine.GetInEnd(nt.PointId)
}

//--------------------------------------------------------
// To: 数据出口
//--------------------------------------------------------

func (nt *natsSource) Stop() {
	nt.status = typex.SOURCE_STOP
	if nt.CancelCTX != nil {
		nt.CancelCTX()
	}
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

}

func (nt *natsSource) DataModels() []typex.XDataModel {
	return nt.XDataModels
}
func (nt *natsSource) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*natsSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*natsSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*natsSource) UpStream([]byte) (int, error) {
	return 0, nil
}

// --------------------------------------------------------------------------------------------------
func (nt *natsSource) subscribeNats() {
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
}

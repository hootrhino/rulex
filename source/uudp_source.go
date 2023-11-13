package source

import (
	"context"
	"net"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type udpSource struct {
	typex.XStatus
	uDPConn    *net.UDPConn
	mainConfig common.RULEXUdpConfig
	status     typex.SourceState
}

func NewUdpInEndSource(e typex.RuleX) typex.XSource {
	u := udpSource{}
	u.mainConfig = common.RULEXUdpConfig{}
	u.RuleEngine = e
	return &u
}
func (u *udpSource) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX

	addr := &net.UDPAddr{IP: net.ParseIP(u.mainConfig.Host), Port: u.mainConfig.Port}
	var err error
	if u.uDPConn, err = net.ListenUDP("udp", addr); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	u.status = typex.SOURCE_UP
	go func(ctx context.Context, u1 *udpSource) {
		if u.mainConfig.MaxDataLength == 0 {
			u.mainConfig.MaxDataLength = 4096
		}
		data := make([]byte, u.mainConfig.MaxDataLength)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			n, remoteAddr, err := u1.uDPConn.ReadFromUDP(data)
			if err != nil {
				glogger.GLogger.Error(err.Error())
				continue
			}
			work, err := u.RuleEngine.WorkInEnd(u.RuleEngine.GetInEnd(u.PointId), string(data[:n]))
			if !work {
				glogger.GLogger.Error(err)
				continue
			}
			// return ok
			_, err = u1.uDPConn.WriteToUDP([]byte("ok"), remoteAddr)
			if err != nil {
				glogger.GLogger.Error(err)
			}
		}
	}(u.Ctx, u)
	u.status = typex.SOURCE_UP
	glogger.GLogger.Infof("UDP source started on [%v]:%v", u.mainConfig.Host, u.mainConfig.Port)
	return nil

}

func (u *udpSource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *udpSource) Test(inEndId string) bool {
	return true
}

func (u *udpSource) Init(inEndId string, configMap map[string]interface{}) error {
	u.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &u.mainConfig); err != nil {
		return err
	}
	return nil
}

func (u *udpSource) DataModels() []typex.XDataModel {
	return u.XDataModels
}

func (u *udpSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (u *udpSource) Stop() {
	u.status = typex.SOURCE_STOP
	if u.CancelCTX != nil {
		u.CancelCTX()
	}
	if u.uDPConn != nil {
		err := u.uDPConn.Close()
		if err != nil {
			glogger.GLogger.Error(err)
		}
	}

}
func (*udpSource) Driver() typex.XExternalDriver {
	return nil
}

// 拓扑
func (*udpSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*udpSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*udpSource) UpStream([]byte) (int, error) {
	return 0, nil
}

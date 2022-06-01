package source

import (
	"context"
	"net"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
)

type udpSource struct {
	typex.XStatus
	uDPConn *net.UDPConn
}
type udpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址" info:""`
	Port          int    `json:"port" validate:"required" title:"服务端口" info:""`
	MaxDataLength int    `json:"maxDataLength" validate:"required" title:"最大数据包" info:""`
}

func NewUdpInEndSource(e typex.RuleX) *udpSource {
	u := udpSource{}
	u.RuleEngine = e
	return &u
}
func (u *udpSource) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX

	config := u.RuleEngine.GetInEnd(u.PointId).Config
	var mainConfig udpConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	addr := &net.UDPAddr{IP: net.ParseIP(mainConfig.Host), Port: mainConfig.Port}
	var err error
	if u.uDPConn, err = net.ListenUDP("udp", addr); err != nil {
		log.Error(err)
		return err
	}

	go func(c context.Context, u1 *udpSource) {
		data := make([]byte, mainConfig.MaxDataLength)
		for {
			n, remoteAddr, err := u1.uDPConn.ReadFromUDP(data)
			if err != nil {
				log.Error(err.Error())
			} else {
				// log.Infof("Receive udp data:<%s> %s\n", remoteAddr, data[:n])
				work, err := u.RuleEngine.WorkInEnd(u.RuleEngine.GetInEnd(u.PointId), string(data[:n]))
				if !work {
					log.Error(err)
				}
				// return ok
				_, err = u1.uDPConn.WriteToUDP([]byte("ok"), remoteAddr)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}(u.Ctx, u)
	log.Infof("UDP source started on [%v]:%v", mainConfig.Host, mainConfig.Port)
	return nil

}

func (u *udpSource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *udpSource) Test(inEndId string) bool {
	return true
}

func (u *udpSource) Init(inEndId string, cfg map[string]interface{}) error {
	u.PointId = inEndId
	return nil
}
func (u *udpSource) Enabled() bool {
	return true
}

func (u *udpSource) DataModels() []typex.XDataModel {
	return u.XDataModels
}

func (u *udpSource) Reload() {
}

func (u *udpSource) Pause() {
}

func (u *udpSource) Status() typex.SourceState {
	return typex.UP
}

func (u *udpSource) Stop() {
	if u.uDPConn != nil {
		err := u.uDPConn.Close()
		if err != nil {
			log.Error(err)
		}
	}
	u.CancelCTX()
}
func (*udpSource) Driver() typex.XExternalDriver {
	return nil
}
func (*udpSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.RULEX_UDP, "RULEX_UDP", udpConfig{})
}

//
// 拓扑
//
func (*udpSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

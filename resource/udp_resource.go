package resource

import (
	"context"
	"net"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
)

type udpResource struct {
	typex.XStatus
	uDPConn *net.UDPConn
}
type udpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址" info:""`
	Port          int    `json:"port" validate:"required" title:"服务端口" info:""`
	MaxDataLength int    `json:"maxDataLength" validate:"required" title:"最大数据包" info:""`
}

func NewUdpInEndResource(e typex.RuleX) *udpResource {
	u := udpResource{}
	u.RuleEngine = e
	return &u
}
func (u *udpResource) Start() error {
	config := u.RuleEngine.GetInEnd(u.PointId).Config
	var mainConfig udpConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	addr := &net.UDPAddr{IP: net.ParseIP(mainConfig.Host), Port: mainConfig.Port}
	var err error
	if u.uDPConn, err = net.ListenUDP("udp", addr); err != nil {
		log.Error(err)
		return err
	}
	go func(c context.Context, u1 *udpResource) {
		data := make([]byte, mainConfig.MaxDataLength)
		for {
			n, remoteAddr, err := u1.uDPConn.ReadFromUDP(data)
			if err != nil {
				log.Error(err.Error())
			} else {
				// log.Infof("Receive udp data:<%s> %s\n", remoteAddr, data[:n])
				work, err := u.RuleEngine.Work(u.RuleEngine.GetInEnd(u.PointId), string(data[:n]))
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
	}(context.Background(), u)
	log.Infof("UDP resource started on [%v]:%v", mainConfig.Host, mainConfig.Port)
	return nil

}
func (u *udpResource) OnStreamApproached(data string) error {
	work, err := u.RuleEngine.Work(u.RuleEngine.GetInEnd(u.PointId), data)
	if !work {
		return err
	}
	return nil
}
func (u *udpResource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *udpResource) Test(inEndId string) bool {
	return true
}

func (u *udpResource) Register(inEndId string) error {
	u.PointId = inEndId
	return nil
}

func (u *udpResource) Enabled() bool {
	return true
}

func (u *udpResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (u *udpResource) Reload() {
}

func (u *udpResource) Pause() {
}

func (u *udpResource) Status() typex.ResourceState {
	return typex.UP
}

func (u *udpResource) Stop() {
	if u.uDPConn != nil {
		u.uDPConn.Close()
	}
}
func (*udpResource) Driver() typex.XExternalDriver {
	return nil
}
func (*udpResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("RULEX_UDP", "", udpConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

//
// 拓扑
//
func (*udpResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

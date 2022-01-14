package resource

import (
	"context"
	"net"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
)

type UdpResource struct {
	typex.XStatus
	uDPConn *net.UDPConn
}
type udpConfig struct {
	Host          string `json:"host" validate:"required" title:"服务地址" info:""`
	Port          int    `json:"port" validate:"required" title:"服务端口" info:""`
	MaxDataLength int    `json:"maxDataLength" validate:"required" title:"最大数据包" info:""`
}

func NewUdpInEndResource(e typex.RuleX) *UdpResource {
	u := UdpResource{}
	u.RuleEngine = e
	return &u
}
func (u *UdpResource) Start() error {
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
	go func(c context.Context, u1 *UdpResource) {
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
func (u *UdpResource) OnStreamApproached(data string) error {
	work, err := u.RuleEngine.Work(u.RuleEngine.GetInEnd(u.PointId), data)
	if !work {
		return err
	}
	return nil
}
func (u *UdpResource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *UdpResource) Test(inEndId string) bool {
	return true
}

func (u *UdpResource) Register(inEndId string) error {
	u.PointId = inEndId
	return nil
}

func (u *UdpResource) Enabled() bool {
	return true
}

func (u *UdpResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (u *UdpResource) Reload() {
}

func (u *UdpResource) Pause() {
}

func (u *UdpResource) Status() typex.ResourceState {
	return typex.UP
}

func (u *UdpResource) Stop() {
	if u.uDPConn != nil {
		u.uDPConn.Close()
	}
}
func (*UdpResource) Driver() typex.XExternalDriver {
	return nil
}
func (*UdpResource) Configs() typex.XConfig {
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
func (*UdpResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

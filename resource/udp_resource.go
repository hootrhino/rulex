package resource

import (
	"context"
	"net"
	"rulex/typex"

	"github.com/ngaut/log"
)

type UdpResource struct {
	typex.XStatus
	uDPConn  *net.UDPConn
	CanWrite bool
}

func NewUdpInEndResource(inEndId string, e typex.RuleX) *UdpResource {
	u := UdpResource{
		CanWrite: false,
	}
	u.RuleEngine = e
	return &u
}
func (u *UdpResource) Start() error {
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 2584}
	var err error
	if u.uDPConn, err = net.ListenUDP("udp", addr); err != nil {
		log.Error(err)
		return err
	}
	u.CanWrite = true
	go func(c context.Context, u1 *UdpResource) {
		data := make([]byte, 1024)
		for u1.CanWrite {
			n, remoteAddr, err := u1.uDPConn.ReadFromUDP(data)
			if err != nil {
				if u1.CanWrite {
					log.Error(err.Error())
				}
			} else {
				// fmt.Printf("Receive udp data:<%s> %s\n", remoteAddr, data[:n])
				work, err := u.RuleEngine.Work(u.RuleEngine.GetInEnd(u.PointId), string(data[:n]))
				if !work {
					log.Error(err)
				}
				_, err = u1.uDPConn.WriteToUDP([]byte("ok"), remoteAddr)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}(context.Background(), u)
	log.Info("UDP resource started on [0.0.0.0]:2584")
	return nil

}
func (u *UdpResource) OnStreamApproached(data string) error {
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
	u.CanWrite = false
	if u.uDPConn != nil {
		u.uDPConn.Close()
	}
}
func (*UdpResource) Driver() typex.XExternalDriver {
	return nil
}
func (*UdpResource) Configs() []typex.XConfig {
	return []typex.XConfig{}
}

//
// 拓扑
//
func (*UdpResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

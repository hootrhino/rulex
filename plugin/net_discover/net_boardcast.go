package netdiscover

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"gopkg.in/ini.v1"
)

/*
*
* 网关之间的发现协议,让同一个网络内的网关能发现对方.这个功能暂时没想好用来做什么.
*
 */
//
// udp://hostname@host:port
//
type gwnode struct {
	Host     string
	Hostname string
	Port     string
	Timeout  int32
}
type NetDiscover struct {
	Neighbors map[string]gwnode // 邻居
}

func NewNetDiscover() typex.XPlugin {
	return &NetDiscover{}
}

func (dm *NetDiscover) Init(config *ini.Section) error {
	return nil
}

func (dm *NetDiscover) Start(typex.RuleX) error {
	// 超时管理器
	go func() {
		for {
			time.Sleep(5 * time.Second)
			for _, nb := range dm.Neighbors {
				if nb.Timeout == 0 {
					// Delete node
					continue
				}
				atomic.AddInt32(&nb.Timeout, -1)
			}

		}
	}()
	go func() {
		listener, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IPv4zero,
			Port: 1994,
		})
		if err != nil {
			glogger.GLogger.Fatal(err)
			return
		}
		data := make([]byte, 1024)
		for {
			n, remoteAddr, err := listener.ReadFromUDP(data)
			if err != nil {
				glogger.GLogger.Infof("read remote: %s, err: %s", remoteAddr.String(), err)
				continue
			}
			// 请求组网的包 CAST:IP-Addr
			// CAST:192.168.001.001
			if n == 4 {
				// 任何一个网关加入集群都会广播
				if string(data[:4]) == ("CAST") {
					glogger.GLogger.Infof("Received CAST from:%s", remoteAddr.String())
				}
				// 集群内同步网关列表
				if string(data[:4]) == ("SYNC") {
					glogger.GLogger.Infof("Received SYNC from:%s", remoteAddr.String())
				}
			}

		}
	}()
	return nil
}
func (dm *NetDiscover) Stop() error {
	return nil
}

func (hh *NetDiscover) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "NetDiscover",
		Version:  "0.0.1",
		Homepage: "www.github.com/i4de/rulex",
		HelpLink: "www.github.com/i4de/rulex",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}
/*
*
* 服务调用接口
*
 */
 func (cs *NetDiscover) Service(arg typex.ServiceArg) error {
	return nil
}

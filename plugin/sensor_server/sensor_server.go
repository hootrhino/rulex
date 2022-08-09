package sensor_server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/i4de/rulex/typex"
	"github.com/plgd-dev/kit/v2/log"
	"gopkg.in/ini.v1"
)

type SensorServer struct {
	ctxMain    context.Context
	cancelMain context.CancelFunc
	tcpPort    int
	httpPort   int
}

func NewSensorServer() *SensorServer {
	return &SensorServer{}
}

func (dm *SensorServer) Init(config *ini.Section) error {
	return nil
}

func (dm *SensorServer) Start(typex.RuleX) error {
	return nil
}
func (dm *SensorServer) Stop() error {
	return nil
}

func (hh *SensorServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "GenericSensorServer",
		Version:  "0.0.1",
		Homepage: "www.rulexgw.io",
		HelpLink: "www.rulexgw.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func formatUrl(v int) string {
	return fmt.Sprintf("0.0.0.0:%d", v)
}
func startSServer(ctxMain context.Context, cancelMain context.CancelFunc) error {
	var listener net.Listener

	var err error
	listener, err = net.Listen("tcp", formatUrl(6690))
	if err != nil {
		log.Fatal("Error listening:", err)
		return err
	}
	defer listener.Close()

	for {
		select {
		case <-ctxMain.Done():
			{
				cancelMain()
				return nil
			}
		default:
			{
			}
		}
		peerConn, err := listener.Accept()
		if err != nil {
			log.Error("Error Listener Accept: ", err)
			continue
		}
		session := NewSession(peerConn)
		ctx, cancel := context.WithCancel(ctxMain)
		go waitForAuth(ctx, cancel, session)

	}

}

/*
*
* 等待认证
*
 */
func waitForAuth(ctx context.Context, cancel context.CancelFunc, s Session) {
	// 等待3秒内是否收到认证报文
	buffer := make([]byte, 64)
	for {
		s.Transport.SetDeadline(time.Now().Add(3 * time.Second))
		n, err := s.Transport.Read(buffer)
		if err != nil {
			log.Error(err)
			return
		}
		id := string(buffer[:n])
		log.Debug("Sensor ready to auth:", id)

	}

}

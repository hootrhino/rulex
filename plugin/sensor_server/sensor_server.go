package sensor_server

import (
	"context"
	"errors"
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
	ctx, cancel := context.WithCancel(context.Background())
	return &SensorServer{ctxMain: ctx, cancelMain: cancel}
}

func (dm *SensorServer) Init(config *ini.Section) error {
	k1, _ := config.GetKey("tcp_port")
	dm.tcpPort = k1.MustInt(30001)
	k2, _ := config.GetKey("http_port")
	dm.httpPort = k2.MustInt(60001)
	return nil
}

func (dm *SensorServer) Start(typex.RuleX) error {
	return dm.startSServer(dm.ctxMain, dm.cancelMain)
}
func (dm *SensorServer) Stop() error {
	dm.cancelMain()
	return nil
}

func (hh *SensorServer) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		Name:     "GenericSensorServer",
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
func (cs *SensorServer) Service(arg typex.ServiceArg) error {
	return nil
}

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------
func formatUrl(v int) string {
	return fmt.Sprintf("0.0.0.0:%d", v)
}
func (hh *SensorServer) startSServer(ctxMain context.Context, cancelMain context.CancelFunc) error {
	var listener net.Listener

	var err error
	listener, err = net.Listen("tcp", formatUrl(hh.tcpPort))
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
* 等待认证: 传感器发送的第一个包必须为ID, 最大不能超过64字节
*
 */
func waitForAuth(ctx context.Context, cancel context.CancelFunc, session Session) {
	buffer := [64]byte{}
	for {
		session.Transport.SetDeadline(time.Now().Add(5 * time.Second))
		n, err := session.Transport.Read(buffer[:])
		session.Transport.SetDeadline(time.Time{})
		if err != nil {
			log.Error(err)
			session.Transport.Close()
			goto END
		}
		sn := string(buffer[:n])
		log.Debug("Sensor ready to auth:", sn)
		// 这里应该加入认证的逻辑 但是目前默认传ID就表示认证成功
		if sn != "" {
			isensor := NewSensor(session)
			isensor.OnRegister(sn)
			if err := isensor.OnRegister(sn); err != nil {
				log.Error(err)
				session.Transport.Close()
				goto END
			} else {
				isensor.OnLine()
				startWorker(ctx, cancel, isensor)
				goto END
			}
		} else {
			log.Error(errors.New("must set sensor sn"))
			session.Transport.Close()
			goto END
		}
	}
END:
	{
	}
}

/*
*
* 启动工作进程
*
 */
func startWorker(ctx context.Context, cancel context.CancelFunc, isensor ISensor) {
	worker := SensorWorker{Ctx: ctx, Cancel: cancel, isensor: isensor}
	worker.Run()
}

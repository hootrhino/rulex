// 拖车驱动
package driver

import (
	"context"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/sidecar"
	"github.com/i4de/rulex/typex"

	"google.golang.org/grpc"
)

type SideCarDriver struct {
	state      typex.DriverState
	RuleEngine typex.RuleX
	client     sidecar.SidecarClient
	config     map[string]string
}

func NewSideCarDriver(e typex.RuleX, grpcConn *grpc.ClientConn) typex.XExternalDriver {
	sideCarDriver := &SideCarDriver{
		state:      typex.DRIVER_STOP,
		RuleEngine: e,
		client:     sidecar.NewSidecarClient(grpcConn),
	}
	return sideCarDriver

}
func (sc *SideCarDriver) Test() error {
	if err := sc.t(); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (sc *SideCarDriver) Init(config map[string]string) error {
	sc.config = config
	_, err := sc.client.Init(context.Background(), &sidecar.Config{
		Kv: sc.config,
	})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (sc *SideCarDriver) Work() error {
	_, err := sc.client.Start(context.Background(), &sidecar.Request{})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (sc *SideCarDriver) State() typex.DriverState {
	if sc.t() != nil {
		return typex.DRIVER_STOP
	}
	return typex.DRIVER_UP
}

/*
*
* 读取
*
 */
func (sc *SideCarDriver) Read(data []byte) (int, error) {
	response, err := sc.client.Read(context.Background(), &sidecar.ReadRequest{})
	if err != nil {
		glogger.GLogger.Error(err)
		return 0, err
	}
	copy(data, response.GetData())
	return len(response.Data), nil
}

/*
*
* 写入
*
 */
func (sc *SideCarDriver) Write(data []byte) (int, error) {
	response, err := sc.client.Write(context.Background(), &sidecar.WriteRequest{
		Data: data,
	})
	if err != nil {
		glogger.GLogger.Error(err)
		return 0, err
	}
	return int(response.Code), nil
}

//---------------------------------------------------
func (sc *SideCarDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "SIDECAR-DRIVER",
		Type:        "SIDECAR",
		Description: "SIDECAR 通用GRPC协议驱动",
	}
}

func (sc *SideCarDriver) Stop() error {
	_, err := sc.client.Stop(context.Background(), &sidecar.Request{})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

//----------------------------------------------------------------------
// 私有函数
//----------------------------------------------------------------------
func (sc *SideCarDriver) t() error {
	_, err := sc.client.Status(context.Background(), &sidecar.Request{})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

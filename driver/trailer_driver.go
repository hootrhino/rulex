// 拖车驱动
package driver

import (
	"context"

	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	"google.golang.org/grpc"
)

type TrailerDriver struct {
	state      typex.DriverState
	RuleEngine typex.RuleX
	client     trailer.TrailerClient
	config     map[string]string
}

func NewTrailerDriver(e typex.RuleX, grpcConn *grpc.ClientConn) typex.XExternalDriver {
	TrailerDriver := &TrailerDriver{
		state:      typex.DRIVER_STOP,
		RuleEngine: e,
		client:     trailer.NewTrailerClient(grpcConn),
	}
	return TrailerDriver

}
func (sc *TrailerDriver) Test() error {
	if err := sc.t(); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (sc *TrailerDriver) Init(config map[string]string) error {
	sc.config = config
	_, err := sc.client.Init(context.Background(), &trailer.Config{
		Kv: []byte{},
	})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (sc *TrailerDriver) Work() error {
	_, err := sc.client.Start(context.Background(), &trailer.Request{})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (sc *TrailerDriver) State() typex.DriverState {
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
func (sc *TrailerDriver) Read(cmd []byte, data []byte) (int, error) {
	response, err := sc.client.Service(context.Background(),
		&trailer.ServiceRequest{Cmd: cmd, Args: data})
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
func (sc *TrailerDriver) Write(cmd []byte, data []byte) (int, error) {
	response, err := sc.client.Service(context.Background(),
		&trailer.ServiceRequest{Cmd: cmd, Args: data})
	if err != nil {
		glogger.GLogger.Error(err)
		return 0, err
	}
	return int(response.Code), nil
}

// ---------------------------------------------------
func (sc *TrailerDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Trailer-DRIVER",
		Type:        "Trailer",
		Description: "Trailer 通用GRPC协议驱动",
	}
}

func (sc *TrailerDriver) Stop() error {
	_, err := sc.client.Stop(context.Background(), &trailer.Request{})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

// ----------------------------------------------------------------------
// 私有函数
// ----------------------------------------------------------------------
func (sc *TrailerDriver) t() error {
	_, err := sc.client.Status(context.Background(), &trailer.Request{})
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

// 拖车驱动
package driver

//   // 初始化, 主要是为了传配置进去
//   rpc Init (Config) returns (Request) {}
//   // 启动
//   rpc Start (Request) returns (Response) {}
//   // 获取状态
//   rpc Status (Request) returns (Response) {}
//   // 读数据
//   rpc Read (ReadRequest) returns (ReadResponse) {}
//   // 写数据
//   rpc Write (WriteRequest) returns (WriteResponse) {}
//   // 停止
//   rpc Stop (Request) returns (Response) {}
//
import (
	"context"
	"rulex/sidecar"
	"rulex/typex"

	"google.golang.org/grpc"
)

type SideCarDriver struct {
	state      typex.DriverState
	device     *typex.Device
	RuleEngine typex.RuleX
	client     sidecar.SidecarClient
}

func NewSideCarDriver(d *typex.Device,
	grpcConn *grpc.ClientConn, e typex.RuleX) typex.XExternalDriver {
	sideCarDriver := &SideCarDriver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		client:     sidecar.NewSidecarClient(grpcConn),
	}
	return sideCarDriver

}
func (sc *SideCarDriver) Test() error {
	_, err := sc.client.Status(context.Background(), &sidecar.Request{})
	if err != nil {
		return err
	}
	return nil
}

func (sc *SideCarDriver) Init() error {
	return nil
}

func (sc *SideCarDriver) Work() error {
	return nil
}

func (sc *SideCarDriver) State() typex.DriverState {
	_, err := sc.client.Status(context.Background(), &sidecar.Request{})
	if err != nil {
		return typex.DRIVER_STOP
	}
	return typex.DRIVER_RUNNING
}

/*
*
* 读取
*
 */
func (sc *SideCarDriver) Read(data []byte) (int, error) {
	response, err := sc.client.Read(context.Background(), &sidecar.ReadRequest{})
	if err != nil {
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
		return 0, err
	}
	return int(response.Code), nil
}

//---------------------------------------------------
func (sc *SideCarDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "SIDECAR-DRIVER",
		Type:        "SIDECAR",
		Description: "SIDECAR 通用GRPC协议驱动",
	}
}

func (sc *SideCarDriver) Stop() error {
	sc.client.Stop(context.Background(), &sidecar.Request{})
	return nil
}

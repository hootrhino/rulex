package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BeatTime/bacnet"
	"github.com/BeatTime/bacnet/btypes"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"time"
)

type bacnetIpCommonConfig struct {
	Ip        string `json:"ip,omitempty" title:"bacnet设备ip"`
	Port      int    `json:"port,omitempty" title:"bacnet端口，通常是47808"`
	LocalPort int    `json:"localPort" title:"本地监听端口，填0表示默认47808（有的模拟器必须本地监听47808才能正常交互）"`
	Interval  int    `json:"interval" title:"采集间隔，单位秒"`
}

type bacnetIpNodeConfig struct {
	IsMstp   int    `json:"isMstp,omitempty" title:"是否为mstp设备，若是则设备id和子网号必须填写"`
	DeviceId int    `json:"deviceId,omitempty" title:"设备id"`
	Subnet   int    `json:"subnet,omitempty" title:"子网号"`
	Tag      string `json:"tag" validate:"required" title:"数据Tag"`
	Type     int    `json:"type,omitempty" title:"object类型"`
	Id       int    `json:"id,omitempty" title:"object的id"`

	property btypes.PropertyData
}

type BacnetIpConfig struct {
	CommonConfig bacnetIpCommonConfig `json:"commonConfig"`
	NodeConfig   []bacnetIpNodeConfig `json:"nodeConfig"`
}

type GenericBacnetIpDevice struct {
	typex.XStatus
	status         typex.DeviceState
	RuleEngine     typex.RuleX
	bacnetIpConfig BacnetIpConfig
	// Bacnet
	bacnetClient bacnet.Client
	remoteDev    btypes.Device
}

func NewGenericBacnetIpDevice(e typex.RuleX) typex.XDevice {
	g := new(GenericBacnetIpDevice)
	g.RuleEngine = e
	g.bacnetIpConfig = BacnetIpConfig{
		CommonConfig: bacnetIpCommonConfig{},
		NodeConfig:   make([]bacnetIpNodeConfig, 0),
	}
	return g
}

func (dev *GenericBacnetIpDevice) Init(devId string, configMap map[string]interface{}) error {
	dev.PointId = devId
	err := utils.BindSourceConfig(configMap, &dev.bacnetIpConfig)
	if err != nil {
		return err
	}
	return nil
}

func (dev *GenericBacnetIpDevice) Start(cctx typex.CCTX) error {
	dev.CancelCTX = cctx.CancelCTX
	dev.Ctx = cctx.Ctx
	// 创建一个bacnetip的本地网络
	client, err := bacnet.NewClient(&bacnet.ClientBuilder{
		Ip:         "0.0.0.0",
		Port:       dev.bacnetIpConfig.CommonConfig.LocalPort,
		SubnetCIDR: 10, // 随便填一个，主要为了能够创建Client
	})

	if err != nil {
		return err
	}
	client.SetLogger(glogger.GLogger.Logger)

	// 将nodeConfig对应的配置信息
	for idx, v := range dev.bacnetIpConfig.NodeConfig {
		tmp := btypes.PropertyData{
			Object: btypes.Object{
				ID: btypes.ObjectID{
					Type:     btypes.ObjectType(v.Type),
					Instance: btypes.ObjectInstance(v.Id),
				},
				Properties: []btypes.Property{
					{
						Type:       btypes.PropPresentValue, // Present value
						ArrayIndex: btypes.ArrayAll,
					},
				},
			},
		}
		dev.bacnetIpConfig.NodeConfig[idx].property = tmp
	}

	mac := make([]byte, 6)
	fmt.Sscanf(dev.bacnetIpConfig.CommonConfig.Ip, "%d.%d.%d.%d", &mac[0], &mac[1], &mac[2], &mac[3])
	port := uint16(dev.bacnetIpConfig.CommonConfig.Port)
	mac[4] = byte(port >> 8)
	mac[5] = byte(port & 0x00FF)
	dev.remoteDev = btypes.Device{
		Addr: btypes.Address{
			MacLen: 6,
			Mac:    mac,
		},
	}
	dev.bacnetClient = client
	go dev.bacnetClient.ClientRun()

	go func(ctx context.Context) {
		interval := dev.bacnetIpConfig.CommonConfig.Interval
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					return
				}
			default:
				{
				}
			}

			read, err2 := dev.read()
			if err2 != nil {
				glogger.GLogger.Error(err2)
			} else {
				dev.RuleEngine.WorkDevice(dev.Details(), string(read))
			}
			<-ticker.C
		}
	}(dev.Ctx)

	dev.status = typex.DEV_UP
	return nil
}

func (dev *GenericBacnetIpDevice) OnRead(cmd []byte, data []byte) (int, error) {
	read, err := dev.read()
	if err != nil {
		return 0, err
	}
	len := copy(data, read)
	return len, nil
}

func (dev *GenericBacnetIpDevice) read() ([]byte, error) {
	retMap := map[string]string{}
	for _, v := range dev.bacnetIpConfig.NodeConfig {
		property, err := dev.bacnetClient.ReadProperty(dev.remoteDev, v.property)
		if err != nil {
			glogger.GLogger.Errorf("read failed. tag = %v, err=%v", v.Tag, err)
			continue
		}
		value := fmt.Sprintf("%v", property.Object.Properties[0].Data)
		retMap[v.Tag] = value
	}
	bytes, _ := json.Marshal(retMap)
	glogger.GLogger.Debugf("%v", retMap)
	return bytes, nil
}

func (dev *GenericBacnetIpDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	//TODO implement me
	return 0, errors.New("not Support")
}

func (dev *GenericBacnetIpDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return nil, errors.New("not Support")
}

func (dev *GenericBacnetIpDevice) Status() typex.DeviceState {
	return dev.status
}

func (dev *GenericBacnetIpDevice) Stop() {
	dev.CancelCTX()
	if dev.bacnetClient != nil {
		dev.bacnetClient.Close()
	}
}

func (dev *GenericBacnetIpDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

func (dev *GenericBacnetIpDevice) Details() *typex.Device {
	return dev.RuleEngine.GetDevice(dev.PointId)
}

func (dev *GenericBacnetIpDevice) SetState(state typex.DeviceState) {
	dev.status = state
}

func (dev *GenericBacnetIpDevice) Driver() typex.XExternalDriver {
	return nil
}

func (dev *GenericBacnetIpDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

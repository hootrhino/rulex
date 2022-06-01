package device

import (
	"context"
	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
	"github.com/ngaut/log"
)

type _modBusConfig struct {
	Timeout   *int       `json:"timeout" validate:"required" title:"连接超时" info:""`
	SlaverIds []byte     `json:"slaverIds" validate:"required" title:"TCP端口" info:""`
	Config    _rtuConfig `json:"config" validate:"required" title:"工作模式" info:""`
}

type _rtuConfig struct {
	Uart     string       `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int          `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int          `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   typex.Parity `json:"parity" validate:"required" title:"校验位" info:"串口通信校验位"`
	StopBits int          `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
type ts200_v_0_2_sensor struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	slaverIds  []byte
}

var _ typex.XDevice = new(ts200_v_0_2_sensor)

func NewTS200Sensor(deviceId string, e typex.RuleX) typex.XDevice {
	tss := new(ts200_v_0_2_sensor)
	tss.PointId = deviceId
	tss.RuleEngine = e
	return tss
}

//  初始化
func (tss *ts200_v_0_2_sensor) Init(devId string, config map[string]interface{}) error {

	return nil
}

// 启动
func (tss *ts200_v_0_2_sensor) Start(cctx typex.CCTX) error {
	tss.Ctx = cctx.Ctx
	tss.CancelCTX = cctx.CancelCTX
	config := tss.RuleEngine.GetDevice(tss.PointId).Config
	var mainConfig _modBusConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	var rtuConfig _rtuConfig
	if errs := mapstructure.Decode(mainConfig.Config, &rtuConfig); errs != nil {
		log.Error(errs)
		return errs
	}

	// 串口配置固定写法
	handler := modbus.NewRTUClientHandler(rtuConfig.Uart)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = time.Duration(*mainConfig.Timeout) * time.Second
	// handler.Logger = golog.New(os.Stdout, "485THerSource: ", log.LstdFlags)
	if err := handler.Connect(); err != nil {
		return err
	}
	client := modbus.NewClient(handler)
	tss.driver = driver.NewTS200_v_0_2_Driver(tss.Details(), tss.RuleEngine, client)
	tss.slaverIds = append(tss.slaverIds, mainConfig.SlaverIds...)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	tss.status = typex.DEV_RUNNING
	for _, slaverId := range tss.slaverIds {
		go func(ctx context.Context, slaverId byte, rtuDriver typex.XExternalDriver, handler *modbus.RTUClientHandler) {
			ticker := time.NewTicker(time.Duration(5) * time.Second)
			defer ticker.Stop()
			buffer := make([]byte, 128) //128字节数据
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						tss.status = typex.DEV_STOP
						return
					}
				default:
					{
					}
				}
				handler.SlaveId = slaverId // 配置ID
				n, err := rtuDriver.Read(buffer)
				if err != nil {
					log.Error(err)
				} else {
					tss.RuleEngine.WorkDevice(tss.RuleEngine.GetDevice(tss.PointId), string(buffer[:n]))
				}
			}

		}(tss.Ctx, slaverId, tss.driver, handler)
	}
	return nil
}

// 从设备里面读数据出来
func (tss *ts200_v_0_2_sensor) OnRead(data []byte) (int, error) {

	n, err := tss.driver.Read(data)
	if err != nil {
		log.Error(err)
		tss.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (tss *ts200_v_0_2_sensor) OnWrite(_ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (tss *ts200_v_0_2_sensor) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (tss *ts200_v_0_2_sensor) Stop() {
	tss.CancelCTX()
}

// 设备属性，是一系列属性描述
func (tss *ts200_v_0_2_sensor) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (tss *ts200_v_0_2_sensor) Details() *typex.Device {
	return tss.RuleEngine.GetDevice(tss.PointId)
}

// 状态
func (tss *ts200_v_0_2_sensor) SetState(status typex.DeviceState) {
	tss.status = status

}

// 驱动
func (tss *ts200_v_0_2_sensor) Driver() typex.XExternalDriver {
	return tss.driver
}

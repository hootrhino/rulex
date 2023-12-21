package device

import (
	"net/http"
	"sync"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type __HttpCommonConfig struct {
	Timeout     *int   `json:"timeout" validate:"required"`
	AutoRequest *bool  `json:"autoRequest" validate:"required"`
	Frequency   *int64 `json:"frequency" validate:"required"`
}
type __HttpMainConfig struct {
	CommonConfig __HttpCommonConfig `json:"commonConfig" validate:"required"`
	HttpConfig   common.HTTPConfig  `json:"httpConfig" validate:"required"`
}

type GenericHttpDevice struct {
	typex.XStatus
	client     http.Client
	status     typex.DeviceState
	RuleEngine typex.RuleX
	mainConfig __HttpMainConfig
	locker     sync.Locker
}

/*
*
* 通用串口透传
*
 */
func NewGenericHttpDevice(e typex.RuleX) typex.XDevice {
	hd := new(GenericHttpDevice)
	hd.locker = &sync.Mutex{}
	hd.client = *http.DefaultClient
	hd.mainConfig = __HttpMainConfig{
		CommonConfig: __HttpCommonConfig{
			AutoRequest: new(bool),
		},
	}
	hd.RuleEngine = e
	return hd
}

//  初始化
func (hd *GenericHttpDevice) Init(devId string, configMap map[string]interface{}) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	return nil
}

// 启动
func (hd *GenericHttpDevice) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX

	if *hd.mainConfig.CommonConfig.AutoRequest {
		ticker := time.NewTicker(
			time.Duration(*hd.mainConfig.CommonConfig.Frequency) * time.Millisecond)
		go func() {
			for {
				select {
				case <-hd.Ctx.Done():
					{
						ticker.Stop()
						return
					}
				default:
					{
					}
				}
				body := utils.Get(hd.client, hd.mainConfig.HttpConfig.Url)
				if body != "" {
					hd.RuleEngine.WorkDevice(hd.Details(), body)
				}
				<-ticker.C
			}
		}()

	}
	hd.status = typex.DEV_UP
	return nil
}

func (hd *GenericHttpDevice) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
func (hd *GenericHttpDevice) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (hd *GenericHttpDevice) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (hd *GenericHttpDevice) Stop() {
	hd.status = typex.DEV_STOP
	if hd.CancelCTX != nil {
		hd.CancelCTX()
	}
}

// 设备属性，是一系列属性描述
func (hd *GenericHttpDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (hd *GenericHttpDevice) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

// 状态
func (hd *GenericHttpDevice) SetState(status typex.DeviceState) {
	hd.status = status

}

// 驱动
func (hd *GenericHttpDevice) Driver() typex.XExternalDriver {
	return nil
}

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------

func (hd *GenericHttpDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (hd *GenericHttpDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

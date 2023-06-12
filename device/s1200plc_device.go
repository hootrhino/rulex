package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/driver"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/robinson/gos7"
)

type s1200plc struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	mainConfig common.S1200Config
	client     gos7.Client
	block      []common.S1200Block // PLC 的DB块
	lock       sync.Mutex
}

/*
*
* 西门子 S1200 系列 PLC
*
 */
func NewS1200plc(e typex.RuleX) typex.XDevice {
	s1200 := new(s1200plc)
	s1200.RuleEngine = e
	s1200.lock = sync.Mutex{}
	s1200.mainConfig = common.S1200Config{}
	return s1200
}

// 初始化
func (s1200 *s1200plc) Init(devId string, configMap map[string]interface{}) error {
	s1200.PointId = devId
	if err := utils.BindSourceConfig(configMap, &s1200.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	// 检查Tag有没有重复
	tags := []string{}
	for _, block := range s1200.mainConfig.Blocks {
		tags = append(tags, block.Tag)
	}
	if utils.IsListDuplicated(tags) {
		return errors.New("tag duplicated")
	}
	return nil
}

// 启动
func (s1200 *s1200plc) Start(cctx typex.CCTX) error {
	s1200.Ctx = cctx.Ctx
	s1200.CancelCTX = cctx.CancelCTX
	//
	handler := gos7.NewTCPClientHandler(
		// 127.0.0.1:8080
		fmt.Sprintf("%s:%d", s1200.mainConfig.Host, *s1200.mainConfig.Port),
		*s1200.mainConfig.Rack,
		*s1200.mainConfig.Slot)
	handler.Timeout = 5 * time.Second
	if err := handler.Connect(); err != nil {
		return err
	}
	handler.Timeout = time.Duration(*s1200.mainConfig.Timeout) * time.Second
	handler.IdleTimeout = time.Duration(*s1200.mainConfig.IdleTimeout) * time.Second
	s1200.client = gos7.NewClient(handler)
	s1200.driver = driver.NewS1200Driver(s1200.Details(), s1200.RuleEngine, s1200.client, s1200.block)
	if !s1200.mainConfig.AutoRequest {
		s1200.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(s1200.mainConfig.Frequency) * time.Millisecond)
		// 数据缓冲区,最大4KB
		dataBuffer := make([]byte, common.T_4KB)
		s1200.driver.Read([]byte{}, dataBuffer) //清理缓存
		for {
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					if s1200.driver != nil {
						s1200.driver.Stop()
					}
					return
				}
			default:
				{
					// Do nothing
				}
			}
			if s1200.driver == nil {
				return
			}
			s1200.lock.Lock()
			n, err := s1200.driver.Read([]byte{}, dataBuffer)
			s1200.lock.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
				return
			}
			ok, err := s1200.RuleEngine.WorkDevice(
				s1200.RuleEngine.GetDevice(s1200.PointId),
				string(dataBuffer[:n]),
			)
			if !ok {
				glogger.GLogger.Error(err)
			}
			<-ticker.C
		}

	}(cctx.Ctx)
	return nil
}

// 从设备里面读数据出来
func (s1200 *s1200plc) OnRead(cmd []byte, data []byte) (int, error) {
	return s1200.driver.Read(cmd, data)
}

// 把数据写入设备
//
// db.Address:int, db.Start:int, db.Size:int, rData[]
// [
//
//	{
//	    "tag":"V",
//	    "address":1,
//	    "start":1,
//	    "size":1,
//	    "value":"AAECAwQ="
//	}
//
// ]
func (s1200 *s1200plc) OnWrite(cmd []byte, data []byte) (int, error) {
	blocks := []common.S1200BlockValue{}
	if err := json.Unmarshal(data, &blocks); err != nil {
		return 0, err
	}
	return s1200.driver.Write(cmd, data)
}

// 设备当前状态
func (s1200 *s1200plc) Status() typex.DeviceState {
	if s1200.driver.State() == typex.DRIVER_UP {
		return typex.DEV_UP
	}
	return typex.DEV_DOWN

}

// 停止设备
func (s1200 *s1200plc) Stop() {
	s1200.status = typex.DEV_DOWN
	s1200.CancelCTX()

}

// 设备属性，是一系列属性描述
func (s1200 *s1200plc) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (s1200 *s1200plc) Details() *typex.Device {
	return s1200.RuleEngine.GetDevice(s1200.PointId)
}

// 状态
func (s1200 *s1200plc) SetState(status typex.DeviceState) {
	s1200.status = status
}

// 驱动
func (s1200 *s1200plc) Driver() typex.XExternalDriver {
	return s1200.driver
}

func (s1200 *s1200plc) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (s1200 *s1200plc) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

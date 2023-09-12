// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package device

import (
	"context"
	"errors"
	"fmt"
	golog "log"

	"sync"
	"time"

	archsupport "github.com/hootrhino/rulex/bspsupport"
	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/driver"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	modbus "github.com/wwhai/gomodbus"
)

// 这是个通用Modbus采集器, 主要用来在通用场景下采集数据，因此需要配合规则引擎来使用
//
// Modbus 采集到的数据如下, LUA 脚本可做解析, 示例脚本可参照 generic_modbus_parse.lua
//
//	{
//	    "d1":{
//	        "tag":"d1",
//	        "function":3,
//	        "slaverId":1,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    },
//	    "d2":{
//	        "tag":"d2",
//	        "function":3,
//	        "slaverId":2,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    }
//	}
type _GMODCommonConfig struct {
	Mode        string `json:"mode" title:"工作模式" info:"RTU/TCP"`
	Timeout     int    `json:"timeout" validate:"required" title:"连接超时"`
	AutoRequest bool   `json:"autoRequest" title:"启动轮询"`
	Frequency   int64  `json:"frequency" validate:"required" title:"采集频率"`
}
type _GMODHostConfig struct {
	Host string `json:"host" title:"服务地址"`
	Port int    `json:"port" title:"服务端口"`
}

type _GMODConfig struct {
	CommonConfig _GMODCommonConfig       `json:"commonConfig" validate:"required"`
	RtuConfig    common.CommonUartConfig `json:"rtuConfig"`
	TcpConfig    _GMODHostConfig         `json:"tcpConfig"`
	Registers    []common.RegisterRW     `json:"registers" validate:"required" title:"寄存器配置"`
}
type generic_modbus_device struct {
	typex.XStatus ``
	status        typex.DeviceState
	RuleEngine    typex.RuleX
	driver        typex.XExternalDriver
	rtuHandler    *modbus.RTUClientHandler
	tcpHandler    *modbus.TCPClientHandler
	mainConfig    _GMODConfig
	locker        sync.Locker
	retryTimes    int
}

/*
*
* 温湿度传感器
*
 */
func NewGenericModbusDevice(e typex.RuleX) typex.XDevice {
	mdev := new(generic_modbus_device)
	mdev.RuleEngine = e
	mdev.locker = &sync.Mutex{}
	mdev.mainConfig = _GMODConfig{
		CommonConfig: _GMODCommonConfig{},
		TcpConfig:    _GMODHostConfig{Host: "127.0.0.1", Port: 502},
		RtuConfig: common.CommonUartConfig{
			Timeout:  3000,
			Uart:     "/tty/s1",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
	}
	mdev.Busy = false
	mdev.status = typex.DEV_DOWN
	return mdev
}

//  初始化
func (mdev *generic_modbus_device) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	// 超时大于30秒无意义
	if mdev.mainConfig.CommonConfig.Timeout > 30000 {
		return errors.New("'timeout' must less than 30 second")
	}
	// 频率不能太快
	if mdev.mainConfig.CommonConfig.Frequency < 50 {
		return errors.New("'frequency' must grate than 50 millisecond")

	}
	// 检查Tag有没有重复
	tags := []string{}
	for _, register := range mdev.mainConfig.Registers {
		tags = append(tags, register.Tag)
	}
	if utils.IsListDuplicated(tags) {
		return errors.New("tag duplicated")
	}
	if !utils.SContains([]string{"RTU", "TCP"}, mdev.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'RTU'")
	}
	// 做端口管理
	Port := archsupport.GetUart(mdev.mainConfig.RtuConfig.Uart)
	if Port.Busy {
		return errors.New(Port.BusyingInfo())
	}
	return nil
}

// 启动
func (mdev *generic_modbus_device) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX

	if mdev.mainConfig.CommonConfig.Mode == "RTU" {
		mdev.rtuHandler = modbus.NewRTUClientHandler(mdev.mainConfig.RtuConfig.Uart)
		mdev.rtuHandler.BaudRate = mdev.mainConfig.RtuConfig.BaudRate
		mdev.rtuHandler.DataBits = mdev.mainConfig.RtuConfig.DataBits
		mdev.rtuHandler.Parity = mdev.mainConfig.RtuConfig.Parity
		mdev.rtuHandler.StopBits = mdev.mainConfig.RtuConfig.StopBits
		// timeout 最大不能超过20, 不然无意义
		mdev.rtuHandler.Timeout = time.Duration(mdev.mainConfig.RtuConfig.Timeout) * time.Millisecond
		if core.GlobalConfig.AppDebugMode {
			mdev.rtuHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus RTU Mode: "+mdev.PointId+", ", golog.LstdFlags)
		}

		if err := mdev.rtuHandler.Connect(); err != nil {
			return err
		}
		client := modbus.NewClient(mdev.rtuHandler)
		mdev.driver = driver.NewModBusRtuDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.rtuHandler,
			client, mdev.mainConfig.CommonConfig.Frequency)
	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		mdev.tcpHandler = modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mdev.mainConfig.TcpConfig.Host, mdev.mainConfig.TcpConfig.Port),
		)
		if core.GlobalConfig.AppDebugMode {
			mdev.tcpHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus TCP Mode: "+mdev.PointId+", ", golog.LstdFlags)
		}

		if err := mdev.tcpHandler.Connect(); err != nil {
			return err
		}
		client := modbus.NewClient(mdev.tcpHandler)
		mdev.driver = driver.NewModBusTCPDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.tcpHandler, client,
			mdev.mainConfig.CommonConfig.Frequency)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	if !mdev.mainConfig.CommonConfig.AutoRequest {
		mdev.status = typex.DEV_UP
		archsupport.UARTBusy(mdev.mainConfig.RtuConfig.Uart)
		return nil
	}
	mdev.retryTimes = 0
	go func(ctx context.Context, Driver typex.XExternalDriver) {
		buffer := make([]byte, common.T_64KB)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			n, err := Driver.Read([]byte{}, buffer)
			if err != nil {
				glogger.GLogger.Error(err)
				mdev.retryTimes++
			} else {
				mdev.RuleEngine.WorkDevice(mdev.Details(), string(buffer[:n]))
			}
		}

	}(mdev.Ctx, mdev.driver)
	mdev.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来
func (mdev *generic_modbus_device) OnRead(cmd []byte, data []byte) (int, error) {

	n, err := mdev.driver.Read(cmd, data)
	if err != nil {
		glogger.GLogger.Error(err)
		mdev.retryTimes++
	}
	return n, err
}

// 把数据写入设备
func (mdev *generic_modbus_device) OnWrite(cmd []byte, data []byte) (int, error) {
	if mdev.Busy {
		return 0, fmt.Errorf("device busing now")
	}
	return mdev.driver.Write(cmd, data)
}

// 设备当前状态
func (mdev *generic_modbus_device) Status() typex.DeviceState {
	// 容错5次
	if mdev.retryTimes > 0 {
		return typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *generic_modbus_device) Stop() {
	mdev.CancelCTX()
	mdev.status = typex.DEV_DOWN
	if mdev.driver != nil {
		mdev.driver.Stop()
	}
	archsupport.UARTFree(mdev.mainConfig.RtuConfig.Uart)
}

// 设备属性，是一系列属性描述
func (mdev *generic_modbus_device) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *generic_modbus_device) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *generic_modbus_device) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *generic_modbus_device) Driver() typex.XExternalDriver {
	return mdev.driver
}
func (mdev *generic_modbus_device) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (mdev *generic_modbus_device) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	serial "github.com/wwhai/tarmserial"
)

// 读出来的字节缓冲默认大小
const __DEFAULT_BUFFER_SIZE = 100

// 传输形式：
// `rawtcp`, `rawudp`, `rawserial`
const rawtcp string = "rawtcp"
const rawudp string = "rawudp"
const rawserial string = "rawserial"

type _CPDCommonConfig struct {
	Transport string `json:"transport" validate:"required"` // 传输协议
	RetryTime int    `json:"retryTime" validate:"required"` // 几次以后重启,0 表示不重启
}

/*
*
* 自定义协议
*
 */
type _CustomProtocolConfig struct {
	CommonConfig _CPDCommonConfig        `json:"commonConfig" validate:"required"`
	UartConfig   common.CommonUartConfig `json:"uartConfig" validate:"required"`
	HostConfig   common.HostConfig       `json:"hostConfig" validate:"required"`
}
type CustomProtocolDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	serialPort *serial.Port // 串口
	tcpcon     net.Conn     // TCP
	mainConfig _CustomProtocolConfig
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启
}

func NewCustomProtocolDevice(e typex.RuleX) typex.XDevice {
	mdev := new(CustomProtocolDevice)
	mdev.RuleEngine = e
	mdev.mainConfig = _CustomProtocolConfig{
		CommonConfig: _CPDCommonConfig{},
		UartConfig:   common.CommonUartConfig{},
		HostConfig:   common.HostConfig{Timeout: 50},
	}
	return mdev

}

// 初始化
func (mdev *CustomProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"N", "E", "O"}, mdev.mainConfig.UartConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	if !utils.SContains([]string{`rawtcp`, `rawudp`, `rawserial`},
		mdev.mainConfig.CommonConfig.Transport) {
		return errors.New("option only one of 'rawtcp','rawudp','rawserial'")
	}
	return nil
}

// 启动
func (mdev *CustomProtocolDevice) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	mdev.errorCount = 0
	mdev.status = typex.DEV_DOWN

	// 现阶段暂时只支持RS485串口, 以后有需求再支持TCP、UDP
	if mdev.mainConfig.CommonConfig.Transport == "rawserial" {
		config := serial.Config{
			Name:        mdev.mainConfig.UartConfig.Uart,
			Baud:        mdev.mainConfig.UartConfig.BaudRate,
			Size:        byte(mdev.mainConfig.UartConfig.DataBits),
			Parity:      serial.Parity(mdev.mainConfig.UartConfig.Parity[0]),
			StopBits:    serial.StopBits(mdev.mainConfig.UartConfig.StopBits),
			ReadTimeout: time.Duration(mdev.mainConfig.UartConfig.Timeout) * time.Millisecond,
		}
		serialPort, err := serial.OpenPort(&config)
		if err != nil {
			glogger.GLogger.Error("serialPort start failed:", err)
			return err
		}
		mdev.serialPort = serialPort
		mdev.status = typex.DEV_UP
		return nil
	}

	// rawtcp
	if mdev.mainConfig.CommonConfig.Transport == "rawtcp" {
		tcpcon, err := net.Dial("tcp",
			fmt.Sprintf("%s:%d", mdev.mainConfig.HostConfig.Host,
				mdev.mainConfig.HostConfig.Port))
		if err != nil {
			glogger.GLogger.Error("tcp connection start failed:", err)
			return err
		}
		mdev.tcpcon = tcpcon
		mdev.status = typex.DEV_UP
		return nil
	}
	return fmt.Errorf("unsupported transport:%s", mdev.mainConfig.CommonConfig.Transport)
}

/*
*
* 数据读出来，对数据结构有要求, 其中Key必须是个数字或者数字字符串, 例如 1 or "1"
*
 */
func (mdev *CustomProtocolDevice) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, errors.New("unknown read command:" + string(cmd))

}

/*
*
* 写进来的数据格式 参考@Protocol
*
 */

// 把数据写入设备
func (mdev *CustomProtocolDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	return 0, errors.New("unknown write command:" + string(cmd))
}

/*
*
* 外部指令交互, 常用来实现自定义协议等
*
 */
func (mdev *CustomProtocolDevice) OnCtrl(cmd []byte, _ []byte) ([]byte, error) {
	glogger.GLogger.Debug("Time slice SliceRequest:", string(cmd))
	return mdev.ctrl(cmd)
}

// 设备当前状态
func (mdev *CustomProtocolDevice) Status() typex.DeviceState {
	if mdev.errorCount >= mdev.mainConfig.CommonConfig.RetryTime {
		mdev.CancelCTX()
		mdev.status = typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *CustomProtocolDevice) Stop() {
	mdev.CancelCTX()
	mdev.status = typex.DEV_DOWN
	if mdev.mainConfig.CommonConfig.Transport == rawtcp {
		if mdev.tcpcon != nil {
			mdev.tcpcon.Close()
		}
	}
	if mdev.mainConfig.CommonConfig.Transport == rawserial {
		if mdev.serialPort != nil {
			mdev.serialPort.Close()
		}
	}
}

// 设备属性，是一系列属性描述
func (mdev *CustomProtocolDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *CustomProtocolDevice) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *CustomProtocolDevice) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *CustomProtocolDevice) Driver() typex.XExternalDriver {
	return nil
}

/*
*
* 设备服务调用
*
 */
func (mdev *CustomProtocolDevice) OnDCACall(_ string, Command string,
	Args interface{}) typex.DCAResult {

	return typex.DCAResult{}
}

// --------------------------------------------------------------------------------------------------
// 内部函数
// --------------------------------------------------------------------------------------------------
func (mdev *CustomProtocolDevice) ctrl(args []byte) ([]byte, error) {
	hexs, err1 := hex.DecodeString(string(args))
	if err1 != nil {
		glogger.GLogger.Error(err1)
		return nil, err1
	}
	glogger.GLogger.Debug("Custom Protocol Device Request:", hexs)
	result := [__DEFAULT_BUFFER_SIZE]byte{}
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(mdev.mainConfig.UartConfig.Timeout)*time.Millisecond)
	var count int = 0
	var errSliceRequest error = nil
	if mdev.mainConfig.CommonConfig.Transport == rawserial {
		count, errSliceRequest = utils.SliceRequest(ctx, mdev.serialPort,
			hexs, result[:], false,
			time.Duration(30)*time.Millisecond /*30ms wait*/)
	}
	if mdev.mainConfig.CommonConfig.Transport == rawtcp {
		mdev.tcpcon.SetReadDeadline(
			time.Now().Add((time.Duration(mdev.mainConfig.HostConfig.Timeout) * time.Millisecond)),
		)
		count, errSliceRequest = utils.SliceRequest(ctx, mdev.tcpcon,
			hexs, result[:], false,
			time.Duration(30)*time.Millisecond /*30ms wait*/)
		mdev.tcpcon.SetReadDeadline(time.Time{})
	}

	cancel()
	if errSliceRequest != nil {
		glogger.GLogger.Error("Custom Protocol Device Request error: ", errSliceRequest)
		mdev.errorCount++
		return nil, errSliceRequest
	}
	dataMap := map[string]string{}
	dataMap["in"] = string(args)
	dataMap["out"] = hex.EncodeToString(result[:count])
	bytes, _ := json.Marshal(dataMap)
	return []byte(bytes), nil
}

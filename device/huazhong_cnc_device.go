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
	"github.com/hootrhino/rulex/component/iotschema"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

// HNC8 NC-Link 接口参口手册 版本 v1.0 2019-02-13
type hnc8_cnc_config struct {
	CNCSerialNumber string `json:"serialNumber" validate:"required"` // CNC 序列号
	Host            string `json:"host" validate:"required"`         // IP:Port
	ApiVersion      int    `json:"apiVersion" validate:"required"`   // API 版本,2 | 3
}

/*
*
* 凯帝恩CNC
*
 */
type HNC8_DEVICE_Device struct {
	typex.XStatus
	mainConfig hnc8_cnc_config
	status     typex.DeviceState
}

func NewHNC8_DEVICE_Device(e typex.RuleX) typex.XDevice {
	hd := new(HNC8_DEVICE_Device)
	hd.RuleEngine = e
	return hd
}

//  初始化
func (hd *HNC8_DEVICE_Device) Init(devId string, configMap map[string]interface{}) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	return nil
}

// 启动
func (hd *HNC8_DEVICE_Device) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX

	hd.status = typex.DEV_UP
	return nil
}

func (hd *HNC8_DEVICE_Device) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
func (hd *HNC8_DEVICE_Device) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (hd *HNC8_DEVICE_Device) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (hd *HNC8_DEVICE_Device) Stop() {
	hd.status = typex.DEV_STOP
	hd.CancelCTX()
}

// 设备属性，是一系列属性描述
func (hd *HNC8_DEVICE_Device) Property() []iotschema.IoTSchema {
	return []iotschema.IoTSchema{}
}

// 真实设备
func (hd *HNC8_DEVICE_Device) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

// 状态
func (hd *HNC8_DEVICE_Device) SetState(status typex.DeviceState) {
	hd.status = status

}

// 驱动
func (hd *HNC8_DEVICE_Device) Driver() typex.XExternalDriver {
	return nil
}

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------

func (hd *HNC8_DEVICE_Device) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (hd *HNC8_DEVICE_Device) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

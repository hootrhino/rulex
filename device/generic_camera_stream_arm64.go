//go:build arm64
// +build arm64

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

	"github.com/hootrhino/rulex/typex"
)

type videoCamera struct {
	typex.XStatus
	status typex.DeviceState
}

/*
*
* ARM32不支持
*
 */
func NewVideoCamera(e typex.RuleX) typex.XDevice {
	hd := new(videoCamera)
	hd.RuleEngine = e
	return hd
}

func (hd *videoCamera) Init(devId string, configMap map[string]interface{}) error {
	hd.PointId = devId
	return nil
}

func (hd *videoCamera) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX

	hd.status = typex.DEV_UP
	return nil
}

func (hd *videoCamera) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
func (hd *videoCamera) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (hd *videoCamera) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (hd *videoCamera) Stop() {
	hd.status = typex.DEV_DOWN
	hd.CancelCTX()
}

// 设备属性，是一系列属性描述
func (hd *videoCamera) Property() []iotschema.IoTSchema {
	return []iotschema.IoTSchema{}
}

// 真实设备
func (hd *videoCamera) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

// 状态
func (hd *videoCamera) SetState(status typex.DeviceState) {
	hd.status = status

}

// 驱动
func (hd *videoCamera) Driver() typex.XExternalDriver {
	return nil
}

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------
func (hd *videoCamera) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (hd *videoCamera) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

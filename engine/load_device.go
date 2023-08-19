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

package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

//--------------------------------------------------------------------------------------------------
// Abstract device
//--------------------------------------------------------------------------------------------------

// 获取设备
func (e *RuleEngine) GetDevice(id string) *typex.Device {
	v, ok := e.Devices.Load(id)
	if ok {
		return v.(*typex.Device)
	} else {
		return nil
	}

}

// 保存设备
func (e *RuleEngine) SaveDevice(dev *typex.Device) {
	e.Devices.Store(dev.UUID, dev)
}

// 获取所有外挂设备
func (e *RuleEngine) AllDevices() *sync.Map {
	return e.Devices

}

// 删除设备
func (e *RuleEngine) RemoveDevice(uuid string) {
	if dev := e.GetDevice(uuid); dev != nil {
		if dev.Device != nil {
			glogger.GLogger.Infof("Device [%v] ready to stop", uuid)
			if dev.Device.Status() == typex.DEV_UP {
				dev.Device.Stop()
			}
			glogger.GLogger.Infof("Device [%v] has been stopped", uuid)
			e.Devices.Delete(uuid)
			glogger.GLogger.Infof("Device [%v] has been deleted", uuid)
		}

	}
}

/*
*
* 加载设备
*
 */
func (e *RuleEngine) LoadDeviceWithCtx(deviceInfo *typex.Device,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	if config := e.DeviceTypeManager.Find(deviceInfo.Type); config != nil {
		return e.loadDevices(config.NewDevice(e), deviceInfo, ctx, cancelCTX)
	}
	return fmt.Errorf("unsupported Device type:%s", deviceInfo.Type)

}

/*
*
* 启动一个和RULEX直连的外部设备
*
 */
func (e *RuleEngine) loadDevices(abstractDevice typex.XDevice, deviceInfo *typex.Device,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	// Bind
	deviceInfo.Device = abstractDevice
	e.SaveDevice(deviceInfo)
	// Load config
	// 要从数据库里面查Config
	config := e.GetDevice(deviceInfo.UUID).Config
	if config == nil {
		e.RemoveDevice(deviceInfo.UUID)
		err := fmt.Errorf("device [%v] config is nil", deviceInfo.Name)
		return err
	}
	if err := abstractDevice.Init(deviceInfo.UUID, config); err != nil {
		e.RemoveDevice(deviceInfo.UUID)
		return err
	}
	startDevice(abstractDevice, e, ctx, cancelCTX)
	glogger.GLogger.Infof("device [%v, %v] load successfully", deviceInfo.Name, deviceInfo.UUID)
	return nil
}

func startDevice(abstractDevice typex.XDevice, e *RuleEngine,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	if err := abstractDevice.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("abstractDevice start error:", err)
		return err
	}
	// LoadNewestDevice
	// 2023-06-14新增： 重启成功后数据会丢失,还得加载最新的Rule到设备中
	device := abstractDevice.Details()
	if device != nil {
		// bind 最新的规则 要从数据库拿刚更新的
		for _, rule := range device.BindRules {
			glogger.GLogger.Debugf("Load rule:%s", rule.Name)
			RuleInstance := typex.NewLuaRule(e,
				rule.UUID,
				rule.Name,
				rule.Description,
				rule.FromSource,
				rule.FromDevice,
				rule.Success,
				rule.Actions,
				rule.Failed)
			if err1 := e.LoadRule(RuleInstance); err1 != nil {
				return err1
			}
		}
	}
	return nil
}

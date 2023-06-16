package engine

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

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
			dev.Device.Stop()
			glogger.GLogger.Infof("Device [%v] has been stopped", uuid)
			e.Devices.Delete(uuid)
			glogger.GLogger.Infof("Device [%v] has been deleted", uuid)
		}

	}
}

/*
* 从0.4.0开始, 可支持用户加载设备
* 加载用户设备， 第一个参数为Worker, 第二个参数为设备原始信息，实际上就是元数据
*
 */
func (e *RuleEngine) LoadUserDevice(abstractDevice typex.XDevice, deviceInfo *typex.Device) error {
	return loadDevices(abstractDevice, deviceInfo, e)
}

// 加载内置设备
func (e *RuleEngine) LoadBuiltinDevice(deviceInfo *typex.Device) error {
	return e.LoadDevice(deviceInfo)
}

/*
*
* 加载设备
*
 */
func (e *RuleEngine) LoadDevice(deviceInfo *typex.Device) error {
	if config := e.DeviceTypeManager.Find(deviceInfo.Type); config != nil {
		return loadDevices(config.Device, deviceInfo, e)
	}
	return fmt.Errorf("unsupported Device type:%s", deviceInfo.Type)

}

/*
*
* 启动一个和RULEX直连的外部设备
*
 */
func loadDevices(abstractDevice typex.XDevice, deviceInfo *typex.Device, e *RuleEngine) error {
	// Bind
	deviceInfo.Device = abstractDevice
	e.SaveDevice(deviceInfo)
	// Load config
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

	// start
	// if err := startDevice(abstractDevice, e); err != nil {
	// 	glogger.GLogger.Error(err)
	// 	e.RemoveDevice(deviceInfo.UUID)
	// 	return err
	// }
	startDevice(abstractDevice, e)
	go func(ctx context.Context) {
		for {
			ticker := time.NewTicker(time.Duration(time.Second * 5))
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
			<-ticker.C
			if abstractDevice.Details() == nil {
				return
			}
			tryIfRestartDevice(abstractDevice, e, deviceInfo.UUID)

		}

	}(typex.GCTX)
	glogger.GLogger.Infof("device [%v, %v] load successfully", deviceInfo.Name, deviceInfo.UUID)
	return nil
}

func tryIfRestartDevice(abstractDevice typex.XDevice, e *RuleEngine, devId string) {
	Status := abstractDevice.Status()
	if Status == typex.DEV_STOP {
		return
	}
	if Status == typex.DEV_DOWN {
		abstractDevice.Details().State = typex.DEV_DOWN
		glogger.GLogger.Warnf("Device [%v, %v] down. try to restart.",
			abstractDevice.Details().UUID, abstractDevice.Details().Name)
		abstractDevice.Stop()
		runtime.Gosched()
		runtime.GC()
		startDevice(abstractDevice, e)
	} else {
		abstractDevice.Details().State = typex.DEV_UP
	}

}

func startDevice(abstractDevice typex.XDevice, e *RuleEngine) error {
	ctx, cancelCTX := typex.NewCCTX()
	if err := abstractDevice.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("abstractDevice start error:", err)
		return err
	}
	// LoadNewestDevice
	// 2023-06-14新增： 重启成功后数据会丢失,还得加载最新的Rule到设备中
	device := abstractDevice.Details()
	if device != nil {
		for _, rule := range device.BindRules {
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

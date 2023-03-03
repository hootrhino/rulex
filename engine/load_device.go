package engine

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
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
		glogger.GLogger.Infof("Device [%v] ready to stop", uuid)
		if dev.Device != nil {
			dev.Device.Stop()
		}
		glogger.GLogger.Infof("Device [%v] has been stopped", uuid)
		e.Devices.Delete(uuid)
		dev = nil
		glogger.GLogger.Infof("Device [%v] has been deleted", uuid)
	}
}

/*
* 从0.4.0开始, 可支持用户加载设备
* 加载用户设备， 第一个参数为Worker, 第二个参数为设备原始信息，实际上就是元数据
*
 */
func (e *RuleEngine) LoadUserDevice(abstractDevice typex.XDevice, deviceInfo *typex.Device) error {
	return startDevices(abstractDevice, deviceInfo, e)
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
		return startDevices(config.Device, deviceInfo, e)
	}
	return fmt.Errorf("unsupported Device type:%s", deviceInfo.Type)

}

/*
*
* 启动一个和RULEX直连的外部设备
*
 */
func startDevices(abstractDevice typex.XDevice, deviceInfo *typex.Device, e *RuleEngine) error {
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
		err := fmt.Errorf("device [%v] Init error", deviceInfo.UUID)
		return err
	}
	// Bind
	deviceInfo.Device = abstractDevice
	// start
	if err := startDevice(abstractDevice, e); err != nil {
		glogger.GLogger.Error(err)
		e.RemoveDevice(deviceInfo.UUID)
		return err
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(time.Second * 5))
		// 5 seconds
	TICKER:
		<-ticker.C
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				return
			}
		default:
			{
				goto CHECK
			}
		}
	CHECK:
		{
			if abstractDevice.Details() == nil {
				return
			}
			tryIfRestartDevice(abstractDevice, e, deviceInfo.UUID)
			goto TICKER
		}

	}(typex.GCTX)
	glogger.GLogger.Infof("device [%v, %v] load successfully", deviceInfo.Name, deviceInfo.UUID)
	return nil
}

func startDevice(abstractDevice typex.XDevice, e *RuleEngine) error {
	ctx, cancelCTX := typex.NewCCTX()
	if err := abstractDevice.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("abstractDevice start error:", err)
		return err
	}
	if abstractDevice.Driver() != nil {
		if abstractDevice.Driver().State() == typex.DRIVER_UP {
			abstractDevice.Driver().Stop()
		}
	}
	//----------------------------------
	// 驱动也要停了, 然后重启
	//----------------------------------
	if abstractDevice.Driver() != nil {
		if abstractDevice.Driver().State() == typex.DRIVER_UP {
			abstractDevice.Driver().Stop()
		}
		// Start driver
		// TODO: map[string]string{} 未来会被替换成真实的驱动配置
		// if driverConfig != nil --> Driver().Init(Cfg)
		//
		if err := abstractDevice.Driver().Init(map[string]string{}); err != nil {
			glogger.GLogger.Error("Driver initial error:", err)
			return errors.New("Driver initial error:" + err.Error())
		}
		glogger.GLogger.Infof("Try to start driver: [%v]", abstractDevice.Driver().DriverDetail().Name)
		if err := abstractDevice.Driver().Work(); err != nil {
			glogger.GLogger.Error("Driver work error:", err)
			return errors.New("Driver work error:" + err.Error())
		}
		glogger.GLogger.Infof("Driver start successfully: [%v]", abstractDevice.Driver().DriverDetail().Name)
	}
	return nil
}

func tryIfRestartDevice(abstractDevice typex.XDevice, e *RuleEngine, devId string) {
	if abstractDevice.Status() == typex.DEV_STOP {
		return
	}
	checkDeviceDriverState(abstractDevice)
	// 当内存里面的设备状态已经停止的时候，及时更新数据库里的
	// 此处本质上是个同步过程
	if abstractDevice.Status() == typex.DEV_DOWN {
		abstractDevice.Details().State = typex.DEV_DOWN
		glogger.GLogger.Warnf("Device %v %v down. try to restart it",
			abstractDevice.Details().UUID, abstractDevice.Details().Name)
		abstractDevice.Stop()
		abstractDevice = nil // 释放资源
		runtime.Gosched()
		runtime.GC()
		startDevice(abstractDevice, e)
	} else {
		abstractDevice.Details().State = typex.DEV_UP
	}

}

/*
*
* 检查是否需要重新拉起资源
* 这里也有优化点：不能手动控制内存回收可能会产生垃圾
*
 */
func checkDeviceDriverState(abstractDevice typex.XDevice) {

	if abstractDevice.Driver() == nil {
		return
	}
	if abstractDevice.Status() == typex.DEV_STOP {
		return
	}
	// 只有资源启动状态才拉起驱动
	if abstractDevice.Status() == typex.DEV_UP {
		// 必须资源启动, 驱动才有重启意义
		if abstractDevice.Driver().State() == typex.DRIVER_DOWN {
			glogger.GLogger.Warn("Driver down:", abstractDevice.Driver().DriverDetail().Name)
			// 只需要把资源给拉闸, 就会触发重启
			abstractDevice.Stop()
		}
	}

}

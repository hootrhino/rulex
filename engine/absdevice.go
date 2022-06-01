package engine

import (
	"context"
	"fmt"
	"rulex/device"
	"rulex/typex"
	"runtime"
	"sync"
	"time"

	"github.com/ngaut/log"
)

//--------------------------------------------------------------------------------------------------
// Abstract device
//--------------------------------------------------------------------------------------------------

//
// 获取设备
//
func (e *RuleEngine) GetDevice(id string) *typex.Device {
	v, ok := e.Devices.Load(id)
	if ok {
		return v.(*typex.Device)
	} else {
		return nil
	}

}

//
// 保存设备
//
func (e *RuleEngine) SaveDevice(dev *typex.Device) {
	e.Devices.Store(dev.UUID, dev)
}

//
// 获取所有外挂设备
//
func (e *RuleEngine) AllDevices() *sync.Map {
	return e.Devices

}

//
// 删除设备
//
func (e *RuleEngine) RemoveDevice(uuid string) {
	if dev := e.GetDevice(uuid); dev != nil {
		dev.Device.Stop()
		e.OutEnds.Delete(uuid)
		dev = nil
		log.Infof("Device [%v] has been deleted", uuid)
	}
}

//
//
//
func (e *RuleEngine) LoadDevice(deviceInfo *typex.Device) error {
	//
	// TODO `SIMPLE` Just for development stage; it will be deleted before tag
	//
	if deviceInfo.Type == "SIMPLE" {
		return startDevices(device.NewSimpleDevice(deviceInfo.UUID, e), deviceInfo, e)
	}
	if deviceInfo.Type == "TS200V02" {
		return startDevices(device.NewTS200Sensor(deviceInfo.UUID, e), deviceInfo, e)
	}
	return fmt.Errorf("unsupported InEnd type:%s", deviceInfo.Type)

}

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
	if err := startDevice(abstractDevice, e, deviceInfo.UUID); err != nil {
		log.Error(err)
		e.RemoveDevice(deviceInfo.UUID)
		return err
	}
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	go func(ctx context.Context) {
		// 5 seconds
	TICKER:
		<-ticker.C
		select {
		case <-ctx.Done():
			{
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
	log.Infof("device [%v, %v] load successfully", deviceInfo.Name, deviceInfo.UUID)
	return nil
}

//
//
//
func startDevice(abstractDevice typex.XDevice, e *RuleEngine, devId string) error {
	ctx, cancelCTX := typex.NewCCTX()
	if err := abstractDevice.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		log.Error("Source start error:", err)
		return err
	}
	return nil
}

func tryIfRestartDevice(abstractDevice typex.XDevice, e *RuleEngine, devId string) {
	// 当内存里面的设备状态已经停止的时候，及时更新数据库里的
	// 此处本质上是个同步过程
	if abstractDevice.Status() == typex.DEV_STOP {
		abstractDevice.Details().State = typex.DEV_STOP
		log.Warnf("Device %v %v down. try to restart it", abstractDevice.Details().UUID, abstractDevice.Details().Name)
		abstractDevice.Stop()
		runtime.Gosched()
		runtime.GC()
		startDevice(abstractDevice, e, devId)
	} else {
		abstractDevice.Details().State = typex.DEV_RUNNING
	}
}

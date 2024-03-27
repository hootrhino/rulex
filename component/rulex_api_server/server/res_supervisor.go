package server

import (
	"context"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/component/internotify"
	"github.com/hootrhino/rulex/component/supervisor"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 南向资源监控器 5秒检查一下状态
*
 */
func StartInSupervisor(InCtx context.Context, in *typex.InEnd, ruleEngine typex.RuleX) {
	UUID := in.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(in.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For InEnd:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)
	for {
		select {
		case <-typex.GCTX.Done():
			{
				glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-InCtx.Done():
			{
				glogger.GLogger.Debugf("Source Context cancel:%v, supervisor exit", UUID)
				return
			}
		default:
			{
			}
		}
		// 被删除后就直接退出监督进程
		currentIn := ruleEngine.GetInEnd(UUID)
		if currentIn == nil {
			glogger.GLogger.Debugf("Source:%v Deleted, supervisor exit", UUID)
			return
		}
		// STOP 设计特殊状态,标记被彻底删除的资源
		// 资源可能不会及时DOWN
		if currentIn.Source.Status() == typex.SOURCE_DOWN {
			info := fmt.Sprintf("Source:%v DOWN, supervisor try to Restart", UUID)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  "SOURCE",
				Event: "event.down",
				Ts:    uint64(time.Now().UnixNano()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestInEnd(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

/*
*
* 北向资源监控器 5秒检查一下状态
*
 */
func StartOutSupervisor(OutCtx context.Context, out *typex.OutEnd, ruleEngine typex.RuleX) {
	UUID := out.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(out.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For OutEnd:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)

	for {
		select {
		case <-typex.GCTX.Done():
			glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
			return
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-OutCtx.Done():
			glogger.GLogger.Debugf("OutEnd Context cancel:%v, supervisor exit", UUID)
			return
		default:
			{
			}
		}
		// 被删除后就直接退出监督进程
		currentOut := ruleEngine.GetOutEnd(UUID)
		if currentOut == nil {
			glogger.GLogger.Debugf("OutEnd:%v Deleted, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentOut.Target.Status() == typex.SOURCE_DOWN {
			info := fmt.Sprintf("OutEnd:%v DOWN, supervisor try to Restart", UUID)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  "TARGET",
				Event: "event.down",
				Ts:    uint64(time.Now().UnixNano()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestOutEnd(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

/*
*
* 设备监控器 5秒检查一下状态
*
 */
func StartDeviceSupervisor(DeviceCtx context.Context, device *typex.Device, ruleEngine typex.RuleX) {
	UUID := device.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer ticker.Stop()
	SuperVisor := supervisor.RegisterSuperVisor(device.UUID)
	glogger.GLogger.Debugf("Register SuperVisor For Device:%s", SuperVisor.SlaverId)
	defer supervisor.UnRegisterSuperVisor(SuperVisor.SlaverId)

	for {
		select {
		case <-typex.GCTX.Done():
			{
				glogger.GLogger.Debugf("Global Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-SuperVisor.Ctx.Done():
			{
				glogger.GLogger.Debugf("SuperVisor Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-DeviceCtx.Done():
			{
				glogger.GLogger.Debugf("Device Context cancel:%v, supervisor exit", UUID)
				return
			}
		default:
			{
			}
		}
		// 被删除后就直接退出监督进程
		currentDevice := ruleEngine.GetDevice(UUID)
		if currentDevice == nil {
			glogger.GLogger.Debugf("Device:%v Deleted, supervisor exit", UUID)
			return
		}

		// 资源可能不会及时DOWN
		if currentDevice.Device.Status() == typex.DEV_DOWN {
			info := fmt.Sprintf("Device:%v DOWN, supervisor try to Restart", UUID)
			glogger.GLogger.Debugf(info)
			internotify.Push(internotify.BaseEvent{
				Type:  "DEVICE",
				Event: "event.down",
				Ts:    uint64(time.Now().UnixNano()),
				Info:  info,
			})
			time.Sleep(4 * time.Second)
			go LoadNewestDevice(UUID, ruleEngine)
			return
		}
		<-ticker.C
	}
}

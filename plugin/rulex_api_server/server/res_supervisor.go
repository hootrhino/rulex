package server

import (
	"context"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

func StartInSupervisor(ctx context.Context, in *typex.InEnd, ruleEngine typex.RuleX) {
	UUID := in.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				glogger.GLogger.Debugf("Source Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-typex.GCTX.Done():
			{
				return
			}
		default:
			{
			}
		}
		currentIn := ruleEngine.GetInEnd(UUID)
		if currentIn == nil {
			glogger.GLogger.Debugf("Source:%v Deleted, supervisor exit", UUID)
			return
		}
		if currentIn.Source.Status() == typex.SOURCE_STOP {
			glogger.GLogger.Debugf("Source:%v Stopped, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentIn.Source.Status() == typex.SOURCE_DOWN {
			glogger.GLogger.Debugf("Source:%v DOWN, supervisor try to Restart", UUID)
			time.Sleep(2 * time.Second)
			go LoadNewestInEnd(UUID, ruleEngine)
			return
		}
		// glogger.GLogger.Debugf("Supervisor Get Source :%v state:%v", UUID, currentIn.Source.Status().String())
		<-ticker.C
	}
}
func StartOutSupervisor(ctx context.Context, out *typex.OutEnd, ruleEngine typex.RuleX) {
	UUID := out.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				glogger.GLogger.Debugf("OutEnd Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-typex.GCTX.Done():
			{
				return
			}
		default:
			{
			}
		}
		currentOut := ruleEngine.GetOutEnd(UUID)
		if currentOut == nil {
			glogger.GLogger.Debugf("OutEnd:%v Deleted, supervisor exit", UUID)
			return
		}
		if currentOut.Target.Status() == typex.SOURCE_STOP {
			glogger.GLogger.Debugf("OutEnd:%v Stopped, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentOut.Target.Status() == typex.SOURCE_DOWN {
			glogger.GLogger.Debugf("OutEnd:%v DOWN, supervisor try to Restart", UUID)
			time.Sleep(5 * time.Second)
			go LoadNewestOutEnd(UUID, ruleEngine)
			return
		}
		// glogger.GLogger.Debugf("Supervisor Get OutEnd :%v state:%v", UUID, currentOut.Target.Status().String())
		<-ticker.C
	}
}
func StartDeviceSupervisor(ctx context.Context, device *typex.Device, ruleEngine typex.RuleX) {
	UUID := device.UUID
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				glogger.GLogger.Debugf("Device Context cancel:%v, supervisor exit", UUID)
				return
			}
		case <-typex.GCTX.Done():
			{
				return
			}
		default:
			{
			}
		}
		currentDevice := ruleEngine.GetDevice(UUID)
		if currentDevice == nil {
			glogger.GLogger.Debugf("Device:%v Deleted, supervisor exit", UUID)
			return
		}
		if currentDevice.Device.Status() == typex.DEV_STOP {
			glogger.GLogger.Debugf("Device:%v Stopped, supervisor exit", UUID)
			return
		}
		// 资源可能不会及时DOWN
		if currentDevice.Device.Status() == typex.DEV_DOWN {
			glogger.GLogger.Debugf("Device:%v DOWN, supervisor try to Restart", UUID)
			time.Sleep(2 * time.Second)
			go LoadNewestDevice(UUID, ruleEngine)
			return
		}
		// glogger.GLogger.Debugf("Supervisor Get Device :%v state:%v", UUID, currentDevice.Device.Status().String())
		<-ticker.C
	}
}

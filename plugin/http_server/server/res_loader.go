package server

import (
	"errors"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/square/go-jose.v2/json"
)

/*
*
* 当资源重启加载的时候，内存里面的数据会丢失，需要重新从数据库加载规则到资源，建立绑定关联。
*
 */

// LoadNewestInEnd
func LoadNewestInEnd(uuid string, ruleEngine typex.RuleX) error {
	mInEnd, _ := service.GetMInEndWithUUID(uuid)
	if mInEnd == nil {
		return errors.New("Inend not exists:" + uuid)
	}
	config := map[string]interface{}{}
	if err1 := json.Unmarshal([]byte(mInEnd.Config), &config); err1 != nil {
		glogger.GLogger.Error(err1)
		return err1
	}
	// 所有的更新都先停止资源,然后再加载
	old := ruleEngine.GetInEnd(uuid)
	if old != nil {
		if old.Source.Status() == typex.SOURCE_UP {
			old.Source.Details().State = typex.SOURCE_STOP
			old.Source.Stop()
		}
	}
	ruleEngine.RemoveInEnd(uuid)
	in := typex.NewInEnd(typex.InEndType(mInEnd.Type),
		mInEnd.Name, mInEnd.Description, mInEnd.GetConfig())
	// Important !!!!!!!! in.Id = mInEnd.UUID
	in.UUID = mInEnd.UUID
	// 未来会支持XDataModel数据模型, 目前暂时留空
	in.DataModelsMap = map[string]typex.XDataModel{}
	BindRules := map[string]typex.Rule{}
	for _, ruleId := range mInEnd.BindRules {
		if ruleId == "" {
			continue
		}
		mRule, err1 := service.GetMRuleWithUUID(ruleId)
		if err1 != nil {
			return err1
		}
		glogger.GLogger.Debugf("Load rule:%s", mRule.Name)
		RuleInstance := typex.NewLuaRule(
			ruleEngine,
			mRule.UUID,
			mRule.Name,
			mRule.Description,
			mRule.FromSource,
			mRule.FromDevice,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		BindRules[mRule.UUID] = *RuleInstance
	}
	// 最新的规则
	in.BindRules = BindRules
	// 最新的配置
	in.Config = mInEnd.GetConfig()
	ctx, cancelCTX := typex.NewCCTX()
	if err2 := ruleEngine.LoadInEndWithCtx(in, ctx, cancelCTX); err2 != nil {
		glogger.GLogger.Error(err2)
		return err2
	}
	go StartInSupervisor(ctx, in, ruleEngine)
	return nil
}

// LoadNewestOutEnd
func LoadNewestOutEnd(uuid string, ruleEngine typex.RuleX) error {
	mOutEnd, err := service.GetMOutEndWithUUID(uuid)
	if err != nil {
		return err
	}

	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	old := ruleEngine.GetOutEnd(uuid)
	if old != nil {
		old.Target.Stop()
	}
	ruleEngine.RemoveOutEnd(uuid)
	out := typex.NewOutEnd(typex.TargetType(mOutEnd.Type),
		mOutEnd.Name, mOutEnd.Description, config)
	// Important !!!!!!!!
	out.UUID = mOutEnd.UUID
	out.Config = mOutEnd.GetConfig()
	ctx, cancelCTX := typex.NewCCTX()
	if err := ruleEngine.LoadOutEndWithCtx(out, ctx, cancelCTX); err != nil {
		return err
	}
	go StartOutSupervisor(ctx, out, ruleEngine)
	return nil

}

/*
*
* 当资源重启加载的时候，内存里面的数据会丢失，需要重新从数据库加载规则到资源，建立绑定关联。
*
 */

// LoadNewestDevice
func LoadNewestDevice(uuid string, ruleEngine typex.RuleX) error {
	mDevice, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		return err
	}
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mDevice.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	old := ruleEngine.GetDevice(uuid)
	if old != nil {
		old.Device.Stop()
	}
	ruleEngine.RemoveDevice(uuid) // 删除内存里面的
	dev := typex.NewDevice(typex.DeviceType(mDevice.Type), mDevice.Name,
		mDevice.Description, mDevice.GetConfig())
	// Important !!!!!!!!
	dev.UUID = mDevice.UUID // 本质上是配置和内存的数据映射起来
	BindRules := map[string]typex.Rule{}
	for _, ruleId := range mDevice.BindRules {
		if ruleId == "" {
			continue
		}
		mRule, err1 := service.GetMRuleWithUUID(ruleId)
		if err1 != nil {
			return err1
		}
		glogger.GLogger.Debugf("Load rule:%s", mRule.Name)
		RuleInstance := typex.NewLuaRule(
			ruleEngine,
			mRule.UUID,
			mRule.Name,
			mRule.Description,
			mRule.FromSource,
			mRule.FromDevice,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		BindRules[mRule.UUID] = *RuleInstance
	}
	// 最新的规则
	dev.BindRules = BindRules
	// 最新的配置
	dev.Config = mDevice.GetConfig()
	// 参数传给 --> startDevice()
	ctx, cancelCTX := typex.NewCCTX()
	if err := ruleEngine.LoadDeviceWithCtx(dev, ctx, cancelCTX); err != nil {
		return err
	}
	go StartDeviceSupervisor(ctx, dev, ruleEngine)
	return nil

}

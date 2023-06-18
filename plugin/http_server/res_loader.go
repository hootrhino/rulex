package httpserver

import (
	"errors"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"gopkg.in/square/go-jose.v2/json"
)

/*
*
* 当资源重启加载的时候，内存里面的数据会丢失，需要重新从数据库加载规则到资源，建立绑定关联。
*
 */

// LoadNewestInEnd
func (hh *HttpApiServer) LoadNewestInEnd(uuid string) error {
	mInEnd, _ := hh.GetMInEndWithUUID(uuid)
	if mInEnd == nil {
		return errors.New("Inend not exists:" + uuid)
	}
	config := map[string]interface{}{}
	if err1 := json.Unmarshal([]byte(mInEnd.Config), &config); err1 != nil {
		glogger.GLogger.Error(err1)
		return err1
	}
	// :mInEnd: {k1 :{k1:v1}, k2 :{k2:v2}} --> InEnd: [{k1:v1}, {k2:v2}]
	// var dataModelsMap map[string]typex.XDataModel
	// if err1 := json.Unmarshal([]byte(mInEnd.XDataModels), &dataModelsMap); err1 != nil {
	// 	glogger.GLogger.Error(err1)
	// 	return err1
	// }
	// 所有的更新都先停止资源,然后再加载
	hh.ruleEngine.RemoveInEnd(uuid)
	in := typex.NewInEnd(typex.InEndType(mInEnd.Type),
		mInEnd.Name, mInEnd.Description, config)
	// Important !!!!!!!! in.Id = mInEnd.UUID
	in.UUID = mInEnd.UUID
	in.Config = mInEnd.GetConfig()
	// 未来会支持XDataModel数据模型, 目前暂时留空
	in.DataModelsMap = map[string]typex.XDataModel{}
	BindRules := map[string]typex.Rule{}
	for _, ruleId := range mInEnd.BindRules {
		if ruleId == "" {
			continue
		}
		mRule, err1 := hh.GetMRuleWithUUID(ruleId)
		if err1 != nil {
			return err1
		}
		glogger.GLogger.Debugf("Load rule:%s", mRule.Name)
		RuleInstance := typex.NewLuaRule(
			hh.ruleEngine,
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
	if err2 := hh.ruleEngine.LoadInEnd(in); err2 != nil {
		glogger.GLogger.Error(err2)
		return err2
	}

	return nil
}

// LoadNewestOutEnd
func (hh *HttpApiServer) LoadNewestOutEnd(uuid string) error {
	mOutEnd, _ := hh.GetMOutEndWithUUID(uuid)
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	hh.ruleEngine.RemoveOutEnd(uuid)
	out := typex.NewOutEnd(typex.TargetType(mOutEnd.Type),
		mOutEnd.Name, mOutEnd.Description, config)
	// Important !!!!!!!!
	out.UUID = mOutEnd.UUID
	out.Config = mOutEnd.GetConfig()
	if err := hh.ruleEngine.LoadOutEnd(out); err != nil {
		return err
	} else {
		return nil
	}

}

/*
*
* 当资源重启加载的时候，内存里面的数据会丢失，需要重新从数据库加载规则到资源，建立绑定关联。
*
 */

// LoadNewestDevice
func (hh *HttpApiServer) LoadNewestDevice(uuid string) error {
	mDevice, err := hh.GetDeviceWithUUID(uuid)
	if err != nil {
		return err
	}
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mDevice.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	hh.ruleEngine.RemoveDevice(uuid) // 删除内存里面的
	dev := typex.NewDevice(typex.DeviceType(mDevice.Type), mDevice.Name,
		mDevice.Description, mDevice.GetConfig())
	// Important !!!!!!!!
	dev.UUID = mDevice.UUID // 本质上是配置和内存的数据映射起来
	BindRules := map[string]typex.Rule{}
	for _, ruleId := range mDevice.BindRules {
		if ruleId == "" {
			continue
		}
		mRule, err1 := hh.GetMRuleWithUUID(ruleId)
		if err1 != nil {
			return err1
		}
		glogger.GLogger.Debugf("Load rule:%s", mRule.Name)
		RuleInstance := typex.NewLuaRule(
			hh.ruleEngine,
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
	if err := hh.ruleEngine.LoadDevice(dev); err != nil {
		return err
	}
	return nil

}

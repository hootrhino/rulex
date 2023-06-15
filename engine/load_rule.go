package engine

import (
	"errors"
	"sync"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

// LoadRule: 每个规则都绑定了资源(FromSource)或者设备(FromDevice)
// 使用MAP来记录RULE的绑定关系, KEY是UUID, Value是规则
func (e *RuleEngine) LoadRule(r *typex.Rule) error {
	// 前置语法验证
	if err := core.VerifyLuaSyntax(r); err != nil {
		return err
	}
	// 前置自定义库校验
	if err := LoadExtLuaLib(e, r); err != nil {
		return err
	}
	e.SaveRule(r)
	//--------------------------------------------------------------
	// Load LoadBuildInLuaLib
	//--------------------------------------------------------------
	LoadBuildInLuaLib(e, r)

	glogger.GLogger.Infof("Rule [%v, %v] load successfully", r.Name, r.UUID)
	// 绑定输入资源
	for _, inUUId := range r.FromSource {
		// 查找输入定义的资源是否存在
		if in := e.GetInEnd(inUUId); in != nil {
			(in.BindRules)[r.UUID] = *r
			return nil
		} else {
			return errors.New("'InEnd':" + inUUId + " is not working now")
		}
	}
	// 绑定设备
	for _, devUUId := range r.FromDevice {
		// 查找输入定义的资源是否存在
		if Device := e.GetDevice(devUUId); Device != nil {
			// 绑定资源和规则，建立关联关系
			(Device.BindRules)[r.UUID] = *r
		} else {
			return errors.New("'Device':" + devUUId + " is not working now")
		}
	}
	return nil

}

// GetRule a rule
func (e *RuleEngine) GetRule(id string) *typex.Rule {
	v, ok := (e.Rules).Load(id)
	if ok {
		return v.(*typex.Rule)
	} else {
		return nil
	}
}

func (e *RuleEngine) SaveRule(r *typex.Rule) {
	e.Rules.Store(r.UUID, r)
}

// RemoveRule and inend--rule bindings
func (e *RuleEngine) RemoveRule(ruleId string) {
	if rule := e.GetRule(ruleId); rule != nil {
		// 清空 InEnd 的 bind 资源
		e.AllInEnd().Range(func(key, value interface{}) bool {
			inEnd := value.(*typex.InEnd)
			for _, r := range inEnd.BindRules {
				if rule.UUID == r.UUID {
					delete(inEnd.BindRules, ruleId)
				}
			}
			return true
		})
		// 清空Device的绑定
		e.AllDevices().Range(func(key, value interface{}) bool {
			Device := value.(*typex.Device)
			for _, r := range Device.BindRules {
				if rule.UUID == r.UUID {
					delete(Device.BindRules, ruleId)
				}
			}
			return true
		})
		e.Rules.Delete(ruleId)
		glogger.GLogger.Infof("Rule [%v] has been deleted", ruleId)
	}
}

func (e *RuleEngine) AllRule() *sync.Map {
	return e.Rules
}

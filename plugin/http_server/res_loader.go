package httpserver

import (
	"errors"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"gopkg.in/square/go-jose.v2/json"
)

//
// LoadNewestInEnd
//
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
	var dataModelsMap map[string]typex.XDataModel
	if err1 := json.Unmarshal([]byte(mInEnd.XDataModels), &dataModelsMap); err1 != nil {
		glogger.GLogger.Error(err1)
		return err1
	}
	// 所有的更新都先停止资源,然后再加载
	hh.ruleEngine.RemoveInEnd(uuid)
	in := typex.NewInEnd(typex.InEndType(mInEnd.Type), mInEnd.Name, mInEnd.Description, config)
	// Important !!!!!!!! in.Id = mInEnd.UUID
	in.UUID = mInEnd.UUID
	in.DataModelsMap = dataModelsMap
	if err2 := hh.ruleEngine.LoadInEnd(in); err2 != nil {
		glogger.GLogger.Error(err2)
		return err2
	} else {
		return nil
	}

}

//
// LoadNewestOutEnd
//
func (hh *HttpApiServer) LoadNewestOutEnd(uuid string) error {
	mOutEnd, _ := hh.GetMOutEndWithUUID(uuid)
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mOutEnd.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	hh.ruleEngine.RemoveOutEnd(uuid)
	out := typex.NewOutEnd(typex.TargetType(mOutEnd.Type), mOutEnd.Name, mOutEnd.Description, config)
	// Important !!!!!!!!
	out.UUID = mOutEnd.UUID
	if err := hh.ruleEngine.LoadOutEnd(out); err != nil {
		return err
	} else {
		return nil
	}

}

//
// LoadNewestDevice
//
func (hh *HttpApiServer) LoadNewestDevice(uuid string) error {
	mDevice, _ := hh.GetDeviceWithUUID(uuid)
	config := map[string]interface{}{}
	if err := json.Unmarshal([]byte(mDevice.Config), &config); err != nil {
		return err
	}
	// 所有的更新都先停止资源,然后再加载
	hh.ruleEngine.RemoveDevice(uuid)
	dev := typex.NewDevice(typex.DeviceType(mDevice.Type), mDevice.Name, mDevice.Description, mDevice.ActionScript, config)
	// Important !!!!!!!!
	dev.UUID = mDevice.UUID // 本质上是配置和内存的数据映射起来
	if err := hh.ruleEngine.LoadDevice(dev); err != nil {
		return err
	} else {
		return nil
	}

}

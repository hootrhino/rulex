# Rulex 基础组件扩展模板

为了开发省事，提取了一些代码模板，可根据实际情况配置进代码模板工具。

## 插件
```golang
package demo_plugin

import "rulex/core"

type DemoPlugin struct {
}

func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
}
func (dm *DemoPlugin) Load() *core.XPluginEnv {
	return core.NewXPluginEnv()
}
func (dm *DemoPlugin) Init(*core.XPluginEnv) error {
	return nil

}
//
func (dm *DemoPlugin) Install(*core.XPluginEnv) (*core.XPluginMetaInfo, error) {
	return &core.XPluginMetaInfo{
		Name:     "DemoPlugin",
		Version:  "0.0.1",
		Homepage: "www.ezlinker.cn",
		HelpLink: "www.ezlinker.cn",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}, nil
}
func (dm *DemoPlugin) Start(*core.XPluginEnv) error {
	return nil
}
func (dm *DemoPlugin) Uninstall(*core.XPluginEnv) error {
	return nil
}
func (dm *DemoPlugin) Clean() {

}
```

## 入口
```golang
package core

import (
	"time"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

type SerialResource struct {
	XStatus
}

func NewSerialResource(inEndId string, e *RuleEngine) *SerialResource {
	s := SerialResource{}
	return &s
}

func (mm *SerialResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{}
}

func (s *SerialResource) Test(inEndId string) bool {
	return true
}

func (s *SerialResource) Register(inEndId string) error {
	return nil
}

func (s *SerialResource) Start() error {
    return nil
}

func (s *SerialResource) Enabled() bool {
	return true
}

func (s *SerialResource) Reload() {
}

func (s *SerialResource) Pause() {

}

func (s *SerialResource) Status() State {
	return UP
}

func (s *SerialResource) Stop() {
}
```

## 出口
```golang
  
package core

import (
	"context"

	"github.com/ngaut/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//
type DemoTarget struct {
	XStatus
}

func NewDemoTarget(e *RuleEngine) *DemoTarget {
	d := new(DemoTarget)
	d.ruleEngine = e
	return mg
}

func (m *DemoTarget) Register(outEndId string) error {
	return nil
}

func (m *DemoTarget) Start() error {
	return nil
}

func (m *DemoTarget) Test(outEndId string) bool {
    return true
}

func (m *DemoTarget) Enabled() bool {
	return m.Enable
}

func (m *DemoTarget) Reload() {
}

func (m *DemoTarget) Pause() {

}

func (m *DemoTarget) Status() State {
	return UP
}

func (m *DemoTarget) Stop() {
}

func (m *DemoTarget) To(data interface{}) error {
	return nil
}
```

## Hook
```golang

type DemoHook struct {
    //
}
func (h*DemoHook) Work(data string) error{
    return nil
}
func (h*DemoHook) Error(err error){

}
func (h*DemoHook) Name() string{
    return "DemoHook"
}
```
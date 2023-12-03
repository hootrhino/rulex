// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package source

import (
	"context"
	"encoding/json"

	"github.com/hootrhino/rulex/component/internotify"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 用来将内部消息总线的事件推到资源脚本
*
 */
type __InternalEventSourceConfig struct {
	// - ALL: 全部事件
	// - SOURCE: 南向事件
	// - DEVICE: 设备事件
	// - TARGET: 北向事件
	// - SYSTEM: 系统事件
	// - HARDWARE: 硬件事件
	Type string `json:"type"`
}
type InternalEventSource struct {
	typex.XStatus
	mainConfig __InternalEventSourceConfig
}

func NewInternalEventSource(r typex.RuleX) typex.XSource {
	s := InternalEventSource{}
	s.mainConfig = __InternalEventSourceConfig{
		Type: "ALL",
	}
	s.RuleEngine = r
	return &s
}

func (u *InternalEventSource) Init(inEndId string, configMap map[string]interface{}) error {
	u.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &u.mainConfig); err != nil {
		return err
	}
	return nil
}

func (u *InternalEventSource) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX
	u.startInternalEventQueue()
	return nil

}
func (u *InternalEventSource) DataModels() []typex.XDataModel {
	return u.XDataModels
}

func (u *InternalEventSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (u *InternalEventSource) Stop() {
	if u.CancelCTX != nil {
		u.CancelCTX()
	}
}
func (*InternalEventSource) Driver() typex.XExternalDriver {
	return nil
}

func (u *InternalEventSource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *InternalEventSource) Test(inEndId string) bool {
	return true
}

// 拓扑
func (*InternalEventSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

// 来自外面的数据
func (*InternalEventSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*InternalEventSource) UpStream([]byte) (int, error) {
	return 0, nil
}

type event struct {
	// - ALL: 全部
	// - SOURCE: 南向事件
	// - DEVICE: 设备事件
	// - TARGET: 北向事件
	// - SYSTEM: 系统内部事件
	// - HARDWARE: 硬件事件
	Type  string      `json:"type"`
	Event string      `json:"event"`
	Ts    uint64      `json:"ts"`
	Info  interface{} `json:"info"`
}

/*
*
* 从内部总线拿数据
*
 */
func (ie *InternalEventSource) startInternalEventQueue() {
	go func(ctx context.Context) {
		internotify.AddSource()
		defer internotify.RemoveSource()
		for {
			select {
			case <-ctx.Done():
				return
			case Event := <-internotify.GetQueue():
				if ie.mainConfig.Type == "ALL" {
					bytes, _ := json.Marshal(event{
						Type:  Event.Type,
						Event: Event.Event,
						Ts:    Event.Ts,
						Info:  Event.Info,
					})
					ie.RuleEngine.WorkInEnd(ie.RuleEngine.GetInEnd(ie.PointId), string(bytes))
					continue
				}
				if ie.mainConfig.Type == Event.Type {
					bytes, _ := json.Marshal(event{
						Type:  Event.Type,
						Event: Event.Event,
						Ts:    Event.Ts,
						Info:  Event.Info,
					})
					ie.RuleEngine.WorkInEnd(ie.RuleEngine.GetInEnd(ie.PointId), string(bytes))
					continue
				}
			}
		}
	}(ie.Ctx)
}

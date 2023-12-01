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
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type __InternalEventSourceConfig struct {
}
type InternalEventSource struct {
	typex.XStatus
	mainConfig __InternalEventSourceConfig
}

func NewInternalEventSource(r typex.RuleX) typex.XSource {
	s := InternalEventSource{}
	s.RuleEngine = r
	return &s
}

func (u *InternalEventSource) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX

	return nil

}

func (u *InternalEventSource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *InternalEventSource) Test(inEndId string) bool {
	return true
}

func (u *InternalEventSource) Init(inEndId string, configMap map[string]interface{}) error {
	u.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &u.mainConfig); err != nil {
		return err
	}
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

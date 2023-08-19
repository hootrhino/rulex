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

package core

import "github.com/hootrhino/rulex/typex"

type TargetTypeManager struct {
	// K: 资源类型
	// V: 伪构造函数
	registry map[typex.TargetType]*typex.XConfig
}

func NewTargetTypeManager() typex.TargetRegistry {
	return &TargetTypeManager{
		registry: map[typex.TargetType]*typex.XConfig{},
	}

}
func (rm *TargetTypeManager) Register(name typex.TargetType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *TargetTypeManager) Find(name typex.TargetType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *TargetTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

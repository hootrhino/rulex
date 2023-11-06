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

package datacenter

import "github.com/hootrhino/rulex/typex"

var __DefaultDataCenter *DataCenter

/*
*
* 留着未来扩充数据中心的功能
*
 */
type DataCenter struct {
	ExternalDb ExternalDb
	LocalDb    LocalDb
	rulex      typex.RuleX
}

func InitDataCenter(rulex typex.RuleX) {
	__DefaultDataCenter = new(DataCenter)
	__DefaultDataCenter.ExternalDb = ExternalDb{}
	__DefaultDataCenter.LocalDb = LocalDb{}
	__DefaultDataCenter.rulex = rulex
}

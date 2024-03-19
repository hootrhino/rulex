// Copyright (C) 2024 wwhai
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

package supervisor

import (
	"context"
	"sync"

	"github.com/hootrhino/rulex/typex"
)

var __DefaultSuperVisorAdmin *SuperVisorAdmin

type SuperVisor struct {
	SlaverId string
	Ctx      context.Context
	Cancel   context.CancelFunc
}
type SuperVisorAdmin struct {
	Ctx         context.Context
	Locker      sync.Locker
	SuperVisors map[string]*SuperVisor
	rulex       typex.RuleX
}

/*
*
* 初始化超级Admin
*
 */
func InitResourceSuperVisorAdmin(rulex typex.RuleX) {
	__DefaultSuperVisorAdmin = &SuperVisorAdmin{
		Ctx:         context.Background(),
		Locker:      &sync.Mutex{},
		SuperVisors: map[string]*SuperVisor{},
		rulex:       rulex,
	}
}

/*
*
* 启动Supervisor的时候注册
*
 */
func RegisterSuperVisor(SlaverId string) *SuperVisor {
	if _, Ok := __DefaultSuperVisorAdmin.SuperVisors[SlaverId]; !Ok {
		Ctx, Cancel := context.WithCancel(context.Background())
		SuperVisor := &SuperVisor{SlaverId, Ctx, Cancel}
		__DefaultSuperVisorAdmin.SuperVisors[SlaverId] = SuperVisor
		return SuperVisor
	}
	return nil
}

/*
*
* Supervisor进程退出的时候执行
*
 */
func UnRegisterSuperVisor(UUID string) {
	__DefaultSuperVisorAdmin.Locker.Lock()
	defer __DefaultSuperVisorAdmin.Locker.Unlock()
	if Sv, Ok := __DefaultSuperVisorAdmin.SuperVisors[UUID]; Ok {
		Sv.Cancel()
		delete(__DefaultSuperVisorAdmin.SuperVisors, UUID)
		Sv = nil
	}
}

/*
*
* 停止一个Supervisor
*
 */
func StopSuperVisor(UUID string) {
	__DefaultSuperVisorAdmin.Locker.Lock()
	defer __DefaultSuperVisorAdmin.Locker.Unlock()
	if Sv, Ok := __DefaultSuperVisorAdmin.SuperVisors[UUID]; Ok {
		Sv.Cancel()
	}
}

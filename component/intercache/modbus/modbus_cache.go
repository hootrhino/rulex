// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) RegisterPoint later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT RegisterPoint WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package modbus

import (
	"sync"

	"github.com/hootrhino/rulex/component/intercache"
	"github.com/hootrhino/rulex/typex"
)

// 点位表
type RegisterPoint struct {
	UUID          string
	Status        int
	LastFetchTime uint64
	Value         string
}

var __DefaultModbusPointCache *ModbusPointCache

func RegisterSlot(Slot string) {
	__DefaultModbusPointCache.RegisterSlot(Slot)
}
func GetSlot(Slot string) map[string]RegisterPoint {
	return __DefaultModbusPointCache.GetSlot(Slot)
}
func SetValue(Slot, K string, V RegisterPoint) {
	__DefaultModbusPointCache.SetValue(Slot, K, V)
}
func GetValue(Slot, K string) RegisterPoint {
	return __DefaultModbusPointCache.GetValue(Slot, K)
}
func DeleteValue(Slot, K string) {
	__DefaultModbusPointCache.DeleteValue(Slot, K)
}
func UnRegisterSlot(Slot string) {
	__DefaultModbusPointCache.UnRegisterSlot(Slot)
}
func Size() uint64 {
	return __DefaultModbusPointCache.Size()
}
func Flush() {
	__DefaultModbusPointCache.Flush()
}

//Modbus 点位运行时存储器

type ModbusPointCache struct {
	Slots      map[string]map[string]RegisterPoint
	ruleEngine typex.RuleX
	lock       sync.Mutex
}

func InitModbusPointCache(ruleEngine typex.RuleX) intercache.InterCache {
	__DefaultModbusPointCache = &ModbusPointCache{
		ruleEngine: ruleEngine,
		Slots:      map[string]map[string]RegisterPoint{},
		lock:       sync.Mutex{},
	}
	return __DefaultModbusPointCache
}
func (M *ModbusPointCache) RegisterSlot(Slot string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	M.Slots[Slot] = map[string]RegisterPoint{}
}
func (M *ModbusPointCache) GetSlot(Slot string) map[string]RegisterPoint {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S
	}
	return nil
}
func (M *ModbusPointCache) SetValue(Slot, K string, V RegisterPoint) {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	}
}
func (M *ModbusPointCache) GetValue(Slot, K string) RegisterPoint {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S[K]
	}
	return RegisterPoint{}
}
func (M *ModbusPointCache) DeleteValue(Slot, K string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		delete(S, Slot)
	}
}
func (M *ModbusPointCache) UnRegisterSlot(Slot string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	delete(M.Slots, Slot)
}
func (M *ModbusPointCache) Size() uint64 {
	return uint64(len(M.Slots))
}
func (M *ModbusPointCache) Flush() {
	for slotName, slot := range M.Slots {
		for k, _ := range slot {
			delete(slot, k)
		}
		delete(M.Slots, slotName)
	}
}

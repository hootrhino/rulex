// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) DataSchemaValue later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT DataSchemaValue WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package dataschema

import (
	"sync"

	"github.com/hootrhino/rulex/component/intercache"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultIotSchemaCache *IotSchemaCache

func RegisterSlot(Slot string) {
	__DefaultIotSchemaCache.RegisterSlot(Slot)
}
func GetSlot(Slot string) map[string]DataSchemaValue {
	return __DefaultIotSchemaCache.GetSlot(Slot)
}
func SetValue(Slot, K string, V DataSchemaValue) {
	__DefaultIotSchemaCache.SetValue(Slot, K, V)
}
func GetValue(Slot, K string) DataSchemaValue {
	return __DefaultIotSchemaCache.GetValue(Slot, K)
}
func DeleteValue(Slot, K string) {
	__DefaultIotSchemaCache.DeleteValue(Slot, K)
}
func UnRegisterSlot(Slot string) {
	__DefaultIotSchemaCache.UnRegisterSlot(Slot)
}
func Size() uint64 {
	return __DefaultIotSchemaCache.Size()
}
func Flush() {
	__DefaultIotSchemaCache.Flush()
}

//Modbus 点位运行时存储器

type IotSchemaCache struct {
	Slots      map[string]map[string]DataSchemaValue
	ruleEngine typex.RuleX
	locker     sync.Mutex
}
type DataSchemaValue struct {
	UUID          string // Cache uuid
	Name          string // 变量关联名
	Value         any    //
	LastFetchTime uint64
}

func InitIotSchemaCache(ruleEngine typex.RuleX) intercache.InterCache {
	__DefaultIotSchemaCache = &IotSchemaCache{
		ruleEngine: ruleEngine,
		Slots:      map[string]map[string]DataSchemaValue{},
		locker:     sync.Mutex{},
	}
	return __DefaultIotSchemaCache
}
func (M *IotSchemaCache) RegisterSlot(Slot string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	M.Slots[Slot] = map[string]DataSchemaValue{}
}
func (M *IotSchemaCache) GetSlot(Slot string) map[string]DataSchemaValue {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S
	}
	return nil
}
func (M *IotSchemaCache) SetValue(Slot, K string, V DataSchemaValue) {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	}
}
func (M *IotSchemaCache) GetValue(Slot, K string) DataSchemaValue {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S[K]
	}
	return DataSchemaValue{}
}
func (M *IotSchemaCache) DeleteValue(Slot, K string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		delete(S, Slot)
	}
}
func (M *IotSchemaCache) UnRegisterSlot(Slot string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	delete(M.Slots, Slot)
	M.Flush()
}
func (M *IotSchemaCache) Size() uint64 {
	return uint64(len(M.Slots))
}
func (M *IotSchemaCache) Flush() {
	for slotName, slot := range M.Slots {
		for k := range slot {
			delete(slot, k)
		}
		delete(M.Slots, slotName)
	}
}

// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) Hnc8RegisterPoint later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT Hnc8RegisterPoint WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package hnccnc

import (
	"sync"

	"github.com/hootrhino/rulex/component/intercache"
	"github.com/hootrhino/rulex/typex"
)

// 点位表
type Hnc8RegisterPoint struct {
	UUID          string
	Status        int
	LastFetchTime uint64
	Value         string
}

var __DefaultHnc8CnCPointCache *Hnc8CnCPointCache

func RegisterSlot(Slot string) {
	__DefaultHnc8CnCPointCache.RegisterSlot(Slot)
}
func GetSlot(Slot string) map[string]Hnc8RegisterPoint {
	return __DefaultHnc8CnCPointCache.GetSlot(Slot)
}
func SetValue(Slot, K string, V Hnc8RegisterPoint) {
	__DefaultHnc8CnCPointCache.SetValue(Slot, K, V)
}
func GetValue(Slot, K string) Hnc8RegisterPoint {
	return __DefaultHnc8CnCPointCache.GetValue(Slot, K)
}
func DeleteValue(Slot, K string) {
	__DefaultHnc8CnCPointCache.DeleteValue(Slot, K)
}
func UnRegisterSlot(Slot string) {
	__DefaultHnc8CnCPointCache.UnRegisterSlot(Slot)
}
func Size() uint64 {
	return __DefaultHnc8CnCPointCache.Size()
}
func Flush() {
	__DefaultHnc8CnCPointCache.Flush()
}

//Hnc8CnC 点位运行时存储器

type Hnc8CnCPointCache struct {
	Slots      map[string]map[string]Hnc8RegisterPoint
	ruleEngine typex.RuleX
	locker     sync.Mutex
}

func InitHnc8CnCPointCache(ruleEngine typex.RuleX) intercache.InterCache {
	__DefaultHnc8CnCPointCache = &Hnc8CnCPointCache{
		ruleEngine: ruleEngine,
		Slots:      map[string]map[string]Hnc8RegisterPoint{},
		locker:     sync.Mutex{},
	}
	return __DefaultHnc8CnCPointCache
}
func (M *Hnc8CnCPointCache) RegisterSlot(Slot string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	M.Slots[Slot] = map[string]Hnc8RegisterPoint{}
}
func (M *Hnc8CnCPointCache) GetSlot(Slot string) map[string]Hnc8RegisterPoint {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S
	}
	return nil
}
func (M *Hnc8CnCPointCache) SetValue(Slot, K string, V Hnc8RegisterPoint) {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	}
}
func (M *Hnc8CnCPointCache) GetValue(Slot, K string) Hnc8RegisterPoint {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S[K]
	}
	return Hnc8RegisterPoint{}
}
func (M *Hnc8CnCPointCache) DeleteValue(Slot, K string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		delete(S, Slot)
	}
}
func (M *Hnc8CnCPointCache) UnRegisterSlot(Slot string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	delete(M.Slots, Slot)
	M.Flush()
}
func (M *Hnc8CnCPointCache) Size() uint64 {
	return uint64(len(M.Slots))
}
func (M *Hnc8CnCPointCache) Flush() {
	for slotName, slot := range M.Slots {
		for k, _ := range slot {
			delete(slot, k)
		}
		delete(M.Slots, slotName)
	}
}

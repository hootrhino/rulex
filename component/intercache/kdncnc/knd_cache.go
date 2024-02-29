// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) KndRegisterPoint later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT KndRegisterPoint WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package kdncnc

import (
	"sync"

	"github.com/hootrhino/rulex/component/intercache"
	"github.com/hootrhino/rulex/typex"
)

// 点位表
type KndRegisterPoint struct {
	UUID          string
	Status        int
	LastFetchTime uint64
	Value         string
}

var __DefaultKdnCnCPointCache *KdnCnCPointCache

func RegisterSlot(Slot string) {
	__DefaultKdnCnCPointCache.RegisterSlot(Slot)
}
func GetSlot(Slot string) map[string]KndRegisterPoint {
	return __DefaultKdnCnCPointCache.GetSlot(Slot)
}
func SetValue(Slot, K string, V KndRegisterPoint) {
	__DefaultKdnCnCPointCache.SetValue(Slot, K, V)
}
func GetValue(Slot, K string) KndRegisterPoint {
	return __DefaultKdnCnCPointCache.GetValue(Slot, K)
}
func DeleteValue(Slot, K string) {
	__DefaultKdnCnCPointCache.DeleteValue(Slot, K)
}
func UnRegisterSlot(Slot string) {
	__DefaultKdnCnCPointCache.UnRegisterSlot(Slot)
}
func Size() uint64 {
	return __DefaultKdnCnCPointCache.Size()
}
func Flush() {
	__DefaultKdnCnCPointCache.Flush()
}

//KdnCnC 点位运行时存储器

type KdnCnCPointCache struct {
	Slots      map[string]map[string]KndRegisterPoint
	ruleEngine typex.RuleX
	lock       sync.Mutex
}

func InitKdnCnCPointCache(ruleEngine typex.RuleX) intercache.InterCache {
	__DefaultKdnCnCPointCache = &KdnCnCPointCache{
		ruleEngine: ruleEngine,
		Slots:      map[string]map[string]KndRegisterPoint{},
		lock:       sync.Mutex{},
	}
	return __DefaultKdnCnCPointCache
}
func (M *KdnCnCPointCache) RegisterSlot(Slot string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	M.Slots[Slot] = map[string]KndRegisterPoint{}
}
func (M *KdnCnCPointCache) GetSlot(Slot string) map[string]KndRegisterPoint {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S
	}
	return nil
}
func (M *KdnCnCPointCache) SetValue(Slot, K string, V KndRegisterPoint) {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	}
}
func (M *KdnCnCPointCache) GetValue(Slot, K string) KndRegisterPoint {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S[K]
	}
	return KndRegisterPoint{}
}
func (M *KdnCnCPointCache) DeleteValue(Slot, K string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		delete(S, Slot)
	}
}
func (M *KdnCnCPointCache) UnRegisterSlot(Slot string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	delete(M.Slots, Slot)
	M.Flush()
}
func (M *KdnCnCPointCache) Size() uint64 {
	return uint64(len(M.Slots))
}
func (M *KdnCnCPointCache) Flush() {
	for slotName, slot := range M.Slots {
		for k, _ := range slot {
			delete(slot, k)
		}
		delete(M.Slots, slotName)
	}
}

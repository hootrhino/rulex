// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) SiemensPoint later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT SiemensPoint WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package siemens

import (
	"sync"

	"github.com/hootrhino/rulex/component/intercache"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultSiemensPointCache *SiemensPointCache

// 点位表
type SiemensPoint struct {
	UUID          string
	Status        int
	LastFetchTime uint64
	Value         string
}

func RegisterSlot(Slot string) {
	__DefaultSiemensPointCache.RegisterSlot(Slot)
}
func GetSlot(Slot string) map[string]SiemensPoint {
	return __DefaultSiemensPointCache.GetSlot(Slot)
}
func SetValue(Slot, K string, V SiemensPoint) {
	__DefaultSiemensPointCache.SetValue(Slot, K, V)
}
func GetValue(Slot, K string) SiemensPoint {
	return __DefaultSiemensPointCache.GetValue(Slot, K)
}
func DeleteValue(Slot, K string) {
	__DefaultSiemensPointCache.DeleteValue(Slot, K)
}
func UnRegisterSlot(Slot string) {
	__DefaultSiemensPointCache.UnRegisterSlot(Slot)
}
func Size() uint64 {
	return __DefaultSiemensPointCache.Size()
}
func Flush() {
	__DefaultSiemensPointCache.Flush()
}

type SiemensPointCache struct {
	Slots      map[string]map[string]SiemensPoint
	ruleEngine typex.RuleX
	lock       sync.Mutex
}

func InitSiemensPointCache(ruleEngine typex.RuleX) intercache.InterCache {
	__DefaultSiemensPointCache = &SiemensPointCache{
		ruleEngine: ruleEngine,
		Slots:      map[string]map[string]SiemensPoint{},
		lock:       sync.Mutex{},
	}
	return __DefaultSiemensPointCache
}
func (M *SiemensPointCache) RegisterSlot(Slot string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	M.Slots[Slot] = map[string]SiemensPoint{}
}
func (M *SiemensPointCache) GetSlot(Slot string) map[string]SiemensPoint {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S
	}
	return nil
}
func (M *SiemensPointCache) SetValue(Slot, K string, V SiemensPoint) {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	}
}
func (M *SiemensPointCache) GetValue(Slot, K string) SiemensPoint {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S[K]
	}
	return SiemensPoint{}
}
func (M *SiemensPointCache) DeleteValue(Slot, K string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		delete(S, Slot)
	}
}
func (M *SiemensPointCache) UnRegisterSlot(Slot string) {
	M.lock.Lock()
	defer M.lock.Unlock()
	delete(M.Slots, Slot)
}
func (M *SiemensPointCache) Size() uint64 {
	return uint64(len(M.Slots))
}
func (M *SiemensPointCache) Flush() {
	for slotName, slot := range M.Slots {
		for k, _ := range slot {
			delete(slot, k)
		}
		delete(M.Slots, slotName)
	}
}

package statistics

import "sync"

type statistics struct {
	InSuccess  int64 `json:"inSuccess"`
	OutSuccess int64 `json:"outSuccess"`
	InFailed   int64 `json:"inFailed"`
	OutFailed  int64 `json:"outFailed"`
}

var statisticsCache statistics

var lock sync.Mutex

func init() {
	statisticsCache = statistics{}
}
func IncIn() {

	statisticsCache.InSuccess = statisticsCache.InSuccess + 1
}
func DecIn() {
	lock.Lock()
	defer lock.Unlock()
	if statisticsCache.InSuccess-1 > 0 {
		statisticsCache.InSuccess = statisticsCache.InSuccess - 1
	}
}
func IncOut() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache.OutSuccess = statisticsCache.OutSuccess + 1
}
func DecOut() {
	lock.Lock()
	defer lock.Unlock()

	if statisticsCache.OutSuccess-1 > 0 {
		statisticsCache.OutSuccess = statisticsCache.OutSuccess - 1
	}
}
func IncInFailed() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache.InFailed = statisticsCache.InFailed + 1
}

func IncOutFailed() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache.OutFailed = statisticsCache.OutFailed + 1
}

func Reset() {
	statisticsCache.InSuccess = 0
	statisticsCache.InFailed = 0
	statisticsCache.OutFailed = 0
	statisticsCache.OutSuccess = 0
}
func AllStatistics() statistics {
	return statisticsCache
}

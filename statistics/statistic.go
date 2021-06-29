package statistics

import "sync"

// import "time"

var statisticsCache map[string]int64

// var recentlyLogCache []string
var lock sync.Mutex

func init() {
	// recentlyLogCache = make([]string, 50)
	statisticsCache = map[string]int64{
		"inSuccess":  0,
		"outSuccess": 0,
		"inFailed":   0,
		"outFailed":  0,
	}
}
func IncIn() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache["inSuccess"] = statisticsCache["inSuccess"] + 1
}
func DecIn() {
	lock.Lock()
	if statisticsCache["inSuccess"]-1 > 0 {
		statisticsCache["inSuccess"] = statisticsCache["inSuccess"] - 1
	}
}
func IncOut() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache["outSuccess"] = statisticsCache["outSuccess"] + 1
}
func DecOut() {
	lock.Lock()
	defer lock.Unlock()

	if statisticsCache["outSuccess"]-1 > 0 {
		statisticsCache["outSuccess"] = statisticsCache["outSuccess"] - 1
	}
}
func IncInFailed() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache["inFailed"] = statisticsCache["inFailed"] + 1
}

func IncOutFailed() {
	lock.Lock()
	defer lock.Unlock()
	statisticsCache["outFailed"] = statisticsCache["outFailed"] + 1
}

func AllStatistics() *map[string]int64 {
	return &statisticsCache
}

// func AddLog(logs string) {
// 	if len(recentlyLogCache) > 50 {
// 		recentlyLogCache = recentlyLogCache[1:]
// 		recentlyLogCache = append(recentlyLogCache, time.Now().Format("2006.01.02 15:04:05")+" => "+logs)

// 	} else {
// 		recentlyLogCache = append(recentlyLogCache, time.Now().Format("2006.01.02 15:04:05")+" => "+logs)
// 	}
// }
// func AllLog() []string {
// 	return recentlyLogCache
// }

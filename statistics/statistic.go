package statistics

// import "time"

var statisticsCache map[string]int64
// var recentlyLogCache []string

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
	statisticsCache["inSuccess"] = statisticsCache["inSuccess"] + 1
}
func DecIn() {
	if statisticsCache["inSuccess"]-1 > 0 {
		statisticsCache["inSuccess"] = statisticsCache["inSuccess"] - 1
	}
}
func IncOut() {
	statisticsCache["outSuccess"] = statisticsCache["outSuccess"] + 1

}
func DecOut() {
	if statisticsCache["outSuccess"]-1 > 0 {
		statisticsCache["outSuccess"] = statisticsCache["outSuccess"] - 1
	}
}
func IncInFailed() {
	statisticsCache["inFailed"] = statisticsCache["inFailed"] + 1
}

func IncOutFailed() {
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

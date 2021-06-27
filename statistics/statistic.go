package statistics

var cache map[string]int64

func init() {
	cache = map[string]int64{
		"in":     0,
		"out":    0,
		"failed": 0,
	}
}
func IncIn() {
	cache["in"] = cache["in"] + 1
}
func DecIn() {
	if cache["in"]-1 > 0 {
		cache["in"] = cache["in"] - 1
	}
}
func IncOut() {
	cache["out"] = cache["out"] + 1

}
func DecOut() {
	if cache["out"]-1 > 0 {
		cache["out"] = cache["out"] - 1
	}
}
func IncFailed() {
	cache["failed"] = cache["failed"] + 1
}
func DecFailed() {
	if cache["failed"]-1 > 0 {
		cache["failed"] = cache["failed"] - 1
	}
}

func AllStatistics() *map[string]int64 {
	return &cache
}

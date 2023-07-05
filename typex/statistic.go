package typex

import (
	"math"
)

type MetricStatistics struct {
	InSuccess  uint64 `json:"inSuccess"`
	OutSuccess uint64 `json:"outSuccess"`
	InFailed   uint64 `json:"inFailed"`
	OutFailed  uint64 `json:"outFailed"`
}

func NewMetricStatistics() *MetricStatistics {
	return &MetricStatistics{
		InSuccess:  0,
		OutSuccess: 0,
		InFailed:   0,
		OutFailed:  0,
	}

}
func (statisticsCache *MetricStatistics) IncIn() {
	if statisticsCache.InSuccess < math.MaxUint64 {
		statisticsCache.InSuccess = statisticsCache.InSuccess + 1
	}
}
func (statisticsCache *MetricStatistics) DecIn() {

	if statisticsCache.InSuccess-1 > 0 {
		statisticsCache.InSuccess = statisticsCache.InSuccess - 1
	}
}
func (statisticsCache *MetricStatistics) IncOut() {

	if statisticsCache.OutSuccess < math.MaxUint64 {
		statisticsCache.OutSuccess = statisticsCache.OutSuccess + 1
	}
}
func (statisticsCache *MetricStatistics) DecOut() {

	if statisticsCache.OutSuccess-1 > 0 {
		statisticsCache.OutSuccess = statisticsCache.OutSuccess - 1
	}
}
func (statisticsCache *MetricStatistics) IncInFailed() {

	if statisticsCache.InFailed < math.MaxUint64 {
		statisticsCache.InFailed = statisticsCache.InFailed + 1
	}
}

func (statisticsCache *MetricStatistics) IncOutFailed() {

	if statisticsCache.InFailed < math.MaxUint64 {
		statisticsCache.OutFailed = statisticsCache.OutFailed + 1
	}
}

func (statisticsCache *MetricStatistics) Reset() {
	statisticsCache.InSuccess = 0
	statisticsCache.InFailed = 0
	statisticsCache.OutFailed = 0
	statisticsCache.OutSuccess = 0
}

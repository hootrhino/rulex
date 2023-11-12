package intermetric

import (
	"math"

	"github.com/hootrhino/rulex/typex"
)

var __DefaultInternalMetric *InternalMetric

type InternalMetric struct {
	rulex      typex.RuleX
	InSuccess  uint64 `json:"inSuccess"`
	OutSuccess uint64 `json:"outSuccess"`
	InFailed   uint64 `json:"inFailed"`
	OutFailed  uint64 `json:"outFailed"`
}

func InitInternalMetric(rulex typex.RuleX) *InternalMetric {
	__DefaultInternalMetric = &InternalMetric{
		InSuccess:  0,
		OutSuccess: 0,
		InFailed:   0,
		OutFailed:  0,
	}
	__DefaultInternalMetric.rulex = rulex
	return __DefaultInternalMetric

}
func GetMetric() InternalMetric {
	return *__DefaultInternalMetric

}
func IncIn() {
	if __DefaultInternalMetric.InSuccess < math.MaxUint64 {
		__DefaultInternalMetric.InSuccess = __DefaultInternalMetric.InSuccess + 1
	}
}
func DecIn() {

	if __DefaultInternalMetric.InSuccess-1 > 0 {
		__DefaultInternalMetric.InSuccess = __DefaultInternalMetric.InSuccess - 1
	}
}
func IncOut() {

	if __DefaultInternalMetric.OutSuccess < math.MaxUint64 {
		__DefaultInternalMetric.OutSuccess = __DefaultInternalMetric.OutSuccess + 1
	}
}
func DecOut() {

	if __DefaultInternalMetric.OutSuccess-1 > 0 {
		__DefaultInternalMetric.OutSuccess = __DefaultInternalMetric.OutSuccess - 1
	}
}
func IncInFailed() {

	if __DefaultInternalMetric.InFailed < math.MaxUint64 {
		__DefaultInternalMetric.InFailed = __DefaultInternalMetric.InFailed + 1
	}
}

func IncOutFailed() {

	if __DefaultInternalMetric.InFailed < math.MaxUint64 {
		__DefaultInternalMetric.OutFailed = __DefaultInternalMetric.OutFailed + 1
	}
}

func Reset() {
	__DefaultInternalMetric.InSuccess = 0
	__DefaultInternalMetric.InFailed = 0
	__DefaultInternalMetric.OutFailed = 0
	__DefaultInternalMetric.OutSuccess = 0
}

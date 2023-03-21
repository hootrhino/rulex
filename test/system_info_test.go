package test

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/v3/mem"
)

// {
//     "total":34061004800,
//     "available":21563211776,
//     "used":12497793024,
//     "usedPercent":36,
//     "free":21563211776,
//     "active":0,
//     "inactive":0,
//     "wired":0,
//     "laundry":0,
//     "buffers":0,
//     "cached":0,
//     "writeBack":0,
//     "dirty":0,
//     "writeBackTmp":0,
//     "shared":0,
//     "slab":0,
//     "sreclaimable":0,
//     "sunreclaim":0,
//     "pageTables":0,
//     "swapCached":0,
//     "commitLimit":0,
//     "committedAS":0,
//     "highTotal":0,
//     "highFree":0,
//     "lowTotal":0,
//     "lowFree":0,
//     "swapTotal":0,
//     "swapFree":0,
//     "mapped":0,
//     "vmallocTotal":0,
//     "vmallocUsed":0,
//     "vmallocChunk":0,
//     "hugePagesTotal":0,
//     "hugePagesFree":0,
//     "hugePageSize":0
// }

func Test_gopsutil(t *testing.T) {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)
}

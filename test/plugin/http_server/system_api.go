package httpserver

import (
	"net/http"
	"rulex/statistics"
	"rulex/typex"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

//
// Get all plugins
//
func Plugins(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	plugins := e.AllPlugins()
	plugins.Range(func(key, value interface{}) bool {
		data = append(data, value.(typex.XPlugin).XPluginMetaInfo())
		return true
	})
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	})
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//
// Get system infomation
//
func System(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	cpuPercent, _ := cpu.Percent(time.Millisecond, true)
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: gin.H{
			"version":    e.Version().Version,
			"diskInfo":   int(diskInfo.UsedPercent),
			"system":     bToMb(m.Sys),
			"alloc":      bToMb(m.Alloc),
			"total":      bToMb(m.TotalAlloc),
			"cpuPercent": cpuPercent,
			"osArch":     runtime.GOOS + "-" + runtime.GOARCH,
		},
	})
}

//
// Get all Drivers
//
func Drivers(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	inEnds := e.AllInEnd()
	inEnds.Range(func(key, value interface{}) bool {
		if value.(*typex.InEnd).Resource.Driver() != nil {
			data = append(data, value.(*typex.InEnd).Resource.Driver().DriverDetail())
		}
		return true
	})
	c.JSON(200, Result{
		Code: 200,
		Msg:  "Success",
		Data: data,
	})
}

//
// Get statistics data
//
func Statistics(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: statistics.AllStatistics(),
	})
}

//
// Get statistics data
//
func ResourceCount(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: map[string]int{
			"inends":  0,
			"outends": 0,
			"rules":   0,
			"plugins": 0,
		},
	})
}

package httpserver

import (
	"runtime"
	"time"

	"github.com/i4de/rulex/source"
	"github.com/i4de/rulex/statistics"
	"github.com/i4de/rulex/target"
	"github.com/i4de/rulex/typex"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"go.bug.st/serial"
)

/*
*
* 健康检查接口, 一般用来监视是否工作
*
 */
func Ping(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.Writer.Write([]byte("PONG"))
	c.Writer.Flush()
}

// Get all plugins
func Plugins(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	plugins := e.AllPlugins()
	plugins.Range(func(key, value interface{}) bool {
		data = append(data, value.(typex.XPlugin).PluginMetaInfo())
		return true
	})
	c.JSON(200, OkWithData(data))
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Get system infomation
func System(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	cpuPercent, _ := cpu.Percent(5*time.Millisecond, true)
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.JSON(200, OkWithData(gin.H{
		"version":     e.Version().Version,
		"diskInfo":    int(diskInfo.UsedPercent),
		"system":      bToMb(m.Sys),
		"alloc":       bToMb(m.Alloc),
		"total":       bToMb(m.TotalAlloc),
		"cpuPercent":  calculateCpuPercent(cpuPercent),
		"osArch":      runtime.GOOS + "-" + runtime.GOARCH,
		"startedTime": StartedTime,
	}))
}

// Get all Drivers
func Drivers(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	e.AllInEnd().Range(func(key, value interface{}) bool {
		if value.(*typex.InEnd).Source.Driver() != nil {
			data = append(data, value.(*typex.InEnd).Source.Driver().DriverDetail())
		}
		return true
	})
	e.AllDevices().Range(func(key, value interface{}) bool {
		if value.(*typex.Device).Device.Driver() != nil {
			data = append(data, value.(*typex.Device).Device.Driver().DriverDetail())
		}
		return true
	})
	c.JSON(200, OkWithData(data))
}

// Get statistics data
func Statistics(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData(statistics.AllStatistics()))
}

// Get statistics data
func SourceCount(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	allInEnd := e.AllInEnd()
	allOutEnd := e.AllOutEnd()
	allRule := e.AllRule()
	plugins := e.AllPlugins()
	var c1, c2, c3, c4 int
	allInEnd.Range(func(key, value interface{}) bool {
		c1 += 1
		return true
	})
	allOutEnd.Range(func(key, value interface{}) bool {
		c2 += 1
		return true
	})
	allRule.Range(func(key, value interface{}) bool {
		c3 += 1
		return true
	})
	plugins.Range(func(key, value interface{}) bool {
		c4 += 1
		return true
	})
	c.JSON(200, OkWithData(map[string]int{
		"inends":  c1,
		"outends": c2,
		"rules":   c3,
		"plugins": c4,
	}))
}

/*
*
* 输入类型配置
*
 */
func RType(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	Type, _ := c.GetQuery("type")
	if Type == "" {
		c.JSON(200, OkWithData(source.SM.All()))
	} else {
		c.JSON(200, OkWithData(source.SM.Find(typex.InEndType(Type))))
	}

}

/*
*
* 输出类型配置
*
 */
func TType(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	Type, _ := c.GetQuery("type")
	if Type == "" {
		c.JSON(200, OkWithData(target.TM.All()))
	} else {
		c.JSON(200, OkWithData(target.TM.Find(typex.TargetType(Type))))
	}

}

/*
*
* 获取本地的串口列表
*
 */
func GetUarts(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	ports, _ := serial.GetPortsList()
	c.JSON(200, OkWithData(ports))
}

/*
*
* 计算开机时间
*
 */
func StartedAt(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData(StartedTime))
}

// 计算CPU平均使用率
func calculateCpuPercent(cpus []float64) float64 {
	var acc float64 = 0
	for _, v := range cpus {
		acc += v
	}
	return acc / float64(len(cpus))
}

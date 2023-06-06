package httpserver

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/statistics"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

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
		pi := value.(typex.XPlugin).PluginMetaInfo()
		data = append(data, pi)
		return true
	})
	c.JSON(200, OkWithData(data))
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// 计算资源数据
func source_count(e typex.RuleX) map[string]int {
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
	return map[string]int{
		"inends":  c1,
		"outends": c2,
		"rules":   c3,
		"plugins": c4,
	}
}

/*
*
* 获取系统指标, Go 自带这个不准, 后期版本需要更换跨平台实现
*
 */
func System(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	cpuPercent, _ := cpu.Percent(5*time.Millisecond, false)
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ip, err0 := utils.HostNameI()
	hardWareInfo := map[string]interface{}{
		"version":    e.Version().Version,
		"diskInfo":   calculateDiskInfo(diskInfo),
		"systemMem":  bToMb(m.Sys),
		"allocMem":   bToMb(m.Alloc),
		"totalMem":   bToMb(m.TotalAlloc),
		"cpuPercent": calculateCpuPercent(cpuPercent),
		"osArch":     runtime.GOOS + "-" + runtime.GOARCH,
		"osDist": func() string {
			v, _ := utils.GetOSDistribution()
			return v
		}(),
		"startedTime": StartedTime,
		"ip": func() []string {
			if err0 != nil {
				glogger.GLogger.Error(err0)
				return []string{"127.0.0.1"}
			}
			return ip
		}(),
		"wsUrl": func() []string {
			if err0 != nil {
				glogger.GLogger.Error(err0)
				return []string{"ws://127.0.0.1:2580/ws"}
			}
			ips := []string{}
			for _, ipp := range ip {
				ips = append(ips, fmt.Sprintf("ws://%s:2580/ws", ipp))
			}
			return ips
		}(),
	}
	c.JSON(200, OkWithData(gin.H{
		"hardWareInfo": hardWareInfo,
		"statistic":    statistics.AllStatistics(),
		"sourceCount":  source_count(e),
	}))
}

/*
*
* SnapshotDump
*
 */
func SnapshotDump(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	c.JSON(200, OkWithData(e.SnapshotDump()))
}

// Get all Drivers
func Drivers(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	data := []interface{}{}
	id := 0
	e.AllInEnd().Range(func(key, value interface{}) bool {
		drivers := value.(*typex.InEnd).Source.Driver()
		if drivers != nil {
			dd := drivers.DriverDetail()
			dd.UUID = fmt.Sprintf("DRIVER:%v", id)
			id++
			data = append(data, dd)
		}
		return true
	})
	e.AllDevices().Range(func(key, value interface{}) bool {
		drivers := value.(*typex.Device).Device.Driver()
		if drivers != nil {
			dd := drivers.DriverDetail()
			dd.UUID = fmt.Sprintf("DRIVER:%v", id)
			id++
			data = append(data, dd)
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
* 设备配置
*
 */
func DType(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	Type, _ := c.GetQuery("type")
	if Type == "" {
		c.JSON(200, OkWithData(device.DM.All()))
	} else {
		c.JSON(200, OkWithData(device.DM.Find(typex.DeviceType(Type))))
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

func calculateDiskInfo(diskInfo *disk.UsageStat) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", diskInfo.UsedPercent), 64)
	return value

}

// 计算CPU平均使用率
func calculateCpuPercent(cpus []float64) float64 {
	var acc float64 = 0
	for _, v := range cpus {
		acc += v
	}
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", acc), 64)
	return value
}

package httpserver

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"

	common "github.com/hootrhino/rulex/plugin/http_server/common"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/typex"

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
func Ping(c *gin.Context, hh *HttpApiServer) {
	c.Writer.Write([]byte("PONG"))
	c.Writer.Flush()
}

// Get all plugins
func Plugins(c *gin.Context, hh *HttpApiServer) {
	data := []interface{}{}
	plugins := hh.ruleEngine.AllPlugins()
	plugins.Range(func(key, value interface{}) bool {
		pi := value.(typex.XPlugin).PluginMetaInfo()
		data = append(data, pi)
		return true
	})
	c.JSON(common.HTTP_OK, common.OkWithData(data))
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
func System(c *gin.Context, hh *HttpApiServer) {
	cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, true)
	diskInfo, _ := disk.Usage("/")
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// ip, err0 := utils.HostNameI()
	hardWareInfo := map[string]interface{}{
		"version":     hh.ruleEngine.Version().Version,
		"diskInfo":    calculateDiskInfo(diskInfo),
		"systemMem":   bToMb(m.Sys),
		"allocMem":    bToMb(m.Alloc),
		"totalMem":    bToMb(m.TotalAlloc),
		"cpuPercent":  calculateCpuPercent(cpuPercent),
		"osArch":      hh.ruleEngine.Version().Arch,
		"osDist":      hh.ruleEngine.Version().Dist,
		"startedTime": StartedTime,
		// "ip": func() []string {
		// 	if err0 != nil {
		// 		glogger.GLogger. common.Error(err0)
		// 		return []string{"127.0.0.1"}
		// 	}
		// 	return ip
		// }(),
		// "wsUrl": func() []string {
		// 	if err0 != nil {
		// 		glogger.GLogger. common.Error(err0)
		// 		return []string{"ws://127.0.0.1:2580/ws"}
		// 	}
		// 	ips := []string{}
		// 	for _, ipp := range ip {
		// 		ips = append(ips, fmt.Sprintf("ws://%s:2580/ws", ipp))
		// 	}
		// 	return ips
		// }(),
	}
	c.JSON(common.HTTP_OK, common.OkWithData(gin.H{
		"hardWareInfo": hardWareInfo,
		"statistic":    hh.ruleEngine.GetMetricStatistics(),
		"sourceCount":  source_count(hh.ruleEngine),
	}))
}

/*
*
* SnapshotDump
*
 */
func SnapshotDump(c *gin.Context, hh *HttpApiServer) {
	c.JSON(common.HTTP_OK, common.OkWithData(hh.ruleEngine.SnapshotDump()))
}

// Get all Drivers
func Drivers(c *gin.Context, hh *HttpApiServer) {
	data := []interface{}{}
	id := 0
	hh.ruleEngine.AllInEnd().Range(func(key, value interface{}) bool {
		drivers := value.(*typex.InEnd).Source.Driver()
		if drivers != nil {
			dd := drivers.DriverDetail()
			dd.UUID = fmt.Sprintf("DRIVER:%v", id)
			id++
			data = append(data, dd)
		}
		return true
	})
	hh.ruleEngine.AllDevices().Range(func(key, value interface{}) bool {
		drivers := value.(*typex.Device).Device.Driver()
		if drivers != nil {
			dd := drivers.DriverDetail()
			dd.UUID = fmt.Sprintf("DRIVER:%v", id)
			id++
			data = append(data, dd)
		}
		return true
	})
	c.JSON(common.HTTP_OK, common.OkWithData(data))
}

// Get statistics data
func Statistics(c *gin.Context, hh *HttpApiServer) {
	c.JSON(common.HTTP_OK, common.OkWithData(hh.ruleEngine.GetMetricStatistics()))
}

// Get statistics data
func SourceCount(c *gin.Context, hh *HttpApiServer) {
	allInEnd := hh.ruleEngine.AllInEnd()
	allOutEnd := hh.ruleEngine.AllOutEnd()
	allRule := hh.ruleEngine.AllRule()
	plugins := hh.ruleEngine.AllPlugins()
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
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]int{
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
func RType(c *gin.Context, hh *HttpApiServer) {
	Type, _ := c.GetQuery("type")
	if Type == "" {
		c.JSON(common.HTTP_OK, common.OkWithData(source.SM.All()))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(source.SM.Find(typex.InEndType(Type))))
	}

}

/*
*
* 输出类型配置
*
 */
func TType(c *gin.Context, hh *HttpApiServer) {
	Type, _ := c.GetQuery("type")
	if Type == "" {
		c.JSON(common.HTTP_OK, common.OkWithData(target.TM.All()))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(target.TM.Find(typex.TargetType(Type))))
	}

}

/*
*
* 设备配置
*
 */
func DType(c *gin.Context, hh *HttpApiServer) {
	Type, _ := c.GetQuery("type")
	if Type == "" {
		c.JSON(common.HTTP_OK, common.OkWithData(device.DM.All()))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(device.DM.Find(typex.DeviceType(Type))))
	}

}

/*
*
* 获取本地的串口列表
*
 */
func GetUarts(c *gin.Context, hh *HttpApiServer) {
	ports, _ := serial.GetPortsList()
	c.JSON(common.HTTP_OK, common.OkWithData(ports))
}

func GetNetInterfaces(c *gin.Context, hh *HttpApiServer) {
	interfaces, err := getAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(interfaces))
	}
}

/*
*
* 计算开机时间
*
 */
func StartedAt(c *gin.Context, hh *HttpApiServer) {
	c.JSON(common.HTTP_OK, common.OkWithData(StartedTime))
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
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", acc/float64(len(cpus))), 64)
	return value
}

type NetInterfaceInfo struct {
	Name string `json:"name,omitempty"`
	Mac  string `json:"mac,omitempty"`
	Addr string `json:"addr,omitempty"`
}

func getAvailableInterfaces() ([]NetInterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netInterfaces := make([]NetInterfaceInfo, 0, len(interfaces))
	for _, inter := range interfaces {
		info := NetInterfaceInfo{
			Name: inter.Name,
			Mac:  inter.HardwareAddr.String(),
		}
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}
		for i := range addrs {
			addr := addrs[i].String()
			cidr, _, _ := net.ParseCIDR(addr)
			if cidr == nil {
				continue
			}
			if cidr.To4() != nil {
				// 找到第一个ipv4地址
				info.Addr = addr
				break
			}
		}
		netInterfaces = append(netInterfaces, info)

	}

	return netInterfaces, nil
}

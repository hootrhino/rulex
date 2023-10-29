package apis

import (
	"fmt"
	"net"

	// "runtime"
	"strconv"
	"time"

	"github.com/hootrhino/rulex/component/appstack"
	"github.com/hootrhino/rulex/component/intermetric"
	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/ossupport"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/utils"

	"github.com/hootrhino/rulex/device"
	"github.com/hootrhino/rulex/source"
	"github.com/hootrhino/rulex/target"
	"github.com/hootrhino/rulex/typex"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"go.bug.st/serial"
)

// 启动时间
var __StartedAt = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

/*
*
* 健康检查接口, 一般用来监视是否工作
*
 */
func Ping(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData("PONG"))
}

// Get all plugins
func Plugins(c *gin.Context, ruleEngine typex.RuleX) {
	data := []interface{}{}
	plugins := ruleEngine.AllPlugins()
	plugins.Range(func(key, value interface{}) bool {
		pi := value.(typex.XPlugin).PluginMetaInfo()
		data = append(data, pi)
		return true
	})
	c.JSON(common.HTTP_OK, common.OkWithData(data))
}

// 计算资源数据
func source_count(e typex.RuleX) map[string]int {
	allInEnd := e.AllInEnd()
	allOutEnd := e.AllOutEnd()
	allRule := e.AllRule()
	plugins := e.AllPlugins()
	devices := e.AllDevices()
	goods := trailer.AllGoods()
	var c1, c2, c3, c4, c5, c6 int
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
	devices.Range(func(key, value interface{}) bool {
		c5 += 1
		return true
	})
	goods.Range(func(key, value interface{}) bool {
		c6 += 1
		return true
	})
	return map[string]int{
		"inends":  c1,
		"outends": c2,
		"rules":   c3,
		"plugins": c4,
		"devices": c5,
		"goods":   c6,
		"apps":    appstack.AppCount(),
	}
}

/*
*
* 获取系统指标, Go 自带这个不准, 后期版本需要更换跨平台实现
*
 */
func System(c *gin.Context, ruleEngine typex.RuleX) {
	cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, true)
	diskInfo, _ := disk.Usage("/")
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// ip, err0 := utils.HostNameI()
	memPercent, _ := service.GetMemPercent()
	hardWareInfo := map[string]interface{}{
		"version":    ruleEngine.Version().Version,
		"diskInfo":   calculateDiskInfo(diskInfo),
		"memPercent": memPercent,
		"cpuPercent": calculateCpuPercent(cpuPercent),
		"osArch":     ruleEngine.Version().Arch,
		"osDist":     ruleEngine.Version().Dist,
		"startedAt":  __StartedAt,
		"osUpTime": func() string {
			result, err := ossupport.GetUptime()
			if err != nil {
				return "0 days 0 hours 0 minutes"
			}
			return result
		}(),
	}
	c.JSON(common.HTTP_OK, common.OkWithData(gin.H{
		"hardWareInfo": hardWareInfo,
		"statistic":    intermetric.GetMetric(),
		"sourceCount":  source_count(ruleEngine),
	}))
}

/*
*
* SnapshotDump
*
 */
func SnapshotDump(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData(ruleEngine.SnapshotDump()))
}

// Get all Drivers
func Drivers(c *gin.Context, ruleEngine typex.RuleX) {
	data := []interface{}{}
	id := 0
	ruleEngine.AllInEnd().Range(func(key, value interface{}) bool {
		drivers := value.(*typex.InEnd).Source.Driver()
		if drivers != nil {
			dd := drivers.DriverDetail()
			dd.UUID = fmt.Sprintf("DRIVER:%v", id)
			id++
			data = append(data, dd)
		}
		return true
	})
	ruleEngine.AllDevices().Range(func(key, value interface{}) bool {
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
func Statistics(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData(intermetric.GetMetric()))
}

// Get statistics data
func SourceCount(c *gin.Context, ruleEngine typex.RuleX) {
	allInEnd := ruleEngine.AllInEnd()
	allOutEnd := ruleEngine.AllOutEnd()
	allRule := ruleEngine.AllRule()
	plugins := ruleEngine.AllPlugins()
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
func RType(c *gin.Context, ruleEngine typex.RuleX) {
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
func TType(c *gin.Context, ruleEngine typex.RuleX) {
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
func DType(c *gin.Context, ruleEngine typex.RuleX) {
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
func GetUarts(c *gin.Context, ruleEngine typex.RuleX) {
	ports, _ := serial.GetPortsList()
	c.JSON(common.HTTP_OK, common.OkWithData(ports))
}

/*
*
* apiV2
*
 */
func GetUartList(c *gin.Context, ruleEngine typex.RuleX) {
	ports, _ := serial.GetPortsList()
	List := []map[string]interface{}{}
	for _, port := range ports {
		// 明确知道有这么多端口，啰嗦代码是为了标记以免以后忘记
		if port == "/dev/ttyS0" {
			continue
		}
		if port == "/dev/ttyS3" {
			continue
		}
		if port == "/dev/ttyUSB0" {
			continue
		}
		if port == "/dev/ttyUSB1" {
			continue
		}
		if port == "/dev/ttyUSB2" {
			continue
		}

		List = append(List, map[string]interface{}{
			"port": port,
			"alias": func(port string) string {
				if port == "/dev/ttyS1" {
					return "RS4851(A1B1)"
				}
				if port == "/dev/ttyS2" {
					return "RS4851(A2B2)"
				}
				return port
			}(port)})

	}
	c.JSON(common.HTTP_OK, common.OkWithData(List))
}

/*
*
* 本地网卡
*
 */
func GetNetInterfaces(c *gin.Context, ruleEngine typex.RuleX) {
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
func StartedAt(c *gin.Context, ruleEngine typex.RuleX) {
	c.JSON(common.HTTP_OK, common.OkWithData(__StartedAt))
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
func CatOsRelease(c *gin.Context, ruleEngine typex.RuleX) {
	r, err := utils.CatOsRelease()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(r))
}

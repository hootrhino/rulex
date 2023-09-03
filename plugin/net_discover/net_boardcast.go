package netdiscover

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"gopkg.in/ini.v1"
)

var startedTime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

type _serverConfig struct {
	Enable bool   `ini:"enable"`
	Host   string `ini:"listen_host"`
	Port   int    `ini:"listen_port"`
}

/*
*
* 网关之间的发现协议,让同一个网络内的网关能发现对方.这个功能暂时没想好用来做什么.
*
 */
//
// udp://hostname@host:port
//
type gwnode struct {
	Host     string
	Hostname string
	Port     string
	Timeout  int32
}
type NetDiscover struct {
	mainConfig _serverConfig
	ctx        context.Context
	cancel     context.CancelFunc
	Neighbors  map[string]gwnode // 邻居
	uuid       string
}

func NewNetDiscover() typex.XPlugin {
	ctx, cancel := context.WithCancel(context.Background())
	return &NetDiscover{mainConfig: _serverConfig{},
		ctx: ctx, cancel: cancel, uuid: "NWT_DISCOVER"}
}

func (dm *NetDiscover) Init(config *ini.Section) error {

	if err := utils.InIMapToStruct(config, &dm.mainConfig); err != nil {
		return err
	}
	return nil
}

func (dm *NetDiscover) Start(typex.RuleX) error {
	// 超时管理器
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			for _, nb := range dm.Neighbors {
				if nb.Timeout == 0 {
					// Delete node
					continue
				}
				atomic.AddInt32(&nb.Timeout, -1)
			}
			time.Sleep(5 * time.Second)

		}
	}(dm.ctx)
	go func(ctx context.Context) {
		udp_addr, err := net.ResolveUDPAddr("udp4",
			fmt.Sprintf("%s:%v", dm.mainConfig.Host, dm.mainConfig.Port))
		if err != nil {
			glogger.GLogger.Fatal(err)
		}
		glogger.GLogger.Info("start net_discover udp listener:",
			fmt.Sprintf("%s:%v", dm.mainConfig.Host, dm.mainConfig.Port))
		listener, err := net.ListenUDP("udp4", udp_addr)
		if err != nil {
			glogger.GLogger.Fatal(err)
		}
		data := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			n, remoteAddr, err := listener.ReadFrom(data)
			if err != nil {
				glogger.GLogger.Errorf("read remote: %s, err: %s", remoteAddr.String(), err)
				continue
			}
			// 请求组网的包 CAST:IP-Addr
			// CAST:192.168.001.001
			if n == 4 {
				// 任何一个网关加入集群都会广播
				if string(data[:n]) == ("CAST") {
					glogger.GLogger.Infof("Received CAST from:%s", remoteAddr.String())
				}
				// 集群内同步网关列表
				if string(data[:n]) == ("SYNC") {
					glogger.GLogger.Infof("Received SYNC from:%s", remoteAddr.String())
				}
			}
			// 获取本地信息
			if string(data[:n]) == "NODE_INFO" {
				glogger.GLogger.Infof("Received NODE_INFO from:%s", remoteAddr.String())
				cpuPercent, _ := cpu.Percent(5*time.Millisecond, true)
				parts, _ := disk.Partitions(true)
				diskInfo, _ := disk.Usage(parts[0].Mountpoint)
				// For info on each, see: https://golang.org/pkg/runtime/#MemStats
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				hardWareInfo := map[string]interface{}{
					"version":     typex.DefaultVersion,
					"diskInfo":    int(diskInfo.UsedPercent),
					"systemMem":   bToMb(m.Sys),
					"allocMem":    bToMb(m.Alloc),
					"totalMem":    bToMb(m.TotalAlloc),
					"cpuPercent":  calculateCpuPercent(cpuPercent),
					"osArch":      runtime.GOOS + "-" + runtime.GOARCH,
					"startedTime": startedTime,
				}
				b, _ := json.Marshal(hardWareInfo)
				listener.WriteTo(b, remoteAddr)
			}

		}
	}(dm.ctx)
	return nil
}
func (dm *NetDiscover) Stop() error {
	dm.cancel()
	return nil
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// 计算CPU平均使用率
func calculateCpuPercent(cpus []float64) float64 {
	var acc float64 = 0
	for _, v := range cpus {
		acc += v
	}
	return acc / float64(len(cpus))
}

func (hh *NetDiscover) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:     hh.uuid,
		Name:     "NetDiscover",
		Version:  "v0.0.1",
		Homepage: "https://hootrhino.github.io",
		HelpLink: "https://hootrhino.github.io",
		Author:   "wwhai",
		Email:    "cnwwhai@gmail.com",
		License:  "MIT",
	}
}

/*
*
* 服务调用接口
*
 */
func (cs *NetDiscover) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}

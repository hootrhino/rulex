package device

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

type IcmpSender struct {
	typex.XStatus
	status     typex.DeviceState
	mainConfig common.IpConfig
	RuleEngine typex.RuleX
}

// Example: 0x02 0x92 0xFF 0x98
/*
*
* 温湿度传感器
*
 */
func NewIcmpSender(e typex.RuleX) typex.XDevice {
	sender := new(IcmpSender)
	sender.RuleEngine = e
	sender.mainConfig = common.IpConfig{}
	return sender
}

//  初始化
func (sender *IcmpSender) Init(devId string, configMap map[string]interface{}) error {
	sender.PointId = devId
	if err := utils.BindSourceConfig(configMap, &sender.mainConfig); err != nil {
		return err
	}
	for _, ip := range sender.mainConfig.Hosts {
		if net.ParseIP(ip) == nil {
			return errors.New("invalid ip:" + ip)
		}
	}
	if !sender.mainConfig.AutoRequest {
		sender.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(sender.mainConfig.Frequency) * time.Second)
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					sender.status = typex.DEV_STOP
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			// 轮询IP地址 然后发送ICMP包
			for _, ip := range sender.mainConfig.Hosts {
				t, err := Pingq(ip, time.Duration(sender.mainConfig.Timeout))
				if err != nil {
					glogger.GLogger.Error(err)
					continue
				}
				datas, _ := json.Marshal(map[string]interface{}{
					"tag":   ip,
					"value": t.String(),
				})
				sender.RuleEngine.WorkDevice(sender.Details(), string(datas))
			}
		}
	}(typex.GCTX)
	return nil
}

// 启动
func (sender *IcmpSender) Start(cctx typex.CCTX) error {
	sender.Ctx = cctx.Ctx
	sender.CancelCTX = cctx.CancelCTX
	//
	return nil
}

// 从设备里面读数据出来
func (sender *IcmpSender) OnRead(cmd int, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
func (sender *IcmpSender) OnWrite(cmd int, _ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (sender *IcmpSender) Status() typex.DeviceState {
	return sender.status
}

// 停止设备
func (sender *IcmpSender) Stop() {

	sender.CancelCTX()
}

// 设备属性，是一系列属性描述
func (sender *IcmpSender) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (sender *IcmpSender) Details() *typex.Device {
	return sender.RuleEngine.GetDevice(sender.PointId)
}

// 状态
func (sender *IcmpSender) SetState(status typex.DeviceState) {
	sender.status = status

}

// 驱动
func (sender *IcmpSender) Driver() typex.XExternalDriver {
	return nil
}

func (sender *IcmpSender) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

//--------------------------------------------------------------------------------------------------
// private
//--------------------------------------------------------------------------------------------------
func Pingq(ip string, timeout time.Duration) (time.Duration, error) {
	const IcmpLen = 8
	msg := [32]byte{
		8, 0, 0, 0, 0, 13, 0, 37,
	}
	check := checkSum(msg[:IcmpLen])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)

	remoteAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}
	conn, err := net.DialIP("ip:icmp", nil, remoteAddr)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	if _, err := conn.Write(msg[:IcmpLen]); err != nil {
		return 0, err
	}
	conn.SetReadDeadline(time.Now().Add(timeout))
	_, err1 := conn.Read(msg[:])
	conn.SetReadDeadline(time.Time{})
	if err1 != nil {
		return 0, err1
	}
	return time.Since(start), nil
}

func checkSum(msg []byte) uint16 {
	sum := 0
	for n := 0; n < len(msg); n += 2 {
		sum += int(msg[n])<<8 + int(msg[n+1])
	}
	sum = (sum >> 16) + sum&0xffff
	sum += sum >> 16
	return uint16(^sum)
}

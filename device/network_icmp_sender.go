package device

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type _IcmpSenderCommonConfig struct {
	Timeout int `json:"timeout" validate:"required"`
	// // Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询"`
	// // Request Frequency, default 5 second
	Frequency int64 `json:"frequency" validate:"required" title:"采集频率"`
}
type _IcmpSenderConfig struct {
	CommonConfig _IcmpSenderCommonConfig `json:"commonConfig" validate:"required"`
	IcmpConfig   common.IcmpConfig       `json:"icmpConfig" validate:"required"`
}
type IcmpSender struct {
	typex.XStatus
	status     typex.DeviceState
	mainConfig _IcmpSenderConfig
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
	sender.mainConfig = _IcmpSenderConfig{}
	return sender
}

//  初始化
func (sender *IcmpSender) Init(devId string, configMap map[string]interface{}) error {
	sender.PointId = devId
	if err := utils.BindSourceConfig(configMap, &sender.mainConfig); err != nil {
		return err
	}
	for _, ip := range sender.mainConfig.IcmpConfig.Hosts {
		if net.ParseIP(ip) == nil {
			return errors.New("invalid ip:" + ip)
		}
	}
	if !sender.mainConfig.CommonConfig.AutoRequest {
		sender.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(sender.mainConfig.CommonConfig.Frequency) * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			// 轮询IP地址 然后发送ICMP包
			for _, ip := range sender.mainConfig.IcmpConfig.Hosts {
				t, err := pingQ(ip, time.Duration(sender.mainConfig.CommonConfig.Timeout))
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
			<-ticker.C

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
func (sender *IcmpSender) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
func (sender *IcmpSender) OnWrite(cmd []byte, _ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (sender *IcmpSender) Status() typex.DeviceState {
	return sender.status
}

// 停止设备
func (sender *IcmpSender) Stop() {
	sender.status = typex.DEV_DOWN
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
func (sender *IcmpSender) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

// --------------------------------------------------------------------------------------------------
// private
// --------------------------------------------------------------------------------------------------
func pingQ(ip string, timeout time.Duration) (time.Duration, error) {
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

// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package target

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type _TcpCommonConfig struct {
	DataMode   string `json:"dataMode" validate:"required"`   // RAW_STRING ; HEX_STRING
	AllowPing  *bool  `json:"allowPing" validate:"required"`  // 是否开启ping
	PingPacket string `json:"pingPacket" validate:"required"` // Ping 包内容, 必填16字符以内
}
type _TcpMainConfig struct {
	CommonConfig _TcpCommonConfig  `json:"commonConfig" validate:"required"`
	HostConfig   common.HostConfig `json:"hostConfig" validate:"required"`
}
type TTcpTarget struct {
	typex.XStatus
	client     *net.TCPConn
	mainConfig _TcpMainConfig
	status     typex.SourceState
}

/*
*
* TCP 透传模式
*
 */
func NewTTcpTarget(e typex.RuleX) typex.XTarget {
	ht := new(TTcpTarget)
	ht.RuleEngine = e
	ht.mainConfig = _TcpMainConfig{
		CommonConfig: _TcpCommonConfig{
			DataMode: "RAW_STRING",
			AllowPing: func() *bool {
				b := true
				return &b
			}(),
			PingPacket: "HR0001", //  默认每隔5秒发送PING包
		},
		HostConfig: common.HostConfig{
			Host:    "127.0.0.1",
			Port:    2585,
			Timeout: 3000,
		},
	}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *TTcpTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}
	return nil

}
func (ht *TTcpTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	var err error
	host := fmt.Sprintf("%s:%d", ht.mainConfig.HostConfig.Host, ht.mainConfig.HostConfig.Port)
	serverAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}
	ht.client, err = net.DialTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}, serverAddr)
	if err != nil {
		return err
	}
	if *ht.mainConfig.CommonConfig.AllowPing {
		go func(ht *TTcpTarget) {
			for {
				select {
				case <-ht.Ctx.Done():
					{
						return
					}
				default:
					{
					}
				}
				ht.client.SetReadDeadline(
					time.Now().Add((time.Duration(ht.mainConfig.HostConfig.Timeout) *
						time.Millisecond)),
				)
				_, err1 := ht.client.Write([]byte(ht.mainConfig.CommonConfig.PingPacket))
				ht.client.SetReadDeadline(time.Time{})
				if err1 != nil {
					glogger.GLogger.Error("TTcpTarget Ping Error:", err1)
					ht.status = typex.SOURCE_DOWN
					return
				}
				time.Sleep(5 * time.Second)
			}
		}(ht)
	}
	ht.status = typex.SOURCE_UP
	glogger.GLogger.Info("TTcpTarget  success connect to:", host)
	return nil
}

func (ht *TTcpTarget) Status() typex.SourceState {
	return ht.status
}

/*
*
* 透传模式：字符串和十六进制
*
 */
func (ht *TTcpTarget) To(data interface{}) (interface{}, error) {
	if ht.client != nil {
		switch s := data.(type) {
		case string:
			if ht.mainConfig.CommonConfig.DataMode == "RAW_STRING" {
				ht.client.SetReadDeadline(
					time.Now().Add((time.Duration(ht.mainConfig.HostConfig.Timeout) *
						time.Millisecond)),
				)
				_, err0 := ht.client.Write([]byte(s))
				ht.client.SetReadDeadline(time.Time{})
				if err0 != nil {
					return 0, err0
				}
			}
			if ht.mainConfig.CommonConfig.DataMode == "HEX_STRING" {
				dByte, err1 := hex.DecodeString(s)
				if err1 != nil {
					return 0, err1
				}
				ht.client.SetReadDeadline(
					time.Now().Add((time.Duration(ht.mainConfig.HostConfig.Timeout) *
						time.Millisecond)),
				)
				_, err0 := ht.client.Write(dByte)
				ht.client.SetReadDeadline(time.Time{})
				if err0 != nil {
					return 0, err0
				}
			}
			return len(s), nil
		default:
			return 0, fmt.Errorf("only support string format")
		}
	}
	return 0, fmt.Errorf("tcp already disconnected")

}

func (ht *TTcpTarget) Stop() {
	ht.status = typex.SOURCE_STOP
	ht.CancelCTX()
	if ht.client != nil {
		ht.client.Close()
	}
}
func (ht *TTcpTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}

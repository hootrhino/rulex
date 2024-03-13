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

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type _UdpMainConfig struct {
	AllowPing  *bool  `json:"allowPing"`
	DataMode   string `json:"dataMode"`
	Host       string `json:"host"`
	PingPacket string `json:"pingPacket"`
	Port       int    `json:"port"`
	Timeout    int    `json:"timeout"`
}

/*
*
* 数据推到UDP
*
 */
type UdpTarget struct {
	typex.XStatus
	mainConfig _UdpMainConfig
	status     typex.SourceState
}

func NewUdpTarget(e typex.RuleX) typex.XTarget {
	udpT := new(UdpTarget)
	udpT.RuleEngine = e
	udpT.mainConfig = _UdpMainConfig{
		DataMode:   "RAW_STRING",
		Timeout:    3000,
		PingPacket: "PING\r\n",
		AllowPing: func() *bool {
			b := true
			return &b
		}(),
	}
	udpT.status = typex.SOURCE_DOWN
	return udpT
}

func (udpT *UdpTarget) Init(outEndId string, configMap map[string]interface{}) error {
	udpT.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &udpT.mainConfig); err != nil {
		return err
	}
	return nil

}
func (udpT *UdpTarget) Start(cctx typex.CCTX) error {
	udpT.Ctx = cctx.Ctx
	udpT.CancelCTX = cctx.CancelCTX
	if *udpT.mainConfig.AllowPing {
		go func(ht *UdpTarget) {
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
				socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
					IP:   net.ParseIP(udpT.mainConfig.Host),
					Port: udpT.mainConfig.Port,
				})
				if err != nil {
					glogger.GLogger.Error(err)
					udpT.status = typex.SOURCE_DOWN
					return
				}
				socket.Close()
				time.Sleep(5 * time.Second)
			}
		}(udpT)
	}
	udpT.status = typex.SOURCE_UP
	glogger.GLogger.Info("UdpTarget started")
	return nil
}

func (udpT *UdpTarget) Status() typex.SourceState {
	if err := udpT.UdpStatus(fmt.Sprintf("%s:%d",
		udpT.mainConfig.Host, udpT.mainConfig.Port)); err != nil {
		return typex.SOURCE_DOWN
	}
	return udpT.status

}
func (udpT *UdpTarget) To(data interface{}) (interface{}, error) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(udpT.mainConfig.Host),
		Port: udpT.mainConfig.Port,
	})
	if err != nil {
		return 0, err
	}
	defer socket.Close()
	switch s := data.(type) {
	case string:
		if udpT.mainConfig.DataMode == "RAW_STRING" {
			socket.SetReadDeadline(
				time.Now().Add((time.Duration(udpT.mainConfig.Timeout) *
					time.Millisecond)),
			)
			_, err0 := socket.Write([]byte(s + "\r\n"))
			socket.SetReadDeadline(time.Time{})
			if err0 != nil {
				return 0, err0
			}
		}
		if udpT.mainConfig.DataMode == "HEX_STRING" {
			dByte, err1 := hex.DecodeString(s)
			if err1 != nil {
				return 0, err1
			}
			socket.SetReadDeadline(
				time.Now().Add((time.Duration(udpT.mainConfig.Timeout) *
					time.Millisecond)),
			)
			dByte = append(dByte, []byte{'\r', '\n'}...)
			_, err0 := socket.Write(dByte)
			socket.SetReadDeadline(time.Time{})
			if err0 != nil {
				return 0, err0
			}
		}
		return len(s), nil
	default:
		return 0, fmt.Errorf("only support string format")
	}
}

func (udpT *UdpTarget) Stop() {
	udpT.status = typex.SOURCE_DOWN
	if udpT.CancelCTX != nil {
		udpT.CancelCTX()
	}
}
func (udpT *UdpTarget) Details() *typex.OutEnd {
	return udpT.RuleEngine.GetOutEnd(udpT.PointId)
}
func (udpT *UdpTarget) UdpStatus(serverAddr string) error {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		return fmt.Errorf("UDP connection failed: %v", err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(udpT.mainConfig.PingPacket))
	if err != nil {
		return fmt.Errorf("failed to send data over UDP: %v", err)
	}
	return nil
}

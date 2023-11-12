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

package modbusscanner

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/component/hwportmanager"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/sirupsen/logrus"
	serial "github.com/wwhai/goserial"
)

/*
*
* CRC 计算
*
 */

func calculateCRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF

	for _, b := range data {
		crc ^= uint16(b)

		for i := 0; i < 8; i++ {
			lsb := crc & 0x0001
			crc >>= 1

			if lsb == 1 {
				crc ^= 0xA001 // 0xA001 是Modbus CRC16多项式的表示
			}
		}
	}

	return crc
}
func uint16ToBytes(val uint16) [2]byte {
	bytes := [2]byte{}
	binary.LittleEndian.PutUint16(bytes[:], val)
	return bytes
}

/*
*
* 服务调用接口
*
 */
func (cs *modbusScanner) Service(arg typex.ServiceArg) typex.ServiceResult {
	if cs.busying {
		if arg.Name == "stop" {
			if cs.cancel != nil {
				cs.cancel()
				cs.busying = false
				return typex.ServiceResult{Out: []map[string]interface{}{
					{"error": "Device busying now"},
				}}
			}
		}
		return typex.ServiceResult{Out: []map[string]interface{}{
			{"error": "Modbus Scanner Busing now"},
		}}
	}

	if arg.Name == "scan" {
		cs.busying = true
		// portUuid 从args里传进来的
		switch s := arg.Args.(type) {
		case string:
			{
				type form struct {
					PortUuid string `json:"portUuid"`
				}
				form1 := form{}
				if err := json.Unmarshal([]byte(s), &form1); err != nil {
					cs.cancel()
					cs.busying = false
					return typex.ServiceResult{Out: []map[string]interface{}{
						{"error": err.Error()},
					}}
				}
				hwPort, err := hwportmanager.GetHwPort(form1.PortUuid)
				if err != nil {
					cs.cancel()
					cs.busying = false
					return typex.ServiceResult{Out: []map[string]interface{}{
						{"error": "port not exists:" + form1.PortUuid},
					}}
				}
				hwPortConfig := hwportmanager.UartConfig{}
				switch tcfg := hwPort.Config.(type) {
				case hwportmanager.UartConfig:
					{
						hwPortConfig = tcfg
					}
				default:
					{
						cs.cancel()
						cs.busying = false
						return typex.ServiceResult{Out: []map[string]interface{}{
							{"error": "port not exists:" + form1.PortUuid},
						}}
					}
				}

				config := serial.Config{
					Address:  hwPortConfig.Uart,
					BaudRate: hwPortConfig.BaudRate,
					DataBits: hwPortConfig.DataBits,
					Parity:   hwPortConfig.Parity,
					StopBits: hwPortConfig.StopBits,
					Timeout:  2 * time.Second,
				}
				serialPort, err := serial.Open(&config)
				if err != nil {
					glogger.GLogger.WithFields(logrus.Fields{
						"topic": "plugin/ModbusScanner/" + cs.uuid,
					}).Info("Serial port open failed:", err)
					cs.busying = false
					return typex.ServiceResult{Out: []map[string]interface{}{
						{"error": err.Error()},
					}}
				}
				ctx, cancel := context.WithCancel(typex.GCTX)
				cs.ctx = ctx
				cs.cancel = cancel
				go func(p serial.Port, cs *modbusScanner) {
					defer p.Close()
					defer func() {
						cs.busying = false
					}()
					for i := 1; i <= 255; i++ {
						select {
						case <-cs.ctx.Done():
							{
								return
							}
						default:
							{
							}
						}
						glogger.GLogger.WithFields(logrus.Fields{
							"topic": "plugin/ModbusScanner/" + cs.uuid,
						}).Info(fmt.Sprintf("Start Scan Addr [%v]", i))
						test_data := [8]byte{byte(i), 0x03, 0x00, 0x00, 0x00, 0x01, 0, 0}
						crc16 := uint16ToBytes(calculateCRC16(test_data[:6]))
						test_data[6] = crc16[0]
						test_data[7] = crc16[1]
						glogger.GLogger.WithFields(logrus.Fields{
							"topic": "plugin/ModbusScanner/" + cs.uuid,
						}).Info("Send test packet:", test_data)
						_, err := serialPort.Write(test_data[:])
						if err != nil {
							glogger.GLogger.WithFields(logrus.Fields{
								"topic": "plugin/ModbusScanner/" + cs.uuid,
							}).Error("Serial port write error:", err)
							continue
						}
						time.Sleep(300 * time.Millisecond)
						received_buffer := [6]byte{}
						n, err := serialPort.Read(received_buffer[:])
						if err != nil {
							glogger.GLogger.WithFields(logrus.Fields{
								"topic": "plugin/ModbusScanner/" + cs.uuid,
							}).Error("Serial port read error:", err)
							continue
						}
						if n > 0 {
							glogger.GLogger.WithFields(logrus.Fields{
								"topic": "plugin/ModbusScanner/" + cs.uuid,
							}).Info(fmt.Sprintf("Addr [%d] Receive response:%v",
								i, received_buffer[:n]))
						}

					}
				}(serialPort, cs)
			}
		default:
			cs.busying = false
			return typex.ServiceResult{Out: []map[string]interface{}{
				{"error": "Invalid Uart config"},
			}}
		}
	}
	return typex.ServiceResult{Out: []map[string]interface{}{}}
}

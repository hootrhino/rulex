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
	"encoding/json"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/sirupsen/logrus"
	serial "github.com/wwhai/goserial"
)

/*
*
* CRC 计算
*
 */
func calculateCRC(data []byte) [2]byte {
	poly := uint16(0x8005)
	// CRC-16 多项式
	crc := uint16(0xFFFF) // 初始值为全 1
	for _, b := range data {
		crc ^= uint16(b) // 按位异或
		for i := 0; i < 8; i++ {
			if crc&0x0001 == 0x0001 {
				crc >>= 1

				crc ^= poly
			} else {
				crc >>= 1
			}
		}
	}
	return [2]byte{byte(crc & 0xFF), byte((crc >> 8) & 0xFF)}
}

/*
*
* 服务调用接口
*
 */
func (cs *modbusScanner) Service(arg typex.ServiceArg) typex.ServiceResult {
	if cs.busying {
		return typex.ServiceResult{Out: "Modbus Scanner Busing now"}
	}
	// 配置
	if arg.Name == "scan" {
		cs.busying = true
		switch s := arg.Args.(type) {
		case string:
			{
				err := json.Unmarshal([]byte(s), &cs.UartConfig)
				if err != nil {
					return typex.ServiceResult{Out: err.Error()}
				}
				if !utils.SContains([]string{"N", "E", "O"}, cs.UartConfig.Parity) {
					return typex.ServiceResult{Out: "parity value only one of 'N','O','E'"}
				}
				config := serial.Config{
					Address:  cs.UartConfig.Uart,
					BaudRate: cs.UartConfig.BaudRate,
					DataBits: cs.UartConfig.DataBits,
					Parity:   cs.UartConfig.Parity,
					StopBits: cs.UartConfig.StopBits,
					Timeout:  time.Duration(cs.UartConfig.Timeout) * time.Second,
				}
				serialPort, err := serial.Open(&config)
				if err != nil {
					glogger.GLogger.WithFields(logrus.Fields{
						"topic": "plugin/ModbusScanner/" + cs.uuid,
					}).Info("Serial port open failed:", err)
					return typex.ServiceResult{Out: err.Error()}
				}
				go func(p serial.Port, cs *modbusScanner) {
					defer p.Close()
					defer func() {
						cs.busying = false
					}()
					for i := 0; i <= 255; i++ {
						select {
						case <-typex.GCTX.Done():
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
						crc16 := calculateCRC(test_data[:6])
						test_data[6] = crc16[1]
						test_data[7] = crc16[0]
						_, err := serialPort.Write(test_data[:])
						if err != nil {
							glogger.GLogger.WithFields(logrus.Fields{
								"topic": "plugin/ModbusScanner/" + cs.uuid,
							}).Info("Serial port write error:", err)
							continue
						}
						received_buffer := []byte{}
						n, err := serialPort.Read(received_buffer)
						if err != nil {
							glogger.GLogger.WithFields(logrus.Fields{
								"topic": "plugin/ModbusScanner/" + cs.uuid,
							}).Info("Serial port read error:", err)
							continue
						}
						if n == 6 {
							glogger.GLogger.WithFields(logrus.Fields{
								"topic": "plugin/ModbusScanner/" + cs.uuid,
							}).Info(fmt.Sprintf("Addr [%d] Receive response:%v",
								i, received_buffer[:n]))
						}

					}
				}(serialPort, cs)
			}
		default:
			return typex.ServiceResult{Out: "Invalid Uart config"}
		}
	}
	return typex.ServiceResult{}
}

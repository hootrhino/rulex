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
	"errors"
	"fmt"
	"sync"
	"time"

	serial "github.com/wwhai/goserial"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type userG776 struct {
	typex.XStatus
	mainConfig common.CommonUartConfig
	status     typex.SourceState
	locker     sync.Mutex
	serialPort serial.Port
	errCount   int
}

func NewUserG776(e typex.RuleX) typex.XTarget {
	g776 := new(userG776)
	g776.RuleEngine = e
	g776.mainConfig = common.CommonUartConfig{
		Timeout:  3000,
		Uart:     "/tty/s1",
		BaudRate: 9600,
		Parity:   "N",
		DataBits: 8,
		StopBits: 0,
	}
	g776.locker = sync.Mutex{}
	g776.errCount = 0
	g776.status = typex.SOURCE_DOWN
	return g776
}

func (g776 *userG776) Init(outEndId string, configMap map[string]interface{}) error {
	g776.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &g776.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"N", "E", "O"}, g776.mainConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	return nil

}
func (g776 *userG776) Start(cctx typex.CCTX) error {
	g776.Ctx = cctx.Ctx
	g776.CancelCTX = cctx.CancelCTX
	config := serial.Config{
		Address:  g776.mainConfig.Uart,
		BaudRate: g776.mainConfig.BaudRate,
		DataBits: g776.mainConfig.DataBits,
		Parity:   g776.mainConfig.Parity,
		StopBits: g776.mainConfig.StopBits,
		Timeout:  time.Duration(g776.mainConfig.Timeout) * time.Millisecond,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("Serial.Open failed:", err)
		return err
	}
	g776.errCount = 0
	g776.serialPort = serialPort
	g776.status = typex.SOURCE_UP
	glogger.GLogger.Info("userG776 started")
	return nil
}

func (g776 *userG776) Status() typex.SourceState {
	if g776.serialPort != nil {
		// https://www.usr.cn/Download/806.html
		//  发送： AT\r
		//  接收： \r\nOK\r\n\r\n
		_, err := g776.serialPort.Write([]byte("AT\r"))
		if err != nil {
			g776.errCount++
			glogger.GLogger.Error(err)
			if g776.errCount > 5 {
				return typex.SOURCE_DOWN
			}
		}
		return typex.SOURCE_UP
	}
	return typex.SOURCE_DOWN
}

/*
*
* 数据写到串口
*
 */
func (g776 *userG776) To(data interface{}) (interface{}, error) {
	if g776.serialPort == nil {
		return 0, fmt.Errorf("serial Port invalid")
	}
	switch t := data.(type) {
	case string:
		return g776.serialPort.Write([]byte(t))
	case []byte:
		return g776.serialPort.Write(t)
	}
	return 0, fmt.Errorf("data type must be byte or string")

}

func (g776 *userG776) Stop() {
	g776.status = typex.SOURCE_DOWN
	g776.CancelCTX()
	g776.serialPort.Close()
}
func (g776 *userG776) Details() *typex.OutEnd {
	return g776.RuleEngine.GetOutEnd(g776.PointId)
}

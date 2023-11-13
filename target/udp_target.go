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
	"net"
	"net/http"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 数据推到UDP
*
 */
type UdpTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig common.HostConfig
	status     typex.SourceState
}

func NewUdpTarget(e typex.RuleX) typex.XTarget {
	udpt := new(UdpTarget)
	udpt.RuleEngine = e
	udpt.mainConfig = common.HostConfig{}
	udpt.status = typex.SOURCE_DOWN
	return udpt
}

func (udpt *UdpTarget) Init(outEndId string, configMap map[string]interface{}) error {
	udpt.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &udpt.mainConfig); err != nil {
		return err
	}

	return nil

}
func (udpt *UdpTarget) Start(cctx typex.CCTX) error {
	udpt.Ctx = cctx.Ctx
	udpt.CancelCTX = cctx.CancelCTX
	udpt.client = http.Client{}
	udpt.status = typex.SOURCE_UP
	glogger.GLogger.Info("UdpTarget started")
	return nil
}

func (udpt *UdpTarget) Status() typex.SourceState {
	return udpt.status

}
func (udpt *UdpTarget) To(data interface{}) (interface{}, error) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(udpt.mainConfig.Host),
		Port: udpt.mainConfig.Port,
	})
	if err != nil {
		return 0, err
	}
	defer socket.Close()
	switch t := data.(type) {
	case string:
		socket.Write([]byte(t))
	case []byte:
		socket.Write(t)
	default:
		glogger.GLogger.Error("unknown type:", t)
	}
	return 0, err
}

func (udpt *UdpTarget) Stop() {
	udpt.status = typex.SOURCE_STOP
	udpt.CancelCTX()
}
func (udpt *UdpTarget) Details() *typex.OutEnd {
	return udpt.RuleEngine.GetOutEnd(udpt.PointId)
}

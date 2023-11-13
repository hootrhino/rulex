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
	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/nats-io/nats.go"
)

type natsTarget struct {
	typex.XStatus
	natsConnector *nats.Conn
	mainConfig    common.NatsConfig
	status        typex.SourceState
}

func NewNatsTarget(e typex.RuleX) typex.XTarget {
	nt := &natsTarget{}
	nt.RuleEngine = e
	nt.mainConfig = common.NatsConfig{}
	return nt
}
func (nt *natsTarget) Init(outEndId string, configMap map[string]interface{}) error {
	nt.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &nt.mainConfig); err != nil {
		return err
	}
	return nil
}
func (nt *natsTarget) Start(cctx typex.CCTX) error {
	nt.Ctx = cctx.Ctx
	nt.CancelCTX = cctx.CancelCTX

	nc, err := nats.Connect(fmt.Sprintf("%s:%v", nt.mainConfig.Host, nt.mainConfig.Port), func(o *nats.Options) error {
		o.User = nt.mainConfig.Username
		o.Password = nt.mainConfig.Password
		return nil
	})
	if err != nil {
		return err
	} else {
		nt.natsConnector = nc
		nt.status = typex.SOURCE_UP
		return nil
	}
}

func (nt *natsTarget) Status() typex.SourceState {
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			return typex.SOURCE_UP
		}
	}
	return typex.SOURCE_DOWN
}

func (nt *natsTarget) Details() *typex.OutEnd {
	return nt.RuleEngine.GetOutEnd(nt.PointId)
}

// --------------------------------------------------------
// To: 数据出口
// --------------------------------------------------------
func (nt *natsTarget) To(data interface{}) (interface{}, error) {
	if nt.natsConnector != nil {
		switch t := data.(type) {
		case string:
			err := nt.natsConnector.Publish(nt.mainConfig.Topic, []byte(t))
			return nil, err
		}
		return nil, errors.New("unsupported data type")
	}
	return nil, errors.New("nats Connector is nil")
}

func (nt *natsTarget) Stop() {
	nt.status = typex.SOURCE_STOP
	nt.CancelCTX()
	if nt.natsConnector != nil {
		if nt.natsConnector.IsConnected() {
			nt.natsConnector.Drain()
			nt.natsConnector.Close()
			nt.natsConnector = nil
		}
	}

}

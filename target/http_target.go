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
	"net/http"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type HTTPTarget struct {
	typex.XStatus
	client     http.Client
	mainConfig common.HTTPConfig
	status     typex.SourceState
}

func NewHTTPTarget(e typex.RuleX) typex.XTarget {
	ht := new(HTTPTarget)
	ht.RuleEngine = e
	ht.mainConfig = common.HTTPConfig{}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *HTTPTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ht.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}

	return nil

}
func (ht *HTTPTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	ht.client = http.Client{}
	ht.status = typex.SOURCE_UP
	glogger.GLogger.Info("HTTPTarget started")
	return nil
}

func (ht *HTTPTarget) Status() typex.SourceState {
	return ht.status

}
func (ht *HTTPTarget) To(data interface{}) (interface{}, error) {
	r, err := utils.Post(ht.client, data, ht.mainConfig.Url, ht.mainConfig.Headers)
	return r, err
}

func (ht *HTTPTarget) Stop() {
	ht.status = typex.SOURCE_STOP
	ht.CancelCTX()
}
func (ht *HTTPTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}


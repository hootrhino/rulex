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
	"context"
	"fmt"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTarget struct {
	typex.XStatus
	client     *mongo.Client
	collection *mongo.Collection
	mainConfig common.MongoConfig
	status     typex.SourceState
}

func NewMongoTarget(e typex.RuleX) typex.XTarget {
	mg := new(mongoTarget)
	mg.mainConfig = common.MongoConfig{}
	mg.RuleEngine = e
	mg.status = typex.SOURCE_DOWN
	return mg
}

func (m *mongoTarget) Init(outEndId string, configMap map[string]interface{}) error {
	m.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &m.mainConfig); err != nil {
		return err
	}
	return nil
}
func (m *mongoTarget) Start(cctx typex.CCTX) error {
	m.Ctx = cctx.Ctx
	m.CancelCTX = cctx.CancelCTX
	clientOptions := options.Client().ApplyURI(m.mainConfig.MongoUrl)
	clientOptions.SetConnectTimeout(3 * time.Second)
	// clientOptions.SetDirect(true)
	client, err0 := mongo.Connect(m.Ctx, clientOptions)
	if err0 != nil {
		return err0
	}
	m.collection = client.Database(m.mainConfig.Database).Collection(m.mainConfig.Collection)
	m.client = client
	m.Enable = true
	m.status = typex.SOURCE_UP
	glogger.GLogger.Info("mongoTarget connect successfully")
	return nil

}

func (m *mongoTarget) Status() typex.SourceState {
	if m.client != nil {
		ctx, cancel := context.WithTimeout(m.Ctx, time.Second*2)
		defer cancel()
		if err1 := m.client.Ping(ctx, nil); err1 != nil {
			glogger.GLogger.Error(err1)
			return typex.SOURCE_DOWN
		}
	}
	return m.status
}

func (m *mongoTarget) Stop() {
	m.CancelCTX()
	m.status = typex.SOURCE_DOWN
	if m.client != nil {
		m.client.Disconnect(m.Ctx)
	}
}

func (m *mongoTarget) To(data interface{}) (interface{}, error) {
	switch t := data.(type) {
	case string:
		// 将 JSON 数据解析为 map
		var data map[string]interface{}

		if err := bson.UnmarshalExtJSON([]byte(t), false, &data); err != nil {
			glogger.GLogger.Error("Mongo To Failed:", err)
			return nil, err
		}
		r, err := m.collection.InsertOne(m.Ctx, data)
		if err != nil {
			glogger.GLogger.Error("Mongo To Failed:", err)
			return nil, err
		}
		return r.InsertedID, nil
	}
	return nil, fmt.Errorf("unsupported Bson type:%s", data)

}
func (m *mongoTarget) Details() *typex.OutEnd {
	return m.RuleEngine.GetOutEnd(m.PointId)
}

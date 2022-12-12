package target

import (
	"context"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

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
	client, err0 := mongo.Connect(m.Ctx, clientOptions)
	if err0 != nil {
		return err0
	}
	ctx, cancel := context.WithTimeout(m.Ctx, 3*time.Second)
	defer cancel()
	if err1 := client.Ping(ctx, nil); err1 != nil {
		return err1
	}
	m.collection = client.Database(m.mainConfig.Database).Collection(m.mainConfig.Collection)
	m.client = client
	m.Enable = true
	m.status = typex.SOURCE_UP
	glogger.GLogger.Info("mongoTarget connect successfully")
	return nil

}

func (m *mongoTarget) Test(outEndId string) bool {
	if m.client != nil {
		ctx, cancel := context.WithTimeout(m.Ctx, 3*time.Second)
		defer cancel()
		if err1 := m.client.Ping(ctx, nil); err1 != nil {
			return false
		} else {
			return true
		}
	}
	return false

}

func (m *mongoTarget) Enabled() bool {
	return m.Enable
}

func (m *mongoTarget) Reload() {
	glogger.GLogger.Info("Mongo target Reload success")
}

func (m *mongoTarget) Pause() {
	glogger.GLogger.Info("Mongo target Pause success")
}

func (m *mongoTarget) Status() typex.SourceState {
	return m.status
}

func (m *mongoTarget) Stop() {
	m.status = typex.SOURCE_STOP
	m.CancelCTX()
	m.client.Disconnect(m.Ctx)
}

func (m *mongoTarget) To(data interface{}) (interface{}, error) {
	document := bson.D{bson.E{Key: "data", Value: data}}
	r, err := m.collection.InsertOne(m.Ctx, document)
	if err != nil {
		glogger.GLogger.Error("Mongo To Failed:", err)
	}
	return r.InsertedID, err
}
func (m *mongoTarget) Details() *typex.OutEnd {
	return m.RuleEngine.GetOutEnd(m.PointId)
}

/*
*
* 配置
*
 */
func (*mongoTarget) Configs() *typex.XConfig {
	return typex.GenOutConfig(typex.MONGO_SINGLE, "MONGO_SINGLE", common.MongoConfig{})
}

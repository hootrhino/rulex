package target

import (
	"context"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct {
	MongoUrl   string `json:"mongoUrl" validate:"required"`
	Database   string `json:"database" validate:"required"`
	Collection string `json:"collection" validate:"required"`
}

//
type MongoTarget struct {
	typex.XStatus
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoTarget(e typex.RuleX) typex.XTarget {
	mg := new(MongoTarget)
	mg.RuleEngine = e
	return mg
}
func (m *MongoTarget) OnStreamApproached(data string) error {
	return nil
}
func (m *MongoTarget) Register(outEndId string) error {
	m.PointId = outEndId
	return nil
}

func (m *MongoTarget) Start() error {
	config := m.RuleEngine.GetOutEnd(m.PointId).Config
	var mainConfig mongoConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	clientOptions := options.Client().ApplyURI(mainConfig.MongoUrl)
	client, err0 := mongo.Connect(context.Background(), clientOptions)
	if err0 != nil {
		return err0
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err1 := client.Ping(ctx, nil); err1 != nil {
		return err1
	}
	m.collection = client.Database(mainConfig.Database).Collection(mainConfig.Collection)
	m.client = client
	m.Enable = true
	log.Info("MongoTarget connect successfully")
	return nil

}

func (m *MongoTarget) Test(outEndId string) bool {
	if m.client != nil {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err1 := m.client.Ping(ctx, nil); err1 != nil {
			return false
		} else {
			return true
		}
	}
	return false

}

func (m *MongoTarget) Enabled() bool {
	return m.Enable
}

func (m *MongoTarget) Reload() {
	log.Info("Mongotarget Reload success")
}

func (m *MongoTarget) Pause() {
	log.Info("Mongotarget Pause success")

}

func (m *MongoTarget) Status() typex.ResourceState {
	err1 := m.client.Ping(context.Background(), nil)
	if err1 != nil {
		log.Error(err1)
		return typex.DOWN
	} else {
		return typex.UP
	}
}

func (m *MongoTarget) Stop() {
	m.client.Disconnect(context.Background())
	log.Info("Mongotarget Stop success")
}

func (m *MongoTarget) To(data interface{}) error {
	document := bson.D{bson.E{Key: "data", Value: data}}
	_, err := m.collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Error("Mongo To Failed:", err)
	}
	return err
}
func (m *MongoTarget) Details() *typex.OutEnd {
	return m.RuleEngine.GetOutEnd(m.PointId)
}

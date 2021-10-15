package target

import (
	"context"
	"rulex/typex"

	"github.com/ngaut/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (m *MongoTarget) Register(outEndId string) error {
	m.PointId = outEndId
	return nil
}

func (m *MongoTarget) Start() error {
	config := m.RuleEngine.GetOutEnd(m.PointId).Config
	var clientOptions *options.ClientOptions
	if (*config)["mongourl"] != nil {
		clientOptions = options.Client().ApplyURI((*config)["mongourl"].(string))
	} else {
		clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	}
	client, err0 := mongo.Connect(context.TODO(), clientOptions)
	if err0 != nil {
		return err0
	}

	if (*config)["database"] != nil {
		if (*config)["collection"] != nil {
			m.collection = client.Database((*config)["database"].(string)).Collection((*config)["collection"].(string))
		} else {
			m.collection = client.Database((*config)["mongourl"].(string)).Collection("rulex_data")
		}
	} else {
		m.collection = client.Database("rulex").Collection("rulex_data")
	}
	m.client = client
	m.Enable = true
	log.Info("Mqtt connect successfully")
	return nil

}

func (m *MongoTarget) Test(outEndId string) bool {
	err1 := m.client.Ping(context.Background(), nil)
	if err1 != nil {
		log.Error(err1)
		return false
	} else {
		return true
	}
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
		m.RuleEngine.GetOutEnd(m.PointId).State = typex.UP
		return typex.UP
	}
}

func (m *MongoTarget) Stop() {
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

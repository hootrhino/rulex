package x

import (
	"context"
	"time"

	"github.com/ngaut/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//
type MongoTarget struct {
	enabled    bool
	outEndId   string
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoTarget() *MongoTarget {

	return &MongoTarget{
		enabled: false,
	}
}

func (m *MongoTarget) Register(outEndId string) error {
	m.outEndId = outEndId
	return nil
}

func (m *MongoTarget) Start(e *RuleEngine) error {
	config := e.GetOutEnd(m.outEndId).Config
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
	err1 := client.Ping(context.TODO(), nil)
	if err1 != nil {
		return err1
	} else {
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
		m.enabled = true
		e.GetOutEnd(m.outEndId).State = 1
		log.Info("Mongodb connect successfully")
		ctx := context.Background()
		go func(ctx context.Context, client *mongo.Client) {
			t := time.NewTicker(time.Duration(time.Second * 5))
			defer t.Stop()
			for {
				<-t.C
				err1 := client.Ping(ctx, nil)
				if err1 != nil {
					log.Error("client disconnected,error:", err1)
					e.GetOutEnd(m.outEndId).State = 0
				} else {
					e.GetOutEnd(m.outEndId).State = 1
				}
			}
		}(ctx, m.client)
		return nil
	}
}

func (m *MongoTarget) Test(outEndId string) bool {
	return true
}

func (m *MongoTarget) Enabled() bool {
	return m.enabled
}

func (m *MongoTarget) Reload() {
	log.Info("Mongotarget Reload success")
}

func (m *MongoTarget) Pause() {
	log.Info("Mongotarget Pause success")

}

func (m *MongoTarget) Status(e *RuleEngine) int {
	return e.GetOutEnd(m.outEndId).State
}

func (m *MongoTarget) Stop() {
	log.Info("Mongotarget Stop success")
}

func (m *MongoTarget) To(data interface{}) error {
	_, err := m.collection.InsertOne(context.TODO(), bson.D{{"data", data}})
	if err != nil {
		log.Error("Mongo To Failed:", err)
	}
	return err
}

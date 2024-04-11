package ebase

import (
	"context"
	"fmt"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newMongo(mongoose config.MongoConfig) (*mongo.Client, error) {

	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoose.Username, mongoose.Password, mongoose.Host, mongoose.Port, mongoose.Dbname)

	return mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

}

func (e *Eb) initMongo() {

	for _, mongoConfig := range e.Config.Mongo {
		mongoClient, err := newMongo(mongoConfig)
		if err != nil {
			logger.Error("mongo  初始化失败 %s %s", mongoConfig.Name, err.Error())
		}
		e.Mongo[mongoConfig.Name] = mongoClient
	}

}

package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mas/utils/config"
	"time"
)

var MongoClient *mongo.Client
var MongoConn *mongo.Database

func InitMongoClient() {
	time.Sleep(5 * time.Second)
	err := NewMongoClient()
	if err == nil {
		fmt.Println("==> Connected to MongoDB!")
	} else {
		fmt.Println("==> Cannot connected to MongoDB! Try again after a few seconds...")
	}
	// 心跳goroutine
	go checkMongoClientConnection()
	return
}

func checkMongoClientConnection () {
	time.Sleep(3 *time.Second)
	for {
		if MongoClient == nil {
			InitMongoClient()
			return
		}
		err := MongoClient.Ping(context.Background(), nil)
		if err != nil {
			for err != nil {
				fmt.Println("==> Cannot connected to MongoDB! Try again after a few seconds...")
				time.Sleep(3 *time.Second)
				err = NewMongoClient()
			}
			fmt.Println("==> Connected to MongoDB!")
		} else {
			time.Sleep(5 *time.Second)
		}
	}
}

func NewMongoClient() (err error) {
	connectString := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		config.SystemConfig.MongoDB.Username,
		config.SystemConfig.MongoDB.Password,
		config.SystemConfig.MongoDB.Host,
		config.SystemConfig.MongoDB.Port,
	)
	mongoClientOptions := options.Client().ApplyURI(connectString)
	mongoClientOptions.SetConnectTimeout(1 * time.Second)
	mongoClientOptions.SetSocketTimeout(1 * time.Second)
	MongoClient, err = mongo.Connect(context.Background(), mongoClientOptions)

	if err != nil {
		return
	}
	err = MongoClient.Ping(context.Background(), nil)
	if err != nil {
		return
	}
	MongoConn = MongoClient.Database(config.SystemConfig.MongoDB.DBName)
	return
}
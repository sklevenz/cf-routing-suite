package mongo

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type MongoDBQuery interface {
	GetMode() string
	GetRequestCounter() uint64
	IncRequestCounter()
	ResetRequestCounter()
}

type mongodb struct {
	client *mongo.Client
}

func (db *mongodb) GetMode() string {
	return "mongodb"
}

func (db *mongodb) IncRequestCounter() {
}

func (db *mongodb) GetRequestCounter() uint64 {
	return uint64(1)
}

func (db *mongodb) ResetRequestCounter() {
}

func Create() MongoDBQuery {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Printf("Error: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Can't ping mongodb. Error: %v", err)
		os.Exit(1)
	}

	q := mongodb{client: client}
	query := &q
	log.Printf("MongoDB connected: %v", client.ConnectionString())

	return query
}

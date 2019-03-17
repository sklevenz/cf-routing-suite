package mongo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

const (
	database = "cfrsdb"
)

type MongoDBQuery interface {
	RecordRequest(requestData *RequestData) ResultData
	ResetAll() ResultData
}

type mongodb struct {
	client *mongo.Client
}

type ResultData struct {
	Error       error        `json:"error"`
	Count       int64        `json:"count"`
	Message     string       `json:"message"`
	MongoData   *MongoData   `json:"mongo"`
	RequestData *RequestData `json:"request-data"`
}

type RequestData struct {
	Url             string    `json:"url"`
	Method          string    `json:"method"`
	Remote          string    `json:"remote"`
	Timestamp       time.Time `json:"timestamp"`
	XForwardedFor   string    `json:"x-forwarded-for"`
	XB3TraceId      string    `json:"x-b3-traceid"`
	XB3SpanId       string    `json:"x-b3-spanid"`
	XB3ParentSpanId string    `json:"x_b3_parentspanid"`
	Tag             string    `json:"tag"`
}

type MongoData struct {
	InsertId interface{} `json:"insert-id"`
}

func (db *mongodb) ResetAll() ResultData {
	mongodb := db.client.Database(database)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := mongodb.Drop(ctx)
	if err != nil {
		log.Printf("Mongodb drop error %v", err)
	}

	return ResultData{
		err,
		0,
		"database dropped",
		nil,
		nil,
	}
}

func (db *mongodb) RecordRequest(requestData *RequestData) ResultData {
	collection := db.client.Database(database).Collection(requestData.Tag)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, requestData)
	if err != nil {
		log.Printf("Mongodb insert one error %v", err)
	}

	var count = int64(-1)
	count, err = collection.Count(ctx, bson.M{})
	if err != nil {
		log.Printf("Mongodb count error %v", err)
	}

	return ResultData{
		err,
		count,
		"request recorded",
		&MongoData{res.InsertedID},
		requestData,
	}
}

func Dial(mode string) (query MongoDBQuery) {
	if mode == "mongodb" {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
		if err != nil {
			log.Panicf("Connect error to mongodb://localhost:27017: %v", err)
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Panicf("Can't ping mongodb. Error: %v", err)
		}

		query = &mongodb{client: client}
		log.Printf("MongoDB connected: %v", client.ConnectionString())
		return query
	} else if mode == "simulator" {
		query = &simulator{}
		log.Printf("MongoDB connected: %v", mode)
		return query
	} else {
		err := fmt.Errorf("Unsupported mode: %v, expected: [simulator|mongodb]", mode)
		log.Panicf(err.Error())
		return nil
	}
}

package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var DB *mongo.Database

func InitMongoDB(uri string, dbName string) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return err
	}

	// Verify connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
		return err
	}

	DB = mongoClient.Database(dbName)
	log.Println("Connected to MongoDB successfully")
	return nil
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}

func GetDatabase() *mongo.Database {
	return DB
}

func CloseMongoConnection(ctx context.Context) error {
	if mongoClient != nil {
		return mongoClient.Disconnect(ctx)
	}
	return nil
}

package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig represents config options for mongo database
type MongoConfig struct {
	Port     string
	Address  string
	Database string
}

var (
	// client is the db connection handle
	client *mongo.Client
	// config holds the database config properties
	config *MongoConfig
)

// OpenConnection opens a mysql database connection
func OpenConnection() {
	if client != nil {
		return
	}

	if config == nil {
		LoadConfig()
	}

	connStr := buildConnectionString(config)
	client, _ = mongo.NewClient(options.Client().ApplyURI(connStr))

	// Background is the root of any context tree, it is never cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateOne(coll string, filter *bson.D, update *bson.D) (string, error) {
	if client == nil {
		OpenConnection()
	}

	collection := client.Database(config.Database).Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return "", err
	}

	if result.UpsertedID == nil {
		return "", err
	}

	return result.UpsertedID.(primitive.ObjectID).Hex(), nil
}

// InsertOne inserts bson document data into collection and returns the insertedID
func InsertOne(coll string, document interface{}) (string, error) {
	if client == nil {
		OpenConnection()
	}

	collection := client.Database(config.Database).Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, document)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func Find(coll string, filter bson.M, opts ...*options.FindOptions) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database(config.Database).Collection(coll)
	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	results := make([]bson.M, 0)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}

// buildConnectionString constructs database conn string from config
func buildConnectionString(conf *MongoConfig) string {
	return fmt.Sprintf("%s:%s", conf.Address, conf.Port)
}

// LoadConfig initializes mongo config
func LoadConfig() {
	if config != nil {
		return
	}

	config = &MongoConfig{
		Port:     os.Getenv("MONGO_PORT"),
		Address:  os.Getenv("MONGO_ADDRESS"),
		Database: os.Getenv("MONGO_DATABASE"),
	}
}

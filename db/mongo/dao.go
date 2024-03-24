package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var Db DBHandler

type DBHandler interface {
	InsertOne(ctx context.Context, dbName, table string, data interface{}) error
	CreateIndex(ctx context.Context, dbName, table string, keys bson.D, options *options.IndexOptions) error
}

type MongoDB struct {
	client *mongo.Client
}

func InitMongoDB(mongoUri string) DBHandler {
	// data
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("connect db error", err)
	}
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal("ping db error", err)
	}
	Db = &MongoDB{client}
	return Db
}

func (d *MongoDB) InsertOne(ctx context.Context, dbName, table string, data interface{}) error {
	_, err := d.client.Database(dbName).Collection(table).InsertOne(ctx, data)
	return err
}

func (d *MongoDB) CreateIndex(ctx context.Context, dbName, table string, keys bson.D, options *options.IndexOptions) error {
	_, err := d.client.Database(dbName).Collection(table).Indexes().CreateOne(ctx, mongo.IndexModel{keys, options})
	return err
}

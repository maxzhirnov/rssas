package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoStorage(connString, database, collection string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	return &MongoStorage{client: client, database: database, collection: collection}, nil
}

func (c *MongoStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}

func (c *MongoStorage) Bootstrap() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"guid": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := c.client.Database(c.database).Collection(c.collection).Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) InsertMany(document []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := s.client.Database(s.database).Collection(s.collection)
	_, err := coll.InsertMany(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

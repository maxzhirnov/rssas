package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client   *mongo.Client
	database string
}

func NewMongoStorage(connString, database string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	return &MongoStorage{client: client, database: database}, nil
}

func (s *MongoStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *MongoStorage) Bootstrap() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"guid": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := s.client.Database(s.database).Collection("items").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) InsertMany(document []interface{}, collection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := s.client.Database(s.database).Collection(collection)
	_, err := coll.InsertMany(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) InsertOne(document interface{}, collection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := s.client.Database(s.database).Collection(collection)
	_, err := coll.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) GetFeedsLinks() ([]string, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	coll := s.client.Database(s.database).Collection("feeds")

	links, err := coll.Distinct(context.TODO(), "link", bson.D{})
	if err != nil {
		return nil, err
	}

	// Конвертируем []interface{} в []string
	var result []string
	for _, link := range links {
		result = append(result, link.(string))
	}

	return result, nil
}

package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"rssas/internal/log"
)

type MongoStorage struct {
	client   *mongo.Client
	database string
	logger   *log.Logger
}

func NewMongoStorage(connString, database string, logger *log.Logger) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	return &MongoStorage{
		client:   client,
		database: database,
		logger:   logger,
	}, nil
}

func (s *MongoStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *MongoStorage) Bootstrap() error {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	if err := s.createIndex("guid", "items"); err != nil {
		return err
	}

	if err := s.createIndex("link", "feeds"); err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) InsertMany(document []interface{}, collection string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	coll := s.client.Database(s.database).Collection(collection)
	_, err := coll.InsertMany(ctx, document, options.InsertMany().SetOrdered(false))
	if err != nil {
		s.logger.Log.Error(err)
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
		s.logger.Log.Error(err)
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
		s.logger.Log.Error(err)
		return nil, err
	}

	// Конвертируем []interface{} в []string
	var result []string
	for _, link := range links {
		result = append(result, link.(string))
	}

	return result, nil
}

func (s *MongoStorage) createIndex(fieldName, collection string) error {
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			fieldName: 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := s.client.Database(s.database).Collection(collection).Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		s.logger.Log.Error(err)
		return err
	}

	return nil
}

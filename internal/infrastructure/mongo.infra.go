package infrastructure

import (
	"context"
	"ecommerce/internal/env"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

// type MongoDriverMethod interface {
//   Connect(env *env.Config) (*mongo.Client, error)
//   Disconnect(client *mongo.Client) error
// }

// type MongoClient struct {
//   logger *zap.Logger
//   timeout *time.Duration
// }

// func NewMongoClient(logger *zap.Logger) *MongoClient {
//   return &MongoClient{ logger: logger }
// }

func (m *MongoClient) MongoClient(env *env.Config) (*mongo.Client, error) {
	if env.DB.ConfigDBUrl == "" {
		return nil, nil
	}
	uriMongo := env.DB.ConfigDBUrl
	client, err := mongo.Connect(options.Client().ApplyURI(uriMongo))
	if err != nil {
		m.logger.Fatal("Failed to connect to MongoDB:", zap.Error(err))
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			m.logger.Fatal("Failed to disconnect from MongoDB:", zap.Error(err))
		}
	}()
	return client, nil
}

type MongoDriverMethod interface {
	Connect(env *env.Config) (*mongo.Client, error)
	Disconnect(client *mongo.Client) error
}

type MongoClient struct {
	logger  *zap.Logger
	timeout time.Duration
}

func NewMongoClient(logger *zap.Logger) *MongoClient {
	return &MongoClient{
		logger:  logger,
		timeout: 10 * time.Second,
	}
}

func (m *MongoClient) Connect(env *env.Config) (*mongo.Client, error) {
	uriMongo := env.DB.ConfigDBUrl
	client, err := mongo.Connect(options.Client().ApplyURI(uriMongo))
	if err != nil {
		m.logger.Fatal("Failed to connect to MongoDB:", zap.Error(err))
	  return nil,err
  }
	return client, nil
}

func (m *MongoClient) Disconnect(client *mongo.Client) error {
	if err := client.Disconnect(context.TODO()); err != nil {
		m.logger.Error("Failed to disconnect Mongo", zap.Error(err))
		return err
	}

	m.logger.Info("Disconnected Mongo successfully")
	return nil
}

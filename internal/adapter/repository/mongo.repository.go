package repository

import (
	"ecommerce/internal/application/shopee"
	"ecommerce/internal/env"

	"go.uber.org/zap"
)

type MongoCollectionRepository struct {
  // collection *mongo.Collection
  ShopeeAuthCollection shopee.ShopeeAuthRepository
  // ShopeeShop *mongo.Collection
}

func NewMongoCollectionRepository(shopeeAuth shopee.ShopeeAuthRepository, logger *zap.Logger, cfg *env.Config ) *MongoCollectionRepository {
  // dbName := cfg.DB.ConfigDBName
  // if dbName == "" { 
  //   dbName = "auth"
  //   logger.Error("Mongo Error : no database name provided") 
  // }
  return &MongoCollectionRepository{
    ShopeeAuthCollection: shopeeAuth,
  }
}

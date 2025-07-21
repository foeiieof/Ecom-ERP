package shopee

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type ShopeeAuthRepository interface {
  GetShopeeAuthByShopId(shopId string) (string, error)
}

type shopeeAuthRepo struct {
  logger *zap.Logger
  db *mongo.Collection
}


func NewShopeeAuthRepository(db *mongo.Collection, log *zap.Logger) ShopeeAuthRepository {
  return &shopeeAuthRepo{ db: db, logger: log, }
}


func (r *shopeeAuthRepo)GetShopeeAuthByShopId(shopId string) (string,error) {

  if shopId == "" {
    return "", errors.New("shopId is required")
  }

  res := r.db.FindOne(context.TODO(), bson.M{"shop_id": shopId})

  if res.Err() != nil {
    return "", res.Err()
  }

  var data ShopeeAuthModel

  if err := res.Decode(&data); err != nil {
    return "", err
  }

  // Create Demo 
  // newShopeeModel := &ShopeeAuthModel{ ShopeeID: shopId, } 
  // r.db.InsertOne(context.TODO(), newShopeeModel)
  // if err := r.db.(context.TODO(), readpref.Primary()); err != nil {
  //   return "", nil
  // }


  return data.AccessToken, nil
}



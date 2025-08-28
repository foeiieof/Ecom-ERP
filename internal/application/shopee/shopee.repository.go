package shopee

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

// -- ShopeeAuthRepository
// -- ShopeeAuthResponseRepository
type ShopeeAuthRepository interface {
	InitRepository() error
	CreateShopeeAuth(partnerID string, shopId string, codeID string, accessToken string, refreshToken string) (*ShopeeAuthModel, error)
	GetShopeeShopAuthByShopId(shopId string) (*ShopeeAuthModel, error)
  UpdateShopeeShopAuth(partnerID string , shopID string ,accessToken string, refreshToken string) (*ShopeeAuthModel, error)
}

type shopeeAuthRepo struct {
	logger *zap.Logger
	db     *mongo.Collection
}

func NewShopeeAuthRepository(db *mongo.Collection, log *zap.Logger) ShopeeAuthRepository {
	return &shopeeAuthRepo{db: db, logger: log}
}

func (r *shopeeAuthRepo) InitRepository() error {
	// indexs := []mongo.IndexModel{
	//   {
	//     Keys: bson.D{{Key:"partner_id", Value: 1}},
	//     Options: options.Index().SetUnique(true),
	//   },
	// }

	// _, err := r.db.Indexes().CreateMany(context.TODO(),indexs)

	// if err != nil {
	//   r.logger.Error("error creating index", zap.Error(err))
	//   return errors.New("ShopeePartnerRepository.InitRepository: failed creating index InitRepository")
	// }
	// r.logger.Info("ShopeePartnerRepository.InitRepository: index created")
	return nil
}

func (r *shopeeAuthRepo) CreateShopeeAuth(partnerID string,shopId string, codeID string, accessToken string, refreshToken string) (*ShopeeAuthModel, error) {
	data := &ShopeeAuthModel{
    PartnerID: partnerID,
		ShopID:       shopId,
		Code:         codeID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Hour * 4),
		CreatedBy:    "admin",
		CreatedAt:    time.Now(),
	}
	if _, err := r.db.InsertOne(context.TODO(), data); err != nil {
		return nil, errors.New("failed to insert shopee auth repository")
	}
	return data, nil
}

func (r *shopeeAuthRepo) GetShopeeShopAuthByShopId(shopId string) (*ShopeeAuthModel, error) {

	if shopId == "" {
		return nil, errors.New("shopId is required")
	}

	res := r.db.FindOne(context.TODO(), bson.M{"shop_id": shopId})

	if res.Err() != nil {
    errorLog := res.Err().Error()
    parseError := strings.SplitN(errorLog, ":", 2)
    if len(parseError) == 2 {
      // key := parseError[0]
      val := parseError[1]
      r.logger.Debug("GetShopeeShopAuthByShopId", zap.String("detail", val))
    }

    // r.logger.Debug("GetShopeeShopAuthByShopId", zap.String("error", res.Err().Error()) )
		return nil, errors.New(parseError[1])
	}

	var data ShopeeAuthModel
	if err := res.Decode(&data); err != nil {
		return nil, err
	}

  // r.logger.Debug("GetShopeeShopAuthByShopId", zap.Any("data:", data))

	return &data, nil
}

func (r *shopeeAuthRepo) UpdateShopeeShopAuth(partnerID string , shopID string ,accessToken string, refreshToken string) (*ShopeeAuthModel, error) {

  if shopID == "" || accessToken == "" || refreshToken == "" || partnerID == ""{
    return nil, errors.New("shopId is required")
  }

  filter := bson.M{"partner_id": partnerID, "shop_id": shopID}
  update := bson.M{
    "$set": bson.M{
      "access_token" : accessToken,
      "refresh_token": refreshToken,
      "expired_at"   : time.Now().Add(time.Hour * 4),
      "modified_at"  : time.Now(),
      "modified_by"  : "admin",

    },
  }
  opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
  var updateShopeeAuth ShopeeAuthModel

  err := r.db.FindOneAndUpdate(context.TODO(), filter, update, opt).Decode(&updateShopeeAuth)
  if err != nil {
    errorLog := err.Error()
    parseError := strings.SplitN(errorLog, ":", 2)
    if len(parseError) == 2 {
      // key := parseError[0]
      val := parseError[1]
      r.logger.Debug("repository.ShopeeAuthRepository.UpdateShopeeShopAuth:", zap.String("error", val))
    }
  return nil, errors.New("repository.ShopeeAuthRepository.UpdateShopeeShopAuth: Failed to update accesses&refresh token")
  }
  return &updateShopeeAuth, nil
}

// -- ShopeeAuthRequestRepository
type ShopeeAuthRequestRepository interface {
	InitRepository() error
	SaveShopeeAuthRequestWithName(partnerId string, partnerKey string, partnerName string, generatedUrl string) (*ShopeeAuthRequestModel, error)
}

type shopeeAuthRequestRepo struct {
	logger *zap.Logger
	db     *mongo.Collection
}

func NewShopeeAuthRequestRepository(db *mongo.Collection, log *zap.Logger) ShopeeAuthRequestRepository {
	return &shopeeAuthRequestRepo{db: db, logger: log}
}

func (r *shopeeAuthRequestRepo) InitRepository() error {
	// indexs := []mongo.IndexModel{
	//   {
	//     Keys: bson.D{{Key:"partner_id", Value: 1}},
	//     Options: options.Index().SetUnique(true),
	//   },
	// }

	// _, err := r.db.Indexes().CreateMany(context.TODO(),indexs)

	// if err != nil {
	//   r.logger.Error("error creating index", zap.Error(err))
	//   return errors.New("ShopeePartnerRepository.InitRepository: failed creating index InitRepository")
	// }
	// r.logger.Info("ShopeePartnerRepository.InitRepository: index created")
	return nil
}

func (r *shopeeAuthRequestRepo) SaveShopeeAuthRequestWithName(partnerId string, partnerKey string, partnerName string, generatedUrl string) (*ShopeeAuthRequestModel, error) {
	data := &ShopeeAuthRequestModel{
		PartnerID:    partnerId,
		PartnerKey:   partnerKey,
		PartnerName:  partnerName,
		GeneratedUrl: generatedUrl,
		CreatedBy:    "admin",
		CreatedAt:    time.Now(),
	}
	if _, err := r.db.InsertOne(context.TODO(), data); err != nil {
		return nil, errors.New("failed to insert shopee auth request")
	}
	return data, nil
}

// -- ShopeePartnerRepository
// type ShopeePartnerRepository interface {
// 	InitRepository() error
// 	CreateShopeePartner(partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error)
// 	GetShopeePartnerByPartnerId(partnerId string) (*ShopeePartnerModel, error)
// }
// type shopeePartnerRepo struct {
// 	logger *zap.Logger
// 	db     *mongo.Collection
// }

// func NewShopeePartnerRepository(db *mongo.Collection, log *zap.Logger) ShopeePartnerRepository {
// 	return &shopeePartnerRepo{db: db, logger: log}
// }

// func (r *shopeePartnerRepo) InitRepository() error {
// 	indexs := []mongo.IndexModel{
// 		{
// 			Keys:    bson.D{{Key: "partner_id", Value: 1}},
// 			Options: options.Index().SetUnique(true),
// 		},
// 	}

// 	_, err := r.db.Indexes().CreateMany(context.TODO(), indexs)
// 	if err != nil {
// 		r.logger.Error("error creating index", zap.Error(err))
// 		return errors.New("ShopeePartnerRepository.InitRepository: failed creating index InitRepository")
// 	}
// 	r.logger.Info("ShopeePartnerRepository.InitRepository: index created")
// 	return nil
// }

// func isSameKey(a, b bson.D) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i := range a {
// 		if a[i].Key != b[i].Key || a[i].Value != b[i].Value {
// 			return false
// 		}
// 	}
// 	return true
// }


// func (r *shopeePartnerRepo) InitRepository() error {

// 	requiredIndexes := []mongo.IndexModel{
// 		{
//       Keys: bson.D{{Key: "partner_id", Value: 1}}, 
//       Options: options.Index().SetUnique(true)},
// 		// {Keys: bson.D{{Key: "shop_id", Value: 1}}, Options: options.Index().SetUnique(true)},
// 	}

// 	_, err := r.db.Indexes().CreateMany(context.TODO(), requiredIndexes)
// 	if err != nil {
// 		r.logger.Error("error creating index", zap.Error(err))
// 		return errors.New("ShopeePartnerRepository.InitRepository: failed creating index InitRepository")
// 	}
// 	r.logger.Info("ShopeePartnerRepository.InitRepository: index created")
// 	return nil
	// existingKeys := []bson.D{}

	// cursor, _ := r.db.Indexes().List(context.TODO())
	// for cursor.Next(context.TODO()) {
	// 	var index bson.M
	// 	_ = cursor.Decode(&index)

	// 	if key, ok := index["key"].(bson.M); ok {
	// 		var keys bson.D
	// 		for k, v := range key {
	// 			keys = append(keys, bson.E{Key: k, Value: v})
	// 		}
	// 		existingKeys = append(existingKeys, keys)
	// 	}
	// }

	// var toCreate []mongo.IndexModel
	// for _, idx := range requiredIndexes {
	// 	found := false
	// 	for _, exist := range existingKeys {
	//        r.logger.Debug("index already exist", zap.String("index", exist.String()))
	// 		if isSameKey(idx.Keys.(bson.D), exist) {

	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		toCreate = append(toCreate, idx)
	// 	}
	// }

// }

// func (r *shopeePartnerRepo) CreateShopeePartner(partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error) {

// 	var existing ShopeePartnerModel

// 	filter := bson.M{"partner_id": partnerId}
// 	err := r.db.FindOne(context.TODO(), filter).Decode(&existing)

// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			data := &ShopeePartnerModel{
// 				PartnerID:   partnerId,
// 				PartnerKey:  partnerKey,
// 				PartnerName: partnerName,
// 				CreatedBy:   "admin",
// 				CreatedAt:   time.Now(),
// 				MoidifiedAt: time.Now(),
// 			}

// 			_, err := r.db.InsertOne(context.TODO(), data)
// 			if err != nil {
// 				r.logger.Error("Insert failed", zap.String("error", err.Error()))
// 				return nil, errors.New("failed to insert")
// 			}
// 			return data, nil
// 		} else {
// 			return nil, err
// 		}
// 	}
// 	return &existing, nil
	// if _, err := r.db.InsertOne(context.TODO(), data); err != nil {
	// 	r.logger.Error("SaveShopeePartner",
	// 		zap.String("component", "repository.SaveShopeePartner"),
	// 		zap.String("error", err.Error()),
	// 	)
	// 	return nil, errors.New("failed to insert shopee partner")
	// }
	// return &data, nil
// }

// func (r *shopeePartnerRepo) GetShopeePartnerByPartnerId(partnerId string) (*ShopeePartnerModel, error) {
// 	var data ShopeePartnerModel

// 	filter := bson.M{"partner_id": partnerId}
// 	err := r.db.FindOne(context.TODO(), filter).Decode(&data)
// 	if err != nil {
// 		return nil, errors.New("not found shopee partner")
// 	}
// 	return &data, nil
// }

// {
// //   Method() error
// // }
// // type shopeePartnerRepo struct {
// //   logger *zap.Logger
// //   db *mongo.Collection
// // }
// // func NewShopeePartnerRepository(db *mongo.Collection, log *zap.Logger) ShopeePartnerRepository {}
// // func (r *shopeePartnerRepo) Method() error {)


// ------------------------------- Repository Template ----------------------------------------------------------
// type ShopeePartnerRepository interface {
//   Method() error
// }
// type shopeePartnerRepo struct {
//   logger *zap.Logger
//   db *mongo.Collection
// }
// func NewShopeePartnerRepository(db *mongo.Collection, log *zap.Logger) ShopeePartnerRepository {}
// func (r *shopeePartnerRepo) Method() error {)
// ------------------------------- End Repository Template -------------------------------------------------------

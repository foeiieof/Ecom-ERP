package shopee

import (
	"context"
	"errors"
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
  CreateShopeeAuth(shopId string, codeID string, accessToken string, refreshToken string) (*ShopeeAuthModel, error) 
	GetShopeeAuthByShopId(shopId string) (*ShopeeAuthModel, error)
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

func (r *shopeeAuthRepo) CreateShopeeAuth(shopId string,codeID string, accessToken string, refreshToken string) (*ShopeeAuthModel, error) {
  data := &ShopeeAuthModel{
    ShopID:      shopId,
    Code:        codeID,
    AccessToken: accessToken,
    RefreshToken: refreshToken,
    ExpiredAt:   time.Now().Add(time.Hour * 4),
    CreatedBy:   "admin",
    CreatedAt:   time.Now(),
  }
  if _, err := r.db.InsertOne(context.TODO(), data); err != nil {
    return nil, errors.New("failed to insert shopee auth repository")
  }
  return data, nil
}

func (r *shopeeAuthRepo) GetShopeeAuthByShopId(shopId string) (*ShopeeAuthModel, error) {

	if shopId == "" { return nil, errors.New("shopId is required") }

	res := r.db.FindOne(context.TODO(), bson.M{"shop_id": shopId})

	if res.Err() != nil {
    return nil, errors.New("No Documents" )
  }

	var data ShopeeAuthModel
	if err := res.Decode(&data); err != nil { 
    return nil, err 
  }

	return &data, nil
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
type ShopeePartnerRepository interface {
	InitRepository() error
	CreateShopeePartner(partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error)
  GetShopeePartnerByPartnerId(partnerId string) (*ShopeePartnerModel, error)
}
type shopeePartnerRepo struct {
	logger *zap.Logger
	db     *mongo.Collection
}

func NewShopeePartnerRepository(db *mongo.Collection, log *zap.Logger) ShopeePartnerRepository {
	return &shopeePartnerRepo{db: db, logger: log}
}

func (r *shopeePartnerRepo) InitRepository() error {
	indexs := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "partner_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := r.db.Indexes().CreateMany(context.TODO(), indexs)
	if err != nil {
		r.logger.Error("error creating index", zap.Error(err))
		return errors.New("ShopeePartnerRepository.InitRepository: failed creating index InitRepository")
	}
	r.logger.Info("ShopeePartnerRepository.InitRepository: index created")
	return nil
}

func (r *shopeePartnerRepo) CreateShopeePartner(partnerId string, partnerKey string, partnerName string) (*ShopeePartnerModel, error) {

	var existing ShopeePartnerModel

	filter := bson.M{"partner_id": partnerId}
	err := r.db.FindOne(context.TODO(), filter).Decode(&existing)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			data := &ShopeePartnerModel{
				PartnerID:   partnerId,
				PartnerKey:  partnerKey,
				PartnerName: partnerName,
				CreatedBy:   "admin",
				CreatedAt:   time.Now(),
				MoidifiedAt: time.Now(),
			}

			_, err := r.db.InsertOne(context.TODO(), data)
			if err != nil {
				r.logger.Error("Insert failed", zap.String("error", err.Error()))
				return nil, errors.New("failed to insert")
			}
			return data, nil
		} else { return nil, err }
  }
  return &existing, nil
	// if _, err := r.db.InsertOne(context.TODO(), data); err != nil {
	// 	r.logger.Error("SaveShopeePartner",
	// 		zap.String("component", "repository.SaveShopeePartner"),
	// 		zap.String("error", err.Error()),
	// 	)
	// 	return nil, errors.New("failed to insert shopee partner")
	// }
	// return &data, nil
}

func (r *shopeePartnerRepo)GetShopeePartnerByPartnerId(partnerId string) (*ShopeePartnerModel, error) {
 var data ShopeePartnerModel

  filter := bson.M{"partner_id": partnerId}
  err := r.db.FindOne(context.TODO(), filter).Decode(&data)
  if err != nil {
  return nil , errors.New("not found shopee partner")
  }
  return &data, nil
}

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

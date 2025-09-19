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
  UpdateShopeeShopAuth(partnerID string , code string,shopID string ,accessToken string, refreshToken string) (*ShopeeAuthModel, error)
}

type shopeeAuthRepo struct {
	logger *zap.Logger
	db     *mongo.Collection
}

func NewShopeeAuthRepository(db *mongo.Collection, log *zap.Logger) ShopeeAuthRepository {
	return &shopeeAuthRepo{db: db, logger: log}
}

func (r *shopeeAuthRepo) InitRepository() error {
	indexs := []mongo.IndexModel{
	  {
	    Keys: bson.D{{Key:"shop_id", Value: 1}},
	    Options: options.Index().SetUnique(true),
	  },
	}

	_, err := r.db.Indexes().CreateMany(context.TODO(),indexs)
	if err != nil {
	  r.logger.Error("error creating index", zap.Error(err))
	  return errors.New("ShopeePartnerRepository.InitRepository: failed creating index InitRepository")
	}

	r.logger.Info("ShopeePartnerRepository.InitRepository: index created")

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

func (r *shopeeAuthRepo) UpdateShopeeShopAuth(partnerID string , code string,shopID string ,accessToken string, refreshToken string) (*ShopeeAuthModel, error) {

  if shopID == "" || accessToken == "" || refreshToken == "" || partnerID == ""{
    return nil, errors.New("shopId is required")
  }

  filter := bson.M{"partner_id": partnerID, "shop_id": shopID}



  update := bson.M{
      "access_token" : accessToken,
      "refresh_token": refreshToken,
      "expired_at"   : time.Now().Add(time.Hour * 4),
      "modified_at"  : time.Now(),
      "modified_by"  : "admin",
  }

  if code != "" {
    update["code"] = code
  }

  opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
  var updateShopeeAuth ShopeeAuthModel

  err := r.db.FindOneAndUpdate(context.TODO(), filter,bson.M{"$set": update} , opt).Decode(&updateShopeeAuth)
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


// -- ShopeeShopDetailsRepository
// -- user : ShopeeShopDetailsModel, ShopeeShopProfileModel 
type ShopeeShopDetailsRepository interface {
  InitRepository() error
  CreateShopeeShopDetails(ctx context.Context, shop *ShopeeShopDetailsEntityDTO) (*ShopeeShopDetailsEntityDTO,error)
  GetAllShopeeShopDetails(ctx context.Context) ([]ShopeeShopDetailsEntityDTO, error)
  GetShopeeShopDetailsByShopID(ctx context.Context, shopID string)(*ShopeeShopDetailsEntityDTO, error) 
  // UpdateShopeeShopDetails(ctx context.Context) error
  // DeleteShopeeShopDetails(ctx context.Context) error
}
type shopeeShopDetailsRepo struct {
  Logger *zap.Logger
  DB     *mongo.Collection
}
func NewShopeeShopDetailsRepository(db *mongo.Collection, log *zap.Logger) ShopeeShopDetailsRepository {
  return &shopeeShopDetailsRepo{ DB: db, Logger : log }
}
func (r *shopeeShopDetailsRepo) InitRepository() error {
  indexs := []mongo.IndexModel{
    {
      Keys: bson.D{{ Key: "shop_id", Value: 1 }} ,
      Options: options.Index().SetUnique(true),
    },
  }

  _,err := r.DB.Indexes().CreateMany(context.TODO(), indexs)
  if err != nil {
    r.Logger.Error("error creating index", zap.Error(err))
    return errors.New("ShopeeShopDetailsRepository.InitRepository: failed creating index InitRepository")
  }

  r.Logger.Info("ShopeePartnerRepository.InitRepository: index created")
  return nil
}

func (r *shopeeShopDetailsRepo) CreateShopeeShopDetails(ctx context.Context, shop *ShopeeShopDetailsEntityDTO) (*ShopeeShopDetailsEntityDTO,error){
  // 0 convert to model
  object := ShopeeShopEntityToModel(*shop)
  // 1. insert to db
  res,err := r.DB.InsertOne(ctx, object)
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments){
      return nil, errors.New("duplicate ShopID")
    }
    return nil, err }
  // for normal string in _id 
  if  oid, ok := res.InsertedID.(bson.ObjectID) ; !ok { object.ID = oid }
  // 2. convert to entity
  objParse := ShopeeShopModelToEntity(object)
  return objParse,nil
}

func (r *shopeeShopDetailsRepo)GetAllShopeeShopDetails(ctx context.Context) ([]ShopeeShopDetailsEntityDTO, error) {

  var models []ShopeeShopDetailsModel

  cursor,err := r.DB.Find(ctx, bson.M{})
  if err != nil { return nil, err}
  defer cursor.Close(ctx)

  if err := cursor.All(ctx,&models); err != nil { return nil, err}

  resModels := make([]ShopeeShopDetailsEntityDTO, len(models))
  for i,u := range models {
    resModels[i] = *ShopeeShopModelToEntity(&u)
  }

  return resModels, nil
} 

func (r *shopeeShopDetailsRepo)GetShopeeShopDetailsByShopID(ctx context.Context, shopID string)(*ShopeeShopDetailsEntityDTO, error) {
  var model ShopeeShopDetailsModel
  filter := bson.M{ "shop_id": shopID}

  err := r.DB.FindOne(ctx, filter).Decode(&model)
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return nil, errors.New("ShopID not found")
    }
    return nil, err
  }

  // parse Model -> DTO
  modelParse := ShopeeShopModelToEntity(&model)
  return modelParse,nil
}


// ----------------- [Repository] - Start.Collection("shop_order") ----------------

type ShopeeOrderRepository interface {
  InitRepository() error
  CrateShopeeOrderWithDetails(ctx context.Context, order *ShopeeOrderEntity) (*ShopeeOrderEntity,error)
  GetShopeeOrderByOrderSN(ctx context.Context, orderSN string) (*ShopeeOrderEntity,error)
}
type shopeeOrderRepository struct {
  Logger *zap.Logger
  DB *mongo.Collection
}
func NewShopeeOrderRepository(db *mongo.Collection, log *zap.Logger) ShopeeOrderRepository {
  return &shopeeOrderRepository{ Logger: log, DB: db, } 
}

func (r *shopeeOrderRepository)InitRepository() error {
  indexs := []mongo.IndexModel{
    {
      Keys: bson.D{{ Key: "order_sn", Value: 1}},
      Options: options.Index().SetUnique(true),
    },
  }

  _,err := r.DB.Indexes().CreateMany(context.TODO(), indexs)
  if err != nil { 
    r.Logger.Error("error creating index", zap.Error(err))
	  return errors.New("ShopeeOrderRepository.InitRepository: failed creating index InitRepository")
  }
	
  r.Logger.Info("ShopeePartnerRepository.InitRepository: index created");return nil
}

func (r *shopeeOrderRepository)CrateShopeeOrderWithDetails(ctx context.Context, order *ShopeeOrderEntity) (*ShopeeOrderEntity,error) {

  orderModel := ShopeeOrderEntityToModel(order)
  orderModel.CreatedAt = time.Now()
  orderModel.UpdatedAt = time.Now()

  res,err := r.DB.InsertOne(ctx,orderModel)
  if err != nil { 
    r.Logger.Debug("repo.ShopeeOrder.CreateShopeeOrderWithDetails", zap.String("res,err", err.Error()))
    return nil , err
  }
  if oid, ok := res.InsertedID.(bson.ObjectID) ; !ok {
    orderModel.ID = oid
  }

  orderEnti := ShopeeOrderModelToEntity(orderModel)

  return orderEnti,nil
}

func (r *shopeeOrderRepository)GetShopeeOrderByOrderSN(ctx context.Context, orderSN string) (*ShopeeOrderEntity,error) {

  var order ShopeeOrderModel
  filter := bson.M{"order_sn": orderSN}
  err := r.DB.FindOne(ctx,filter).Decode(&order)
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments){
      return nil, errors.New("OrderSN not found")
    }
    return nil, err
  }
  // parse to Entity
  entity := ShopeeOrderModelToEntity(&order)
  return entity,nil 
}

// ----------------- [Repository] - End.Collection("shop_order") ----------------


// -- ShopeePartnerRepository
// type ShopeePartnerRepository interface k
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


// ----------------- [Repository/usecase/Handler] - Start.Collection("Shop_rder") ----------------

// ----------------- [Repository/usecase/Handler] - End.Collection("Shop_rder") ----------------

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


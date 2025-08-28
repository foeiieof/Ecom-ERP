package partner

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type ShopeePartnerModel struct {
  ID          bson.ObjectID `bson:"_id"`
  PartnerName string        `bson:"partner_name"`
  PartnerID   string        `bson:"partner_id"`
  SecretKey   string        `bson:"secret_key"`
  CreatedAt   time.Time     `bson:"created_at"`
  CreatedBy   string        `bson:"created_by"`
  UpdatedAt   time.Time     `bson:"updated_at"`
  UpdatedBy   string        `bson:"updated_by"`
}

type ShopeePartnerEntity struct {
  ID          string
  PartnerName string        
  PartnerID   string      
  SecretKey   string          
  CreatedAt   time.Time
  CreatedBy   string
  UpdatedAt   time.Time
  UpdatedBy   string
}

type ShopeePartnerDTO struct {
  ID          string
  PartnerName string        
  PartnerID   string        
  SecretKey   string
  CreatedAt   time.Time
  CreatedBy   string
  UpdatedAt   time.Time
  UpdatedBy   string
}

func ShopeePartnerModelToEntity(model ShopeePartnerModel) *ShopeePartnerEntity {
  return &ShopeePartnerEntity{
    ID: model.ID.Hex(),
    PartnerName: model.PartnerName,
    PartnerID: model.PartnerID,
    SecretKey: model.SecretKey,
    CreatedAt: model.CreatedAt,
    UpdatedAt: model.UpdatedAt,
    CreatedBy: model.CreatedBy,
    UpdatedBy: model.UpdatedBy,
  }
}

func ShopeePartnerEntityToDTO(enti ShopeePartnerEntity) *ShopeePartnerDTO {
  return &ShopeePartnerDTO{
    ID: enti.ID,
    PartnerID: enti.PartnerID,
    PartnerName: enti.PartnerName,
    SecretKey: enti.SecretKey,
    CreatedAt: enti.CreatedAt,
    CreatedBy: enti.CreatedBy,
    UpdatedAt: enti.UpdatedAt,
    UpdatedBy: enti.UpdatedBy,
  }
}

func ShopeePartnerEntityToModel(enti ShopeePartnerEntity) *ShopeePartnerModel{
  modelId := bson.NewObjectID()
  if enti.ID != "" {
    if ok,err := bson.ObjectIDFromHex(enti.PartnerID); err == nil {
      modelId = ok 
    } 
  }

  return &ShopeePartnerModel{
    ID: modelId,
    PartnerID: enti.PartnerID,
    PartnerName: enti.PartnerName,
    SecretKey: enti.SecretKey,
    CreatedAt: enti.CreatedAt,
    UpdatedAt: enti.UpdatedAt,
    CreatedBy: enti.CreatedBy,
    UpdatedBy: enti.UpdatedBy,
  }
}

func ShopeePartnerDTOToEntity(dto ShopeePartnerDTO) *ShopeePartnerEntity {
  return &ShopeePartnerEntity{
    ID: dto.ID,
    PartnerID: dto.PartnerID,
    PartnerName: dto.PartnerName,
    SecretKey: dto.SecretKey,
    CreatedAt: dto.CreatedAt,
    UpdatedAt: dto.UpdatedAt,
    CreatedBy: dto.CreatedBy,
    UpdatedBy: dto.CreatedBy,
  }
}


type ShopeePartnerRepository interface {
  InitRepository() (error)
  CreateShopeePartner (ctx context.Context,partner ShopeePartnerEntity) (*ShopeePartnerEntity,error)
  GetShopeePartnerByID(ctx context.Context,partner string)  (*ShopeePartnerEntity,error)
  UpdateShopeePartner (ctx context.Context,partner ShopeePartnerEntity)   (*ShopeePartnerEntity,error)
  DeleteShopeePartner (ctx context.Context,partner string) (*ShopeePartnerEntity, error)
}

type shopeePartner struct {
  Logger *zap.Logger
  DB *mongo.Collection
}

func NewShopeePartnerRepository(db *mongo.Collection, log *zap.Logger) ShopeePartnerRepository {
  return &shopeePartner{ Logger: log, DB: db, } }

func (r *shopeePartner)InitRepository() error {
  indexs := []mongo.IndexModel{
    { Keys: bson.D{{Key: "partner_id", Value: 1}}, },
  }
  _, err := r.DB.Indexes().CreateMany(context.TODO(), indexs)
  if err != nil { return errors.New("ShopeePartnerRepository.InitRepository: failed create indexs") }
  r.Logger.Info("ShopeePartnerRepository.InitRepository: index created")
  return nil
}

// func 

func (r *shopeePartner)CreateShopeePartner (ctx context.Context,partner ShopeePartnerEntity) (*ShopeePartnerEntity,error) {

  partnerCreate := ShopeePartnerEntityToModel(partner)
  res,err := r.DB.InsertOne(ctx, partnerCreate)
  if err != nil { return nil, err }

  if oid, ok := res.InsertedID.(bson.ObjectID) ; !ok {
  partnerCreate.ID = oid     }

  partnerParse := ShopeePartnerModelToEntity(*partnerCreate)
  
  return partnerParse,nil
}

func (r *shopeePartner)GetShopeePartnerByID(ctx context.Context,partner string)  (*ShopeePartnerEntity,error){
  var model ShopeePartnerModel
  filter := bson.M{"partner_id":partner }
  err := r.DB.FindOne(ctx, filter).Decode(&model)
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) { 
      return nil,mongo.ErrNoDocuments  
    }
    return nil, err
  }
  resParse := ShopeePartnerModelToEntity(model)
  return resParse, nil
}
func (r *shopeePartner)UpdateShopeePartner (ctx context.Context,partner ShopeePartnerEntity)   (*ShopeePartnerEntity,error) {
  var updated ShopeePartnerModel
  update := ShopeePartnerEntityToModel(partner)

  filter := bson.M{"partner_id": update.PartnerID}
  opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

  err := r.DB.FindOneAndUpdate(ctx, filter,bson.M{"$set":update},opts) .Decode(&updated) 
  if err  != nil {
    if errors.Is(err,mongo.ErrNoDocuments){
      return nil, mongo.ErrNoDocuments
    }
    return nil,err
  }
  
  updatedParse := ShopeePartnerModelToEntity(updated)
  return updatedParse,nil
}
func (r *shopeePartner)DeleteShopeePartner (ctx context.Context,partner string) (*ShopeePartnerEntity, error) {
  var deleted ShopeePartnerModel
  
  filter := bson.M{"partner_id" : partner}
  err := r.DB.FindOneAndDelete(ctx,filter).Decode(&deleted)
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return nil, mongo.ErrNoDocuments
    }
    return nil, err
  }

  deletedParse := ShopeePartnerModelToEntity(deleted)
  return deletedParse, nil
}




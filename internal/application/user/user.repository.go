package user

import (
	"context"
	"ecommerce/internal/infrastructure"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type UserRepository interface {
  InitRepository() error
  CreateUser(ctx context.Context,user *UserEntity) (*UserEntity, error)
  GetUserDetailByUsername(ctx context.Context,id string) (*UserEntity,error)
  UpdateUserDetail(ctx context.Context, user *UserModel) (*UserEntity, error)
}

type userRepo struct {
  logger *zap.Logger
  db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection, log *zap.Logger) UserRepository {
  return &userRepo{db:db, logger: log }
}

func (r *userRepo) InitRepository() error {

  requiredIndexs := []mongo.IndexModel{
    {
      Keys: bson.D{{Key: "username", Value: 1}}, 
      Options: options.Index().SetUnique(true)},
    }

  _, err := r.db.Indexes().CreateMany(context.TODO(), requiredIndexs)
  if err != nil {
    r.logger.Error("error creating index", zap.Error(err))
    return errors.New("UserRepository.InitRepository: failed create index")
  }

  r.logger.Info("UserRepository.InitRepository: index created")
  return nil
}

func (r *userRepo) CreateUser(ctx context.Context,userEntityParam *UserEntity) (*UserEntity, error) {

  user,err := infrastructure.MapStruct[UserEntity, UserModel](*userEntityParam)
  if err != nil {
    return nil ,errors.New("failed to convert Entity_Model")
  }

  var userTarget UserModel
  filter := bson.M{"username": user.Username}
  err = r.db.FindOne(ctx, filter).Decode(&userTarget)

  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      res, err := r.db.InsertOne(context.TODO(), user)
      if err != nil {
        r.logger.Error("Insert failed", zap.String("error", err.Error()))
        return nil, errors.New("failed to insert")
      }

      if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
        user.ID = oid
      }

      resUser,err := infrastructure.MapStruct[UserModel,UserEntity](user)
      if err!=nil {
        return nil, err
      }

      return &resUser,nil
    } else {
      return nil,err
    }
  }

  return nil , errors.New("username already exits")
}

func (r *userRepo) GetUserDetailByUsername(ctx context.Context,id string) (*UserEntity,error) {
  var res UserModel

  filter := bson.M{"username": id}
  err := r.db.FindOne(ctx, filter).Decode(&res)

  if err != nil {
    return nil, errors.New("Username not found")
  }

  resParse,err := infrastructure.MapStruct[UserModel,UserEntity](res)


  return &resParse , nil
}

func (r *userRepo) UpdateUserDetail(ctx context.Context, user *UserModel) (*UserEntity, error) {

  filter := bson.M{"username" : user.Username}
  update := bson.M{
    "$set": bson.M{
      "full_name": user.FullName,
      "avatar_url": user.AvatarURL,
      "status": user.Status,
      "updated_at": time.Now(),
    },
  }

  opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
  var updatedUser UserModel
  err := r.db.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)

  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return nil, errors.New("username not found")
    }
    return nil , err
  }

  updateUserParse ,err := infrastructure.MapStruct[UserModel, UserEntity](updatedUser)

  return &updateUserParse,nil
}

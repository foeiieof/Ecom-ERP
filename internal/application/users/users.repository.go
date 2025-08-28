package users

import (
	"context"
	"ecommerce/internal/pkg"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"

)
type UserRepository interface {
  InitRepository() error
  CreateUser(ctx context.Context,user UserEntity) (*UserEntity, error)
  GetAllUserDetail(ctx context.Context) ([]UserEntity, error)
  GetUserDetailByUsername(ctx context.Context,id string) (*UserEntity,error)
  UpdateUserDetail(ctx context.Context, user UserEntity) (*UserEntity, error)
  DeleteUser(ctx context.Context, user string) (*UserEntity, error)
}

type userRepo struct {
  logger *zap.Logger
  db     *mongo.Collection
}

func NewUserRepository(db *mongo.Collection, log *zap.Logger) UserRepository {
  return &userRepo{db:db, logger: log }
}

func (r *userRepo) InitRepository() error {

  requiredIndexs := []mongo.IndexModel{
    {
      Keys: bson.D{{Key: "username", Value: 1}}, 
      Options: options.Index().SetUnique(true),
    },
    {
      Keys: bson.D{{Key: "email", Value: 1}},
      Options: options.Index().SetUnique(true),
    },
  }

  _, err := r.db.Indexes().CreateMany(context.TODO(), requiredIndexs)
  if err != nil {
    r.logger.Error("error creating index", zap.Error(err))
    return errors.New("UserRepository.InitRepository: failed create index")
  }

  r.logger.Info("UserRepository.InitRepository: index created")
  return nil
}

func (r *userRepo) CreateUser(ctx context.Context,userEntityParam UserEntity) (*UserEntity, error) {

  user,err := pkg.MapStruct[UserEntity, UserModel](userEntityParam)

  r.logger.Info("before:", zap.String("Source", user.ID.Hex()))
  user.ID = bson.NewObjectID()
  r.logger.Info("after:", zap.String("Source", user.ID.Hex()))

  if err != nil {
    return nil ,errors.New("failed to convert Entity_Model")
  }

  var userTarget UserModel
  filter := bson.M{"username": user.Username}
  err = r.db.FindOne(ctx, filter).Decode(&userTarget)
  r.logger.Error("users.repo:CreateUser -- new", zap.Error(err))
  // r.logger.Info("users.repo:CreateUser", zap.String("tarhet : username",userTarget.Username) )
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      res, err := r.db.InsertOne(ctx, user)
      if err != nil {
        r.logger.Error("Insert failed", zap.String("error", err.Error()))
        return nil, errors.New("failed to insert")
      }

      if oid, ok := res.InsertedID.(bson.ObjectID); ok {
        user.ID = oid
      }

      resUser := UserModelToEntity(user)

      return &resUser,nil
    } else {
      return nil,err
    }
  }

  return nil , errors.New("username already exits")
}

func (r *userRepo) GetAllUserDetail(ctx context.Context) ([]UserEntity, error) {

  var res []UserModel
  // filter := bson.M{}

  // var resOne UserModel
  // erro := r.db.FindOne(ctx, bson.M{"username": "foeiieof"}).Decode(&resOne)
  // if erro != nil { 
  //   return nil , erro
  // }
  // r.logger.Info("resEntity", zap.String("string" , resOne.ID.Hex() ) )


  cursor,err := r.db.Find(ctx,bson.M{})
  if err != nil { return nil,err }
  defer cursor.Close(ctx)

  if err := cursor.All(ctx, &res); err != nil {
    // r.logger.Info("resEntity", zap.Error(err) )
    return nil,err
  }

  // resParse,err := pkg.MapSliceStruct[UserModel, UserEntity](res)
 
  resEntity := make([]UserEntity, len(res))
  for i, u := range res {
    // r.logger.Info("resEntity", zap.String("string" , u.ID.Hex() ) )
    resEntity[i] = UserModelToEntity(u)
  }  
  
  return resEntity,nil
}

func (r *userRepo) GetUserDetailByUsername(ctx context.Context,id string) (*UserEntity,error) {

  var res UserModel
  filter := bson.M{"username": id}
  err := r.db.FindOne(ctx, filter).Decode(&res)

  if err != nil {
    return nil, errors.New("Username not found")
  }


  // resParse,err := pkg.MapStruct[UserModel,UserEntity](res)

  resParse := UserModelToEntity(res)

  return &resParse , nil
}

func (r *userRepo) UpdateUserDetail(ctx context.Context, user UserEntity) (*UserEntity, error) {

  filter := bson.M{"username" : user.Username}

  updateFields := bson.M{}
  if user.FullName != "" { updateFields["full_name"] = &user.FullName }
  if user.AvatarURL != "" { updateFields["avatar_url"] = &user.AvatarURL }
  if user.Status != "" { updateFields["status"] = &user.Status }
  if user.LastLoginAt != nil { updateFields["last_login_at"] = &user.LastLoginAt}

  updateFields["is_deleted"] = &user.IsDeleted
  updateFields["updated_at"] = time.Now()

  update := bson.M{ "$set": updateFields }

  opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
  var updatedUser UserModel
  err := r.db.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)

  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      return nil, errors.New("username not found")
    }
    return nil , err
  }

  // updateUserParse ,err := pkg.MapStruct[UserModel, UserEntity](updatedUser)

  updateUserParse := UserModelToEntity(updatedUser)

  return &updateUserParse,nil
}

func (r *userRepo) DeleteUser(ctx context.Context, user string) (*UserEntity, error) {


  filter := bson.M{"username" : user}
  opts := options.FindOneAndDelete()
  var deleteUser UserModel
  err := r.db.FindOneAndDelete(ctx, filter, opts).Decode(&deleteUser)

  if err != nil {
    if errors.Is(err,mongo.ErrNoDocuments) {
        return nil, errors.New("username not found")
    }
    return nil, err
  }

  deleteUserParse,err := pkg.MapStruct[UserModel,UserEntity](deleteUser)

  if err != nil { return nil, errors.New("Error repository.DeleteUser: failed to parse user ")}
  return &deleteUserParse,nil

}

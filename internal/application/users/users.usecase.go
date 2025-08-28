package users

import (
	"context"
	"ecommerce/internal/env"
	"ecommerce/internal/pkg"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
  CreateUser(ctx context.Context, user IReqCreateUserDTO) (*UserDTO,error)
  GetUsers(ctx context.Context) (*[]UserDTO, error)
  GetUserByUsername(ctx context.Context,user string) (*UserDTO,error) 
  UpdateUser(ctx context.Context, userName string,userUpdate IReqUpdateUserDTO) (*UserDTO, error)
  SoftDeleteUserByUsername(ctx context.Context, user string) (*UserDTO,error)
}

type userService struct {
  Config *env.Config
  Logger *zap.Logger

  UserRepository UserRepository
}

func NewUserService(cfg *env.Config, log *zap.Logger, 
  userRepo UserRepository,
) IUserService {
  return &userService{
    Config: cfg,
    Logger: log,
    UserRepository: userRepo,
  }
}

func (s *userService) CreateUser(ctx context.Context, userCreate IReqCreateUserDTO) (*UserDTO,error) {

  if userCreate.Username == "" || userCreate.Email == "" {
    return nil, errors.New("username or email is required")
  }

  passwordHash, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
  if err != nil {
    return nil, err
  }

  // idGenrate := primitive.NewObjectID()
  // s.Logger.Info("id-generate", zap.String("string", idGenrate.Hex()) )

  newUser := &UserEntity{
    // ID: ,
    Username: userCreate.Username,
    Email: userCreate.Email,
    PasswordHash: string(passwordHash),
    Status: "active",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  savedUser, err := s.UserRepository.CreateUser(ctx, *newUser)
  if err != nil {
    // s.Logger.Error("usecase.CreateUser:",zap.Error(err)) 
    return nil,err }

  // -> To DTO
  respDTO,err := pkg.MapStruct[UserEntity,UserDTO](*savedUser)

  if err != nil {
    return nil, errors.New("Error user.Usecase.CreateUser : Convert Entity to DTO")
  }

  return &respDTO,nil
}

func (s *userService)GetUsers(ctx context.Context) (*[]UserDTO, error) {

  resUsers,err := s.UserRepository.GetAllUserDetail(ctx)
  if err != nil { return nil ,err}

  // s.Logger.Info("usecase.GetUser:",)

  resUsersParse ,err := pkg.MapSliceStruct[UserEntity, UserDTO](resUsers)
  if err != nil { return nil, errors.New("Error usecase.GetUsers:parse to UserDTO")}

  return &resUsersParse,nil
}

func (s *userService)GetUserByUsername(ctx context.Context,user string) (*UserDTO,error) {
  if user == "" {
    return nil, errors.New("username is required")
  }

  resUser,err := s.UserRepository.GetUserDetailByUsername(ctx, user)
  if err != nil { return nil, err }

  resUserParse,err := pkg.MapStruct[UserEntity,UserDTO](*resUser)

 return &resUserParse,nil 
}

func (s *userService)UpdateUser(ctx context.Context, userName string,userUpdate IReqUpdateUserDTO) (*UserDTO, error) {

  if userUpdate.Username == "" {
    return nil, errors.New("username is required") 
  }

  // Check username
  _,err := s.UserRepository.GetUserDetailByUsername(ctx, userName)
  if err != nil { return nil, err }

  userUpdateParse,err := pkg.MapStruct[IReqUpdateUserDTO, UserEntity](userUpdate)
  userUpdateParse.Username = userName

  resUser,err := s.UserRepository.UpdateUserDetail(ctx,userUpdateParse)
  if err != nil {return nil, err}


  resUserParse,err := pkg.MapStruct[UserEntity, UserDTO](*resUser)
  if err != nil { return nil, errors.New("Error usecase.UpdateUser: parse to DTO")}

  return &resUserParse,nil
} 

func (s *userService)SoftDeleteUserByUsername(ctx context.Context, user string) (*UserDTO,error) {
  if user == "" {
    return nil ,errors.New("Error repository.DeleteUserByUsername: username is required") 
  }

  // Check user
  _,err := s.UserRepository.GetUserDetailByUsername(ctx,user)
  if err != nil { return nil, err}

  var softDelete = &UserEntity{ Username: user, IsDeleted: true , Status: StatusInactive}

  userDeleted, err := s.UserRepository.UpdateUserDetail(ctx, *softDelete)
  if err != nil {
    return nil, err
  }

  // userDeleted,err := s.UserRepository.DeleteUser(ctx, user) 
  // if err != nil { return nil, err}

  userDeletedParse,err := pkg.MapStruct[UserEntity,UserDTO](*userDeleted)
  return &userDeletedParse,nil 
}

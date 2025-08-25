package user

import (
	"context"
	"ecommerce/internal/env"
	"ecommerce/internal/infrastructure"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
  CreateUser(ctx context.Context, user IReqCreateUserDTO) (*UserDTO,error)
  GetUserByUsername(ctx context.Context,user string) (*UserDTO,error) 

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

  newUser := &UserEntity{
    ID: uuid.NewString(),
    Username: userCreate.Username,
    Email: userCreate.Email,
    PasswordHash: string(passwordHash),
    Status: "active",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  savedUser, err := s.UserRepository.CreateUser(ctx, newUser)
  if err != nil { return nil,err }

  // -> To DTO
  respDTO,err := infrastructure.MapStruct[UserEntity,UserDTO](*savedUser)

  if err != nil {
    return nil, errors.New("Error user.Usecase.CreateUser : Convert Entity to DTO")
  }

  return &respDTO,nil
}

func (s *userService)GetUserByUsername(ctx context.Context,user string) (*UserDTO,error) {
  if user == "" {
    return nil, errors.New("username is required")
  }

  resUser,err := s.UserRepository.GetUserDetailByUsername(ctx, user)
  if err != nil { return nil, err }

  resUserParse,err := infrastructure.MapStruct[UserEntity,UserDTO](*resUser)


 return &resUserParse,nil 

}



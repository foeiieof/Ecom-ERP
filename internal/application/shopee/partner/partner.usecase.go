package partner

import (
	"context"
	"ecommerce/internal/env"
	"time"

	"go.uber.org/zap"
)

type IShopeePartnerService interface {
  AddShopeePartner(ctx context.Context,dto *IReqShopeePartnerDTO) (*ShopeePartnerDTO,error)
  GetShopeePartnerByID(ctx context.Context, partner string)       (*ShopeePartnerDTO,error)
  GetAllShopeePartner(ctx context.Context) ([]ShopeePartnerDTO,error)
  UpdateShopeePartner(ctx context.Context, dto *IReqShopeePartnerDTO) (*ShopeePartnerDTO,error)
  DeleteShopeePartnerByID(ctx context.Context, partner string) (*ShopeePartnerDTO,error)
}

type shopeePartnerService struct {
  Config *env.Config
  Logger *zap.Logger

  ShopeePartnerRepository ShopeePartnerRepository
}

func NewShopeePartnerService(cfg *env.Config, log *zap.Logger,
  shopeePartner ShopeePartnerRepository,
) IShopeePartnerService {
  return &shopeePartnerService{
    Config: cfg,
    Logger: log,
    ShopeePartnerRepository: shopeePartner,
  }
}

func (s *shopeePartnerService)AddShopeePartner(ctx context.Context,dto *IReqShopeePartnerDTO) (*ShopeePartnerDTO,error) {
  // username :=  
  entities := &ShopeePartnerEntity{
    PartnerID: dto.PartnerID,
    PartnerName: dto.PartnerName,
    SecretKey: dto.SecretKey,
    Validate: false,
    CreatedAt: time.Now(),
    CreatedBy: *dto.Username,
    UpdatedAt: time.Now(),
    UpdatedBy: *dto.Username,
  }
  add,err := s.ShopeePartnerRepository.CreateShopeePartner(ctx, entities) 
  if err != nil { return nil,err }

  addParse := ShopeePartnerEntityToDTO(*add)
  return addParse,nil
}

func (s *shopeePartnerService)GetShopeePartnerByID(ctx context.Context, partner string) (*ShopeePartnerDTO, error) {
  
  partnerObject,err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx,partner)
  if err != nil { return nil, err }

  partnerParse := ShopeePartnerEntityToDTO(*partnerObject)

  return partnerParse,nil
}

func (s *shopeePartnerService)GetAllShopeePartner(ctx context.Context) ([]ShopeePartnerDTO,error) {

  objectPartner, err := s.ShopeePartnerRepository.GetAllShopeePartner(ctx)
  if err != nil { return nil, err}


  objectParse  := make([]ShopeePartnerDTO, len(objectPartner))

  for i,u := range objectPartner {
    objectParse[i] = *ShopeePartnerEntityToDTO(u)
  } 
  
  return objectParse, nil
}

func (s *shopeePartnerService)UpdateShopeePartner(ctx context.Context, dto *IReqShopeePartnerDTO) (*ShopeePartnerDTO,error) {
  
  partnerBef,err := s.ShopeePartnerRepository.GetShopeePartnerByID(ctx, dto.PartnerID)
  if err != nil { return nil, err }

  if dto.PartnerName != "" {
    partnerBef.PartnerName = dto.PartnerName 
  }

  if dto.SecretKey != "" {
    partnerBef.SecretKey = dto.SecretKey
  }

  if dto.Username != nil {
    s.Logger.Info("usecase.UpdateShopeePartner", zap.String("params", *dto.Username))
    partnerBef.UpdatedBy = *dto.Username
  }

  partnerUpdated, err := s.ShopeePartnerRepository.UpdateShopeePartner(ctx,partnerBef)
  if err != nil { return nil,err }

  partnerUpdatedParse := ShopeePartnerEntityToDTO(*partnerUpdated)

  return partnerUpdatedParse, nil
}

func (s *shopeePartnerService)DeleteShopeePartnerByID(ctx context.Context, partner string) (*ShopeePartnerDTO,error) {

  deleted, err := s.ShopeePartnerRepository.DeleteShopeePartner(ctx, partner)
  if err != nil  { return nil, err }

  deletedParse := ShopeePartnerEntityToDTO(*deleted)

  return deletedParse, nil
}



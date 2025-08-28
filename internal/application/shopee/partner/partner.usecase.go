package partner

import (
	"context"
	"ecommerce/internal/env"

	"go.uber.org/zap"
)

type IShopeePartnerService interface {
  AddShopeePartner(ctx context.Context,dto ShopeePartnerDTO) (*ShopeePartnerDTO,error)
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

func (s *shopeePartnerService)AddShopeePartner(ctx context.Context,dto ShopeePartnerDTO) (*ShopeePartnerDTO,error) {

  entities := ShopeePartnerDTOToEntity(dto)
  add,err := s.ShopeePartnerRepository.CreateShopeePartner(ctx, *entities) 
  if err != nil { return nil,err }

  addParse := ShopeePartnerEntityToDTO(*add)
  return addParse,nil
}





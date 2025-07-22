package repository

import (
	"ecommerce/internal/application/shopee"
)

type IMongoCollectionRepository interface {
	ShopeeAuthCollection() shopee.ShopeeAuthRepository
  ShopeeAuthRequestCollection() shopee.ShopeeAuthRequestRepository
  ShopeePartnerCollection() shopee.ShopeePartnerRepository
}

type mongoCollectionRepository struct {
	shopeeAuthRepo shopee.ShopeeAuthRepository
  shopeeAuthRequestRepo shopee.ShopeeAuthRequestRepository
  shopeePartnerRepo shopee.ShopeePartnerRepository
}

func NewMongoCollectionRepository(
	shopeeAuth shopee.ShopeeAuthRepository,
  shopeeAuthReq shopee.ShopeeAuthRequestRepository,
  shopeePartner shopee.ShopeePartnerRepository,
  // logger *zap.Logger, cfg *env.Config,
) IMongoCollectionRepository {
	return &mongoCollectionRepository{
		shopeeAuthRepo: shopeeAuth,
    shopeeAuthRequestRepo: shopeeAuthReq,
    shopeePartnerRepo: shopeePartner,
	}
}

func (m *mongoCollectionRepository) ShopeeAuthCollection() shopee.ShopeeAuthRepository {
	return m.shopeeAuthRepo
}

func (m *mongoCollectionRepository) ShopeeAuthRequestCollection() shopee.ShopeeAuthRequestRepository {
  return m.shopeeAuthRequestRepo
}

func (m *mongoCollectionRepository) ShopeePartnerCollection() shopee.ShopeePartnerRepository {
  return m.shopeePartnerRepo
}


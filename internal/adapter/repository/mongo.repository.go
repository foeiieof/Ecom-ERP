package repository

import (
	"ecommerce/internal/application/shopee"
	"ecommerce/internal/application/shopee/partner"
	"ecommerce/internal/application/users"
)

type IMongoCollectionRepository interface {
	ShopeeAuthCollection() shopee.ShopeeAuthRepository
  ShopeeAuthRequestCollection() shopee.ShopeeAuthRequestRepository
  ShopeePartnerCollection() partner.ShopeePartnerRepository
  UsersCollection() users.UserRepository
}

type mongoCollectionRepository struct {
	shopeeAuthRepo shopee.ShopeeAuthRepository
  shopeeAuthRequestRepo shopee.ShopeeAuthRequestRepository
  shopeePartnerRepo partner.ShopeePartnerRepository
  userRepo users.UserRepository
}

func NewMongoCollectionRepository(
	shopeeAuth shopee.ShopeeAuthRepository,
  shopeeAuthReq shopee.ShopeeAuthRequestRepository,
  shopeePartner partner.ShopeePartnerRepository,
  users users.UserRepository,
  // logger *zap.Logger, cfg *env.Config,
) IMongoCollectionRepository {
	return &mongoCollectionRepository{
		shopeeAuthRepo: shopeeAuth,
    shopeeAuthRequestRepo: shopeeAuthReq,
    shopeePartnerRepo: shopeePartner,
    userRepo: users,
	}
}

func (m *mongoCollectionRepository) ShopeeAuthCollection() shopee.ShopeeAuthRepository {
	return m.shopeeAuthRepo
}

func (m *mongoCollectionRepository) ShopeeAuthRequestCollection() shopee.ShopeeAuthRequestRepository {
  return m.shopeeAuthRequestRepo
}

func (m *mongoCollectionRepository) ShopeePartnerCollection() partner.ShopeePartnerRepository{
  return m.shopeePartnerRepo
}

func (m *mongoCollectionRepository) UsersCollection() users.UserRepository {
  return m.userRepo
}


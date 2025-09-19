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
  ShopeeShopCollection() shopee.ShopeeShopDetailsRepository
  ShopeeOrderCollection() shopee.ShopeeOrderRepository
}

type mongoCollectionRepository struct {
	shopeeAuthRepo shopee.ShopeeAuthRepository
  shopeeAuthRequestRepo shopee.ShopeeAuthRequestRepository
  shopeePartnerRepo partner.ShopeePartnerRepository
  userRepo users.UserRepository
  shopeeShopRepo shopee.ShopeeShopDetailsRepository
  shopeeOrderRepo shopee.ShopeeOrderRepository
}

func NewMongoCollectionRepository(
	shopeeAuth shopee.ShopeeAuthRepository,
  shopeeAuthReq shopee.ShopeeAuthRequestRepository,
  shopeePartner partner.ShopeePartnerRepository,
  users users.UserRepository,
  shop shopee.ShopeeShopDetailsRepository,
  shopeeOrder shopee.ShopeeOrderRepository,
  // logger *zap.Logger, cfg *env.Config,
) IMongoCollectionRepository {
	return &mongoCollectionRepository{
		shopeeAuthRepo: shopeeAuth,
    shopeeAuthRequestRepo: shopeeAuthReq,
    shopeePartnerRepo: shopeePartner,
    userRepo: users,
    shopeeShopRepo: shop,
    shopeeOrderRepo: shopeeOrder,
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

func (m *mongoCollectionRepository) ShopeeShopCollection() shopee.ShopeeShopDetailsRepository {
  return m.shopeeShopRepo
}
func (m *mongoCollectionRepository) ShopeeOrderCollection() shopee.ShopeeOrderRepository{
  return m.shopeeOrderRepo
}

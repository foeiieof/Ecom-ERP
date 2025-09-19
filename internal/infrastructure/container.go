package infrastructure

import (
	"context"
	"ecommerce/internal/adapter"
	"ecommerce/internal/adapter/repository"
	"ecommerce/internal/application/auth"
	"ecommerce/internal/application/demo"
	"ecommerce/internal/application/health"
	"ecommerce/internal/application/shopee"
	"ecommerce/internal/application/shopee/partner"
	"ecommerce/internal/application/users"
	"ecommerce/internal/delivery/http/handler"
	"ecommerce/internal/delivery/http/middleware"

	"ecommerce/internal/application/swagger"
	"ecommerce/internal/env"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type MiddlewareHandle struct {
	Log     *middleware.LogHandler
	Error   *middleware.ErrorHandler
  Auth    middleware.IAuthMiddleware
	Shopee  *middleware.ShopeeMiddleware
}

type Repositories struct {
	MongoRepository repository.IMongoCollectionRepository
}

type Adapter struct {
	ShopeeAdapter adapter.IShopeeService
}

// Container holds all dependencies
type Container struct {
	Config      *env.Config
	Logger      *zap.Logger
	Valid       *validator.Validate
	MongoClient *mongo.Client

	Repository *Repositories
	Middleware *MiddlewareHandle
	Adapter    *Adapter
}

func NewContainer(cfg *env.Config, mongo *mongo.Client, logger *zap.Logger, valid *validator.Validate) *Container {
	return &Container{
		Config:      cfg,
		MongoClient: mongo,
		Logger:      logger,
		Valid:       valid,
	}
}

func (c *Container) InitRepositories() {
	// c.MongoClient = mongoClient

  // for DB name : auth
	auth := c.MongoClient.Database(c.Config.DB.ConfigDBAuthName)
  db := c.MongoClient.Database(c.Config.DB.ConfigDBName)

	shopeePartnerCollection := auth.Collection("shopee_partner")
	shopeePartner := partner.NewShopeePartnerRepository(shopeePartnerCollection, c.Logger)
	shopeePartner.InitRepository()

	shopeeAuthCollection := auth.Collection("shopee_shop_auth")
	shopeeAuth := shopee.NewShopeeAuthRepository(shopeeAuthCollection, c.Logger)
	shopeeAuth.InitRepository()

	shopeeAuthReqCollection := auth.Collection("shopee_auth_request")
	shopeeAuthReq := shopee.NewShopeeAuthRequestRepository(shopeeAuthReqCollection, c.Logger)
	shopeeAuthReq.InitRepository()

  userCollection := db.Collection("users")
  userReq := users.NewUserRepository(userCollection, c.Logger)
  userReq.InitRepository()

  shopeeShopCollection := db.Collection("shopee_shop")
  shopeeShop := shopee.NewShopeeShopDetailsRepository(shopeeShopCollection, c.Logger)
  shopeeShop.InitRepository()

  ShopeeOrderCollection := db.Collection("shopee_order")
  shopeeOrder := shopee.NewShopeeOrderRepository(ShopeeOrderCollection, c.Logger)
  shopeeOrder.InitRepository()

	c.Repository = &Repositories{
		MongoRepository: repository.NewMongoCollectionRepository(shopeeAuth, shopeeAuthReq, shopeePartner,userReq, shopeeShop, shopeeOrder),
	}
  // next using in handle()
}

// initHandlers initializes HTTP handlers

// initMiddleware initializes middleware
func (c *Container) InitMiddleware() {

	logMiddleware := middleware.NewLogHandler(c.Logger, fiberLog.New())
	errorMiddleware := middleware.NewErrorHandler(c.Logger)

  authMiddleware := middleware.NewAuthMiddleware(c.Config, c.Logger)

	db := c.MongoClient.Database(c.Config.DB.ConfigDBName)
	shopeeCollection := db.Collection("shopee_auth")
	shopeeAuthCollection := shopee.NewShopeeAuthRepository(shopeeCollection, c.Logger)

	shopeeMiddleware := middleware.NewShopeeMiddleware(c.Logger, shopeeAuthCollection)

	c.Middleware = &MiddlewareHandle{
		Log:    logMiddleware,
		Error:  errorMiddleware,
		Shopee: shopeeMiddleware,
    Auth:   authMiddleware,
	}
}

func (c *Container) InitHandlers(g fiber.Router) {

	// db := mongo.Connect("...")
	// userRepository := repository.NewUserRepository(db)
	// userService := service.NewUserService(userRepository)
	// userHandler := handler.NewUserHandler(userService)

	health := health.NewHealthHandler(c.Logger)

	swagger := swagger.NewSwaggerHandler()

	demo := demo.NewDemoHandler(c.Repository.MongoRepository)

  // repository
	shopeeRepo := c.Repository.MongoRepository.ShopeeAuthCollection()
	shopeeReqRepo := c.Repository.MongoRepository.ShopeeAuthRequestCollection()
	shopeePartnerRepo := c.Repository.MongoRepository.ShopeePartnerCollection()
  userRepo := c.Repository.MongoRepository.UsersCollection()
  shopeeShopRepo := c.Repository.MongoRepository.ShopeeShopCollection()
  shopeeOrderRepo := c.Repository.MongoRepository.ShopeeOrderCollection()
  // waiting
  // shopeeShopRepo := c.Repository.MongoRepository.ShopeePartnerCollection()
  //usecase
  shopeePartnerUsecase := partner.NewShopeePartnerService(c.Config, c.Logger, shopeePartnerRepo)
	shopeeUsecase := shopee.NewShopeeService(c.Config, c.Logger, c.Adapter.ShopeeAdapter, shopeeRepo, shopeeReqRepo, shopeePartnerRepo, shopeeShopRepo, shopeeOrderRepo)
  usersUsecase := users.NewUserService(c.Config,c.Logger,userRepo)
  authUsecase := auth.NewAuthService(c.Config,c.Logger,userRepo)

  // shopeeShop := shopee.NewShopeeShopDetailsService () 
  // handler
	shopee := shopee.NewShopeeHandler(shopeeUsecase, shopeePartnerUsecase,c.Logger, c.Valid)
  shopeePartner := partner.NewShopeePartnerHandler(c.Logger, c.Valid,shopeePartnerUsecase)
  users := users.NewUserHandler(usersUsecase,c.Logger, c.Valid)
  auth := auth.NewAuthHandle(c.Config,authUsecase, c.Logger, c.Valid)

	h := handler.NewRouterHandler(
    c.Middleware.Auth.Handler(),
    c.Middleware.Shopee.Handler(),
    health, swagger, demo, shopee, shopeePartner,auth,users)
	h.RegisterHandlers(g)
}

func (c *Container) InitAdapter() {
	shopeeAdapter := adapter.NewShopeeAPI(c.Config,c.Config.Shopee.ShopeeApiBaseUrl, c.Config.Shopee.ShopeeApiBasePrefix, c.Logger)
	c.Adapter = &Adapter{ShopeeAdapter: shopeeAdapter}
}

// Close cleans up resources
func (c *Container) Close() {

	if c.MongoClient != nil {
		if err := c.MongoClient.Disconnect(context.TODO()); err != nil {
			c.Logger.Error("Failed to disconnect Mongo", zap.Error(err))
		}
	}
}

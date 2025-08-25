package infrastructure

import (
	"context"
	"ecommerce/internal/adapter"
	"ecommerce/internal/adapter/repository"
	"ecommerce/internal/application/demo"
	"ecommerce/internal/application/health"
	"ecommerce/internal/application/shopee"
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
	Log    *middleware.LogHandler
	Error  *middleware.ErrorHandler
	Shopee *middleware.ShopeeMiddleware
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

	db := c.MongoClient.Database(c.Config.DB.ConfigDBName)

	shopeePartnerCollection := db.Collection("shopee_partner")
	shopeePartner := shopee.NewShopeePartnerRepository(shopeePartnerCollection, c.Logger)
	shopeePartner.InitRepository()

	shopeeAuthCollection := db.Collection("shopee_shop_auth")
	shopeeAuth := shopee.NewShopeeAuthRepository(shopeeAuthCollection, c.Logger)
	shopeeAuth.InitRepository()

	shopeeAuthReqCollection := db.Collection("shopee_auth_request")
	shopeeAuthReq := shopee.NewShopeeAuthRequestRepository(shopeeAuthReqCollection, c.Logger)
	shopeeAuthReq.InitRepository()

	c.Repository = &Repositories{
		MongoRepository: repository.NewMongoCollectionRepository(shopeeAuth, shopeeAuthReq, shopeePartner),
	}
}

// initHandlers initializes HTTP handlers

// initMiddleware initializes middleware
func (c *Container) InitMiddleware() {

	logMiddleware := middleware.NewLogHandler(c.Logger, fiberLog.New())
	errorMiddleware := middleware.NewErrorHandler(c.Logger)

	db := c.MongoClient.Database(c.Config.DB.ConfigDBName)
	shopeeCollection := db.Collection("shopee_auth")
	shopeeAuthCollection := shopee.NewShopeeAuthRepository(shopeeCollection, c.Logger)

	shopeeMiddleware := middleware.NewShopeeMiddleware(c.Logger, shopeeAuthCollection)

	c.Middleware = &MiddlewareHandle{
		Log:    logMiddleware,
		Error:  errorMiddleware,
		Shopee: shopeeMiddleware,
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

  // repositori
	shopeeRepo := c.Repository.MongoRepository.ShopeeAuthCollection()
	shopeeReqRepo := c.Repository.MongoRepository.ShopeeAuthRequestCollection()
	shopeePartnerRepo := c.Repository.MongoRepository.ShopeePartnerCollection()

	shopeeUsecase := shopee.NewShopeeService(c.Config, c.Logger, c.Adapter.ShopeeAdapter, shopeeRepo, shopeeReqRepo, shopeePartnerRepo)
	shopee := shopee.NewShopeeHandler(shopeeUsecase, c.Logger, c.Valid)



	h := handler.NewRouterHandler(health, swagger, demo, shopee)
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

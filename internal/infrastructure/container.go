package infrastructure

import (
	"context"
	"ecommerce/internal/adapter/repository"
	"ecommerce/internal/application/demo"
	"ecommerce/internal/application/health"
	"ecommerce/internal/application/shopee"

	"ecommerce/internal/application/swagger"
	"ecommerce/internal/delivery/http/handler"
	"ecommerce/internal/delivery/http/middleware"
	"ecommerce/internal/env"

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
	MongoRepository *repository.MongoCollectionRepository
}

// Container holds all dependencies
type Container struct {
	Config      *env.Config
	Logger      *zap.Logger
	MongoClient *mongo.Client

	Repository *Repositories
	Middleware *MiddlewareHandle
}

func NewContainer(cfg *env.Config, mongo *mongo.Client, logger *zap.Logger) *Container {
	return &Container{
		Config:      cfg,
		MongoClient: mongo,
		Logger:      logger,
	}
}

func (c *Container) InitRepositories() {
	// c.MongoClient = mongoClient

	db := c.MongoClient.Database(c.Config.DB.ConfigDBName)

	shopeeCollection := db.Collection("shopee_auth")
	shopeeAuthCollection := shopee.NewShopeeAuthRepository(shopeeCollection, c.Logger)

	c.Repository = &Repositories{
		MongoRepository: repository.NewMongoCollectionRepository(shopeeAuthCollection, c.Logger, c.Config),
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



  shopeeRepo := c.Repository.MongoRepository.ShopeeAuthCollection
  shopeeUsecase := shopee.NewShopeeService(shopeeRepo ,c.Logger, c.Config)
  shopee := shopee.NewShopeeHandler(shopeeUsecase, c.Logger)

	h := handler.NewRouterHandler(health, swagger, demo, shopee)
	h.RegisterHandlers(g)

}

// Close cleans up resources
func (c *Container) Close() {

	if c.MongoClient != nil {
		if err := c.MongoClient.Disconnect(context.TODO()); err != nil {
			c.Logger.Error("Failed to disconnect Mongo", zap.Error(err))
		}
	}
}

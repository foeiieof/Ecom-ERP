package main

import (
	"context"
	"log"
	"os"


	"ecommerce/internal/env"
	"ecommerce/internal/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	_ "ecommerce/internal/docs"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func main() {
  envSet := os.Getenv("ENV")
	if envSet == "" { envSet = "dev" }

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg, err := env.LoadEnv(envSet, logger)
	if err != nil { log.Fatal("Failed to load configuration:", err) }
// Demoinstant - Mongo
  var mongoDriver infrastructure.MongoDriverMethod = infrastructure.NewMongoClient(logger)  
  mongoClient,err := mongoDriver.Connect(cfg)
  if err != nil { logger.Error("Mongo Error:", zap.Error(err)) }
  defer mongoDriver.Disconnect(mongoClient)

	container := infrastructure.NewContainer(cfg, mongoClient,logger)

  container.InitMiddleware()
  middlewareConfig := container.Middleware

  container.InitRepositories() 

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok { code = e.Code }

			logger.Error("Error:", zap.Error(err))

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		DisableStartupMessage: false,
	})


	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	if envSet == "dev" { app.Use(middlewareConfig.Log.ReqLogOriginal())
	} else { app.Use(middlewareConfig.Log.ReqLogZap()) }


  // app init shopee middleware

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to E-commerce API",
			"version": cfg.Server.AppVersion,
			"status":  "running",
      "environment" : cfg.Server.AppEnv,
		})
	})

  app.Get("/mongo-check" , func (c *fiber.Ctx) error {
    if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
      return c.JSON(fiber.Map{
        "message": "MongoDB is not running",
      })
    }
    return c.JSON(fiber.Map{
      "message": "MongoDB is running",
    })
  })


  // Shopee Middleware
  app.Use(middlewareConfig.Shopee.Handler())

  // Prfix declareration
	api := app.Group(cfg.Server.Prefix)

	// Router
	// health  := handler.NewHealthHandler(logger)
	// swagger := handler.NewSwaggerHandler()
	// demo    := handler.NewDemoHandler(container.Repository.MongoRepository)

	// Centralization Router
	// router := handler.NewRouterHandler(health, swagger, demo)
	// router.RegisterHandlers(api)

  container.InitHandlers(api) 

	// app.Use(container.OAuthMiddleware.Handler())
	// Auth routes
	// auth := app.Group("/auth")
	// auth.Post("/register", container.AuthHandler.Register)
	// auth.Post("/login", container.AuthHandler.Login)
	// auth.Post("/refresh", container.AuthHandler.RefreshToken)
	// auth.Post("/logout", container.AuthHandler.Logout)

	// Protected routes (authentication required)
	// protected := app.Group("/api")
	// protected.Use(container.OAuthMiddleware.RequireAuth()) // Strict authentication

	// protected.Get("/profile", container.AuthHandler.Profile)

	// Setup graceful shutdown
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// go func() {
	// 	<-c
	// 	fmt.Println("Gracefully shutting down...")
	// 	_ = app.Shutdown()
	// }()

	app.Use(middlewareConfig.Error.NotFoundHandler())

	// Start server
	address := cfg.Server.Host + cfg.Server.Port
	if err := app.Listen(address); err != nil { log.Fatal("Server failed to start:", err) }
}

// validateConfig validates required configuration values
// func validateConfig(cfg *env.Config) error {
// 	if cfg.JWT.AuthJWTSecretKey == "your-super-secret-jwt-key-change-this-in-production" {
// 		return fmt.Errorf("JWT secret key must be changed from default value")
// 	}

// 	return nil
// }

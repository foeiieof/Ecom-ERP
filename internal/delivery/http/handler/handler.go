package handler

import (
	"ecommerce/internal/application/demo"
	"ecommerce/internal/application/health"
	"ecommerce/internal/application/shopee"
	"ecommerce/internal/application/swagger"

	"github.com/gofiber/fiber/v2"
)

type RouterHandler struct {
	healthHandler  *health.HealthHandler
	swaggerHandler *swagger.SwaggerHandler
	demoHandler    *demo.DemoHandler
	shopeeHandler  shopee.IShopeeHandler
}

func NewRouterHandler(
	health *health.HealthHandler,
	swagger *swagger.SwaggerHandler,
	demo *demo.DemoHandler,
	shopee shopee.IShopeeHandler,
) *RouterHandler {
	return &RouterHandler{
		healthHandler:  health,
		swaggerHandler: swagger,
		demoHandler:    demo,
		shopeeHandler:  shopee,
	}
}

// SWAGGER : init
// swag init --parseDependency -g internal/delivery/http/handler/handler.go -o internal/docs

func (r *RouterHandler) RegisterHandlers(router fiber.Router) {
	health := router.Group("/health")
	health.Get("/", r.healthHandler.HealthCheck)

	swagger := router.Group("/swagger")
	swagger.Get("/*", r.swaggerHandler.SwaggerIndex)

	demo := router.Group("/demo")
	demo.Get("/", r.demoHandler.DemoCheck)

  // Shopee Handle
	shopee := router.Group("/shopee")
	shopee.Get("/", r.shopeeHandler.GetShopeeAuthByShopId)

	// Generate Link for Auth + add to DB
	// shopee.Post("/shop/add_auth_partner", r.shopeeHandler.PostShopAuthPartner)
  shopee.Post("/shop/auth_partner", r.shopeeHandler.PostShopAuthPartner)
  shopee.Post("/shop/auth_token",  r.shopeeHandler.PostShopeeTokenAuthPartner)

  shopee.Get("/shop/auth_token/:shopeeShopID",  r.shopeeHandler.GetShopeeTokenAuthPartnerByShopId)

  // webhook - auth
  shopee.Get("/webhook/auth_partner/:partnerId", r.shopeeHandler.GetWebHookAuthPartner)

}

// Prototype :
// shopee.Post("/shop/auth_partner",  func(c *fiber.Ctx) error { return c.SendString("OK")} )



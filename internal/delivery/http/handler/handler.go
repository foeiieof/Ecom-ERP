package handler

import (
	"ecommerce/internal/application/auth"
	"ecommerce/internal/application/demo"
	"ecommerce/internal/application/health"
	"ecommerce/internal/application/shopee"
	"ecommerce/internal/application/shopee/partner"
	"ecommerce/internal/application/swagger"
	"ecommerce/internal/application/users"

	"github.com/gofiber/fiber/v2"
)

type RouterHandler struct {
  callback       fiber.Handler
  shopeeMiddleware fiber.Handler
	healthHandler  health.HealthHandler
	swaggerHandler swagger.SwaggerHandler
	demoHandler    demo.DemoHandler
	shopeeHandler  shopee.IShopeeHandler
  partnerHandler partner.IShopeePartnerHandler
  authHandler    auth.AuthHandler
  usersHandle    users.IUserHandler 
  // userHandle     user.IUserHandler
}

func NewRouterHandler(
  fn fiber.Handler,
  shop  fiber.Handler,

	health  health.HealthHandler,
	swagger swagger.SwaggerHandler,
	demo    demo.DemoHandler,
	shopee  shopee.IShopeeHandler,
  partner partner.IShopeePartnerHandler,
  auth    auth.AuthHandler,
  user    users.IUserHandler, // user *user.
) *RouterHandler {
	return &RouterHandler{
    callback: fn,
    shopeeMiddleware: shop,

		healthHandler:  health,
		swaggerHandler: swagger,
		demoHandler:    demo, shopeeHandler:  shopee,
    partnerHandler: partner,
    authHandler: auth,
    usersHandle: user,
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

  // Auth Handler 
  auth := router.Group("/auth")
  // auth.Get("/", r.authHandler.CheckAuth)
  auth.Post("/login", r.authHandler.PostUserAuthLogin )
  auth.Post("/refresh", r.authHandler.PostUserAuthRefresh )
  // auth/refresh
  // auth/logout
  // auth/register
  // auth/me

  // auth.Get("/", func(c *fiber.Ctx)error {return c.SendString("ok!")} )

  me := router.Group("/user",r.callback)
  me.Get("/me", r.usersHandle.GetUserMe)

  user := router.Group("/users")
  user.Get("/", r.usersHandle.GetUsers)
  user.Get("/:userId", r.usersHandle.GetUserByID)
  user.Post("/", r.usersHandle.CreateUser)
  user.Patch("/:userId", r.usersHandle.UpdateUserByID) 
  user.Delete("/:userId", r.usersHandle.DeleteUserByID)

  // Shopee Handle
	shopee := router.Group("/shopee", r.callback)
	// shopee.Get("/", r.shopeeHandler.GetShopeeAuthByShopId)

	// Generate Link for Auth + add to DB
	// shopee.Post("/shop/add_auth_partner", r.shopeeHandler.PostShopAuthPartner)
  shopee.Post("/shop/auth_partner", r.shopeeHandler.PostShopAuthPartner)
  shopee.Post("/shop/auth_token",  r.shopeeHandler.PostShopeeTokenAuthPartnerWithCode)

  shopee.Get("/shop/auth_token/:shopeeShopID",  r.shopeeHandler.GetShopeeTokenAuthPartnerByShopId)

  // webhook - auth
  shopee.Get("/webhook/auth_partner/:partnerId", r.shopeeHandler.GetWebHookAuthPartner)
   


  partner := shopee.Group("/partner")
  // Crud Shopee Partner
  partner.Post("/", r.partnerHandler.CreateShopeePartner)
  partner.Get("/", r.partnerHandler.GetAllShopeePartner)
  partner.Get("/:partnerID", r.partnerHandler.GetShopeePartnerByID)
  partner.Patch("/:partnerID", r.partnerHandler.UpdateShopeePartnerByID)
  partner.Delete("/:partnerID", r.partnerHandler.DeleteShopeePartnerByID)

  partner.Get("/:partnerID/shops", r.shopeeHandler.GetShopeeShopListByPartnerID)

  // new webhook - auth 
  // to send code and shop id to request asccess and refresh from Shopee
  partner.Get("/:partnerID/webhook",r.shopeeHandler.GetWebHookAuthPartner, r.shopeeMiddleware )

  // waiting to update struct 
  // --> to partner check all shop is under manage
  // shopee.Get("/shop_list/:partnerID", r.shopeeHandler.GetShopeeShopListByPartnerID )
  
  // waiting 
  // shopee.Get("/partner/shop_detail/:shopID", func(c *fiber.Ctx) error { return c.SendString("OK")})

  shopee.Get("/shop/:shopeeShopID/details", r.shopeeHandler.GetShopeeShopDetails )

  // |----> shopee.Get("/shop/order_list/:shopeeShopID", r.shopeeHandler.GetShopeeOrderListByShopID )
  shopee.Get("/shop/:shopeeShopID/orders", r.shopeeHandler.GetShopeeOrderListByShopID )
  
  // |----> shopee.Get("/shop/order_detail/:shopeeShopID/:orderSN", )
  shopee.Get("/shop/:shopeeShopID/orders/:orderSN", r.shopeeHandler.GetShopeeOrderDetailsByShopIDAndOrderSN )

  // shoperPartner := router.Group("/shopee-partner")
  // shoperPartner.Get("/", func(c *fiber.Ctx) error { return c.SendString("ok !")} )

}

// Prototype :
// shopee.Post("/shop/auth_partner",  func(c *fiber.Ctx) error { return c.SendString("OK")} )



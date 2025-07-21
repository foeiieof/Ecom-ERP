package middleware

import (
	"ecommerce/internal/application/shopee"

	"ecommerce/internal/env"
	"strings"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
)

type ShopeeRequest struct {
	ShopId string `json:"shop_id"`
}

type ShopeeMiddleware struct {
	logger               *zap.Logger
	shopeeAuthCollection shopee.ShopeeAuthRepository
	configApp            *env.Config
}

func NewShopeeMiddleware(log *zap.Logger, shopee shopee.ShopeeAuthRepository) *ShopeeMiddleware {
	return &ShopeeMiddleware{logger: log, shopeeAuthCollection: shopee}
}

func (m ShopeeMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		path := c.Path()
		shopId := c.Query("shopeeId")

		if strings.HasPrefix(path, "/api/v1/shopee") {

			// m.logger.Sugar().Infof("Raw URL: %s", c.OriginalURL())
			// m.logger.Sugar().Infof("shopId : %s", shopId)
			// m.logger.Sugar().Infof("Path : %s", path)

			// var reqBody ShopeeRequest
			// if err := c.BodyParser(&reqBody); err != nil {
			// 	return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
			// }

			// if err != nil { m.logger.Error("Failed to get shopee access token", zap.Error(err)) }

			// if shopeeAccessToken == "" {
			// 	m.logger.Error("Shopee access token not found")
			// 	return response.ErrorResponse(c, fiber.StatusUnauthorized, "Shopee access token not found", nil)
			// }
			shopeeAccessToken, _ := m.shopeeAuthCollection.GetShopeeAuthByShopId(shopId)
			if shopeeAccessToken != "" { c.Locals("shopeeAccessToken", shopeeAccessToken) }
      m.logger.Info("Middleware:", zap.String("path", (path) ))
		}

		return c.Next()
	}
}

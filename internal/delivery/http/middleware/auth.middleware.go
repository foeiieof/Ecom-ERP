package middleware

import (
	"ecommerce/internal/env"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type IAuthMiddleware interface {
  Handler() fiber.Handler
}

type authMiddleware struct {
  Config *env.Config
  Logger *zap.Logger
}

type AccessAuthEntity struct {
  Sub      string `json:"sub"`
  Type     string `json:"type"`
  Username string `json:"username"`
  jwt.RegisteredClaims
} 

func NewAuthMiddleware(cfg *env.Config, lgs *zap.Logger ) IAuthMiddleware {
  return &authMiddleware{ Config: cfg , Logger: lgs}
}

func (m *authMiddleware) Handler() fiber.Handler {
  return  func (c *fiber.Ctx) error {
    tokenStr := c.Get("Authorization")
    if tokenStr == "" {
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "missing auth",
      })
    }

    parse := strings.Split(tokenStr, " ")
    if len(parse) != 2 || parse[0] != "Bearer" {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "missing auth",
      })
    }

    tokenJwt := parse[1]
    tokenClaims := &AccessAuthEntity{}
    token,err := jwt.ParseWithClaims(tokenJwt, tokenClaims,func(token *jwt.Token) (interface{} ,error) {
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
      }
      return []byte(m.Config.JWT.AuthJWTSecretKey),nil
    })


    // m.Logger.Warn("invalid jwt", zap.Any("", tokenClaims)) 

    if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{ "error": "invalid or expired token", })
    }


    if tokenClaims.ExpiresAt!=nil&& tokenClaims.ExpiresAt.Before(time.Now()){
      return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {"error": "token expired" })
    }

    // if possible
    c.Locals("user_id",tokenClaims.Sub)
    c.Locals("username",tokenClaims.Username)
    return c.Next()
  }
}


// // AuthHandler handles authentication endpoints
// type AuthHandler struct {
// 	authService auth.AuthService
// }

// // NewAuthHandler creates a new authentication handler
// func NewAuthHandler(authService auth.AuthService) *AuthHandler {
// 	return &AuthHandler{
// 		authService: authService,
// 	}
// }

// // Register handles user registration
// func (h *AuthHandler) Register(c *fiber.Ctx) error {
// 	var req auth.RegisterRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	user, token, err := h.authService.Register(c.Context(), &req)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"user":  user,
// 		"token": token,
// 	})
// }

// // Login handles user login
// func (h *AuthHandler) Login(c *fiber.Ctx) error {
// 	var req auth.LoginRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	user, token, err := h.authService.Login(c.Context(), &req)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"user":  user,
// 		"token": token,
// 	})
// }

// // Profile returns the current user's profile
// func (h *AuthHandler) Profile(c *fiber.Ctx) error {
// 	user, ok := GetUserFromContext(c)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not authenticated",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"user": user,
// 	})
// }

// // RefreshToken refreshes the access token
// func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
// 	userID, ok := GetUserIDFromContext(c)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not authenticated",
// 		})
// 	}

// 	token, err := h.authService.RefreshToken(c.Context(), userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"token": token,
// 	})
// }

// // Logout logs out the current user
// func (h *AuthHandler) Logout(c *fiber.Ctx) error {
// 	userID, ok := GetUserIDFromContext(c)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not authenticated",
// 		})
// 	}

// 	if err := h.authService.Logout(c.Context(), userID); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Successfully logged out",
// 	})
// }

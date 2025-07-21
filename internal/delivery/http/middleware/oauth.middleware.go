package middleware

// import (
// 	"strings"

// 	"ecommerce/domain/auth"
// 	jwtService "ecommerce/internal/infrastructure/jwt"

// 	"github.com/gofiber/fiber/v2"
// )

// // OAuthMiddleware handles JWT authentication
// type OAuthMiddleware struct {
// 	jwtService    *jwtService.Service
// 	authService   auth.AuthService
// 	tokenRepo     auth.TokenRepository
// 	skipPaths     map[string]bool
// 	optionalPaths map[string]bool
// }

// // Config for OAuth middleware
// type OAuthConfig struct {
// 	JWTService    *jwtService.Service
// 	AuthService   auth.AuthService
// 	TokenRepo     auth.TokenRepository
// 	SkipPaths     []string // Paths that don't require authentication
// 	OptionalPaths []string // Paths where authentication is optional
// }

// // NewOAuthMiddleware creates a new OAuth middleware with dependency injection
// func NewOAuthMiddleware(config OAuthConfig) *OAuthMiddleware {
// 	skipPaths := make(map[string]bool)
// 	for _, path := range config.SkipPaths {
// 		skipPaths[path] = true
// 	}

// 	optionalPaths := make(map[string]bool)
// 	for _, path := range config.OptionalPaths {
// 		optionalPaths[path] = true
// 	}

// 	return &OAuthMiddleware{
// 		jwtService:    config.JWTService,
// 		authService:   config.AuthService,
// 		tokenRepo:     config.TokenRepo,
// 		skipPaths:     skipPaths,
// 		optionalPaths: optionalPaths,
// 	}
// }

// // Handler returns the middleware handler function
// func (m *OAuthMiddleware) Handler() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		path := c.Path()
// 		
// 		// Skip authentication for certain paths
// 		if m.skipPaths[path] {
// 			return c.Next()
// 		}

// 		// Extract token from Authorization header
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			if m.optionalPaths[path] {
// 				return c.Next()
// 			}
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Authorization header required",
// 			})
// 		}

// 		// Parse Bearer token
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			if m.optionalPaths[path] {
// 				return c.Next()
// 			}
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid authorization header format",
// 			})
// 		}

// 		tokenString := tokenParts[1]

// 		// Check if token is blacklisted
// 		if m.tokenRepo != nil {
// 			blacklisted, err := m.tokenRepo.IsTokenBlacklisted(c.Context(), tokenString)
// 			if err != nil {
// 				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 					"error": "Failed to validate token",
// 				})
// 			}
// 			if blacklisted {
// 				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 					"error": "Token has been revoked",
// 				})
// 			}
// 		}

// 		// Validate JWT token
// 		claims, err := m.jwtService.ValidateToken(tokenString)
// 		if err != nil {
// 			if m.optionalPaths[path] {
// 				return c.Next()
// 			}
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		// Set user information in context
// 		c.Locals("user_id", claims.UserID)
// 		c.Locals("user_email", claims.Email)
// 		c.Locals("token", tokenString)

// 		// Get user from service
// 		user, err := m.authService.ValidateToken(c.Context(), tokenString)
// 		if err != nil {
// 			if m.optionalPaths[path] {
// 				return c.Next()
// 			}
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		c.Locals("user", user)

// 		return c.Next()
// 	}
// }

// // RequireAuth is a stricter middleware that always requires authentication
// func (m *OAuthMiddleware) RequireAuth() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Authorization header required",
// 			})
// 		}

// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid authorization header format",
// 			})
// 		}

// 		tokenString := tokenParts[1]

// 		// Check if token is blacklisted
// 		if m.tokenRepo != nil {
// 			blacklisted, err := m.tokenRepo.IsTokenBlacklisted(c.Context(), tokenString)
// 			if err != nil {
// 				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 					"error": "Failed to validate token",
// 				})
// 			}
// 			if blacklisted {
// 				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 					"error": "Token has been revoked",
// 				})
// 			}
// 		}

// 		// Validate JWT token
// 		claims, err := m.jwtService.ValidateToken(tokenString)
// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		// Set user information in context
// 		c.Locals("user_id", claims.UserID)
// 		c.Locals("user_email", claims.Email)
// 		c.Locals("token", tokenString)

// 		// Get user from service
// 		user, err := m.authService.ValidateToken(c.Context(), tokenString)
// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Invalid token",
// 			})
// 		}

// 		c.Locals("user", user)

// 		return c.Next()
// 	}
// }

// // GetUserFromContext extracts user from Fiber context
// func GetUserFromContext(c *fiber.Ctx) (*auth.User, bool) {
// 	user, ok := c.Locals("user").(*auth.User)
// 	return user, ok
// }

// // GetUserIDFromContext extracts user ID from Fiber context
// func GetUserIDFromContext(c *fiber.Ctx) (string, bool) {
// 	userID, ok := c.Locals("user_id").(string)
// 	return userID, ok
// }

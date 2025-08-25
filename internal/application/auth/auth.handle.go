package auth

import (
	"ecommerce/internal/delivery/http/response"

	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
  CheckAuth(c *fiber.Ctx) error

  PostUserAuthLogin(c *fiber.Ctx) error
}

type AuthHandler struct { }

func NewAuthHandle() IAuthHandler {
  return &AuthHandler{}
} 

func (d *AuthHandler) CheckAuth(c *fiber.Ctx) error{

  return response.SuccessResponse(c,"check-auth","")
}

type IReqUserLogin struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

func (d *AuthHandler)PostUserAuthLogin(c *fiber.Ctx) error {
    var req IReqUserLogin

    if err := c.BodyParser(&req); err != nil {
      return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostUserAuthLogin", err)
    }

    // user, perms, err := service.AuthenticateUser(req.Username, req.Password)
    // if err != nil {
    //     return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
    // }

    // claims := JWTClaims{
    //     UserID:      user.ID,
    //     Roles:       user.RoleNames,
    //     Permissions: perms,
    //     RegisteredClaims: jwt.RegisteredClaims{
    //         ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
    //     },
    // }

    // token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET")))
    // if err != nil {
    //     return err
    // }

    return response.SuccessResponse(c,"login-complete", "")
}

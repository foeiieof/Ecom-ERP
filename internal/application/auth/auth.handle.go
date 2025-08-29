package auth

import (
	"ecommerce/internal/delivery/http/response"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler interface {
  CheckAuth(c *fiber.Ctx) error

  PostUserAuthLogin(c *fiber.Ctx) error
  PostUserAuthRefresh(c *fiber.Ctx) error
}

type authHandler struct {
  service IAuthService  
  logger  *zap.Logger
  validate *validator.Validate

}

func NewAuthHandle(src IAuthService, log *zap.Logger, valid *validator.Validate) AuthHandler {
  return &authHandler{
    service: src,
    logger: log,
    validate: valid,
  }
} 

func (d *authHandler) CheckAuth(c *fiber.Ctx) error{
  return response.SuccessResponse(c,"check-auth","")
}

type IReqUserLogin struct {
  Username string `json:"username" validate:"required"`
  Password string `json:"password" validate:"required"`
}

func (d *authHandler)PostUserAuthLogin(c *fiber.Ctx) error {
  var req IReqUserLogin

  if err := c.BodyParser(&req); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostUserAuthLogin", err)
  }

  if err := d.validate.Struct(req) ; err != nil { return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.PostUserAuthLogin", "Invalidate Body") }


  res,err := d.service.GetJwtFromLogin(c.Context(), req.Username, req.Password)
  if err != nil {return response.ErrorResponse(c,fiber.StatusBadRequest, "handler.PostUserAuthLogin", "username or password invalid")}

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

    // Test to Set Coockie
    c.Cookie(&fiber.Cookie{
    Name: "refresh_token",
    Value: res.RefreshToken,
    Expires: time.Now().Add(time.Hour * 24 * 7 ),
    HTTPOnly: true,
    Secure: true,
    SameSite: "Strict",
    Path: "/api/v1/auth/refresh"}) 

    return response.SuccessResponse(c,"handler.PostUserAuthLogin", res)
}

func (d *authHandler)PostUserAuthRefresh(c *fiber.Ctx) error {

  refreshToken := c.Cookies("refresh_token")

  return response.SuccessResponse(c,"handler.PostUserAuthRefresh", refreshToken)
}

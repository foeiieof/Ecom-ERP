package auth

import (
	"ecommerce/internal/delivery/http/response"
	"ecommerce/internal/env"
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
  Config  *env.Config
  Service IAuthService  
  Logger  *zap.Logger
  Validate *validator.Validate
}

func NewAuthHandle(cfg *env.Config,src IAuthService, log *zap.Logger, valid *validator.Validate) AuthHandler {
  return &authHandler{
    Config: cfg,
    Service: src,
    Logger: log,
    Validate: valid,
  }
} 

func (d *authHandler) CheckAuth(c *fiber.Ctx) error{ return response.SuccessResponse(c,"check-auth","") }

type IReqUserLogin struct {
  Username string `json:"username" validate:"required"`
  Password string `json:"password" validate:"required"`
}

func (d *authHandler)PostUserAuthLogin(c *fiber.Ctx) error {
  var req IReqUserLogin

  if err := c.BodyParser(&req); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body : PostUserAuthLogin", err)
  }

  if err := d.Validate.Struct(req) ; err != nil { return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.PostUserAuthLogin", "Invalidate Body") }


  res,err := d.Service.GetJwtFromLogin(c.Context(), req.Username, req.Password)
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
    Expires: time.Now().Add(time.Duration(d.Config.JWT.AuthJWTRefreshIN) * time.Minute),
    HTTPOnly: true, Secure: true,
    SameSite: fiber.CookieSameSiteStrictMode,
    Path: "/api/v1/auth/refresh"}) 

    return response.SuccessResponse(c,"handler.PostUserAuthLogin", res)
}

func (d *authHandler)PostUserAuthRefresh(c *fiber.Ctx) error {
  refreshToken := c.Cookies("refresh_token")
  if refreshToken ==""{
    return response.ErrorResponse(c,fiber.StatusBadGateway,"handler.PostUserAuthRefresh","refresh token not found!") }

  res,err := d.Service.GetJwtFromRefresh(c.Context(),refreshToken)
  if err != nil {return response.ErrorResponse(c,fiber.StatusBadGateway, "handler.PostUserAuthRefresh", "Authurization not permission")}

  c.Cookie(&fiber.Cookie{
    Name: "refresh_token",
    Value: res.RefreshToken,
    Expires: time.Now().Add(time.Duration(d.Config.JWT.AuthJWTRefreshIN) * time.Minute),
    HTTPOnly: true , Secure: true,
    SameSite: fiber.CookieSameSiteStrictMode,
    Path: "/api/v1/auth/refresh",
  })

  return response.SuccessResponse(c,"handler.PostUserAuthRefresh", res)
}

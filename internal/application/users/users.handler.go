package users

import (
	"ecommerce/internal/delivery/http/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type IUserHandler interface {
  GetUsers(c *fiber.Ctx)    error
  CreateUser(c *fiber.Ctx)  error
  GetUserMe(c *fiber.Ctx)   error
  GetUserByID(c *fiber.Ctx) error
  UpdateUserByID(c *fiber.Ctx) error
  DeleteUserByID(c *fiber.Ctx) error
}

type userHandler struct {
  service IUserService
  logger  *zap.Logger
  validate *validator.Validate
}

func NewUserHandler(
  service IUserService, 
  logger *zap.Logger,
  valid *validator.Validate,
) IUserHandler {
  return &userHandler{
    service: service,
    logger:  logger,
    validate: valid,
  }
}

func (d *userHandler) GetUsers(c *fiber.Ctx) error {

  data,err := d.service.GetUsers(c.Context())

  if err != nil {return response.ErrorResponse(c, fiber.StatusNotFound, "handler.GetAllUser", err.Error() )} 
 
  return response.SuccessResponse(c,"handler.GetAllUser", data)
}

func (d *userHandler) CreateUser(c *fiber.Ctx) error {

  var req IReqCreateUserDTO

  if err := c.BodyParser(&req); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest,"handler.CreateUser", "")
  }

  if err := d.validate.Struct(req); err != nil {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.CreateUser", "")
  }

  resUser, er := d.service.CreateUser(c.Context(),req)
  if er != nil {
     return response.ErrorResponse(c, fiber.StatusInternalServerError, "handler.CreateUser", er.Error())
  }

  return response.SuccessResponse(c, "success-handler.CreateUser", resUser)
}

func (d *userHandler)GetUserMe(c *fiber.Ctx)   error {
  // userID := c.Locals("user_id")
  userMe := c.Locals("username")

  user ,err := d.service.GetUserByUsername(c.Context(), userMe.(string) )
  if err != nil { return response.ErrorResponse(c,fiber.StatusNotFound, "handler.UserMe", err.Error())}
  return response.SuccessResponse(c, "handler.GetUseme" , user)
}


func (d *userHandler)GetUserByID(c *fiber.Ctx) error {

  userID := c.Params("userId")
  if userID == "" { return response.ErrorResponse(c,fiber.StatusBadRequest,"handler.user", "user is required") }

  resUser,err := d.service.GetUserByUsername(c.Context(), userID)
  if err != nil { return response.ErrorResponse(c, fiber.StatusNotFound, "handler.user", "userDetails not found") }

  return response.SuccessResponse(c, "handler.user", resUser)
}

func (d *userHandler)UpdateUserByID(c *fiber.Ctx) error {
  userName := c.Params("userID")
  var bodyParse IReqUpdateUserDTO
  if err := c.BodyParser(&bodyParse) ; err != nil {
    return response.ErrorResponse(c,fiber.StatusBadRequest,"handler.UpdateUserByID", "Invalid body")
  }

  if err := d.validate.Struct(bodyParse) ; err != nil {
    return response.ErrorResponse(c,fiber.StatusBadRequest,"handler.UpdateUserByID", "Invalid body")
  }

  res,err := d.service.UpdateUser(c.Context(),userName,bodyParse)
  if err != nil { return response.ErrorResponse(c, fiber.StatusInternalServerError,"handler.UpdateUser","Failed Update user") }

  return response.SuccessResponse(c,"handler.UpdateUserByID", res)
}

func (d *userHandler)DeleteUserByID(c *fiber.Ctx) error {
  userName := c.Params("userId")
  if userName == "" {
    return response.ErrorResponse(c, fiber.StatusBadRequest, "handler.DeleteUserByID", "user is required")
  }

  res,err := d.service.SoftDeleteUserByUsername(c.Context(),userName)
  if err != nil {
    return response.ErrorResponse(c,fiber.StatusBadRequest, "handler.DeleteUserByUsername", err.Error())
  }

  return response.SuccessResponse(c, "handler.DeleteUserByUsername", res)
}

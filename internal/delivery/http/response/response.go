package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ConvertHeaderTraceID(i any) string {
	val, ok := i.(string)
	if ok {
		return val
	} else {
		return "unknow"
	}
}

func SuccessResponse[T any](c *fiber.Ctx, message string, data T) error {
	timeNow := time.Now()
	reqID := ConvertHeaderTraceID(c.Locals("request_id"))
	return c.Status(fiber.StatusOK).JSON(APIResponse[T]{
		Success:        true,
		RequestID:      reqID,
		Message:        message,
		Data:           data,
		TimestampUnix:  timeNow.Unix(),
		TimestampUTC:   timeNow.UTC().Format(time.RFC3339),
		TimestampLocal: timeNow.Local().Format(time.RFC3339),
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, errMsg string, errDetail any) error {
	timeNow := time.Now()
	reqID := ConvertHeaderTraceID(c.Locals("request_id"))
	return c.Status(statusCode).JSON(APIResponse[any]{
		Success:        false,
		RequestID:      reqID,
		Message:        errMsg,
		Error:          errDetail,
		TimestampUnix:  timeNow.Unix(),
		TimestampUTC:   timeNow.UTC().Format(time.RFC3339),
		TimestampLocal: timeNow.Local().Format(time.RFC3339),
	})
}

package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
  "github.com/google/uuid"
	"go.uber.org/zap"
)

type LogHandler struct {
	logger   *zap.Logger
	localLog fiber.Handler
}

func NewLogHandler(logger *zap.Logger, locallog fiber.Handler) *LogHandler {
	return &LogHandler{logger: logger, localLog: locallog}
}

func (m *LogHandler) ReqLogZap() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		m.logger.Info("HTTP-Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().Header.StatusCode()),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
		)
		return err
	}
}

func (m *LogHandler) ReqLogOriginal() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return m.localLog(c)
	}
}

func (m *LogHandler) TraceLog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
    // m.logger.Info("TraceLog", zap.String("request_id", requestID))
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Locals("request_id", requestID)
		c.Set("X-Request-ID", requestID) 
		return c.Next()
	}
}

// func RequestLoggerMiddleware(logger *zap.Logger) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		start := time.Now()
// 		requestID := c.Locals("request_id")
// 		if requestID == nil {
// 			requestID = "unknown"
// 		}

// 		logger.Info("request_start",
// 			zap.String("request_id", requestID.(string)),
// 			zap.String("method", c.Method()),
// 			zap.String("path", c.Path()),
// 			zap.String("client_ip", c.IP()),
// 		)

// 		err := c.Next()

// 		latency := time.Since(start)
// 		status := c.Response().StatusCode()
// 		level := zap.InfoLevel
// 		var errMsg string
// 		if err != nil || status >= 400 {
// 			level = zap.ErrorLevel
// 			if err != nil {
// 				errMsg = err.Error()
// 			}
// 		}

// 		logger.Check(level, "request_end").Write(
// 			zap.String("request_id", requestID.(string)),
// 			zap.String("method", c.Method()),
// 			zap.String("path", c.Path()),
// 			zap.Int("status", status),
// 			zap.Duration("latency", latency),
// 			zap.String("error", errMsg),
// 		)

// 		return err
// 	}
// }

package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LogHandler struct {
    logger *zap.Logger
    localLog fiber.Handler
}

func NewLogHandler(logger *zap.Logger, locallog fiber.Handler ) *LogHandler {
    return &LogHandler{ logger: logger, localLog: locallog, }
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



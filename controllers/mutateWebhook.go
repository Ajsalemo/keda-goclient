package controllers

import (
	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
)

func MutateWebhook(c *fiber.Ctx) error {
	c.Request().Header.VisitAll(func(key, value []byte) {
		zap.L().Info("req headerKey", zap.String("key", string(key)), zap.String("value", string(value)))
	})
	zap.L().Info("MutateWebhook called")
	return c.JSON(fiber.Map{"test": "test"})
}

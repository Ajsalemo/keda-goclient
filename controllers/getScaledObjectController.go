package controllers

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func GetScaledObject(c *fiber.Ctx) error {
	scaledObjectClient := config.DynammicKubeConfig("scaledObject")
	scaledObjectName := c.Params("scaledObjectName")

	if scaledObjectName == "" {
		zap.L().Error("ScaledObject name is required")
		return c.Status(400).JSON(fiber.Map{"error": "ScaledObject name is required"})
	}

	scaledObject, err := scaledObjectClient.Namespace("apps").Get(context.TODO(), scaledObjectName, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "scaledobjects.keda.sh") && strings.Contains(err.Error(), scaledObjectName) && strings.Contains(err.Error(), "not found") {
			zap.L().Error("ScaledObject not found", zap.Error(err))
			return c.Status(404).JSON(fiber.Map{"error": "ScaledObject not found"})
		}
		zap.L().Error("Error retrieving scaledObjects", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": scaledObject.Object})
}

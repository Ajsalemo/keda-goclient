package controllers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func DeleteScaledObject(c *fiber.Ctx) error {
	scaledObjectClient := config.DynammicKubeConfig("scaledObject")
	scaledObjectName := c.Params("scaledObjectName")

	if scaledObjectName == "" {
		zap.L().Error("ScaledObject name is required")
		return c.Status(400).JSON(fiber.Map{"error": "ScaledObject name is required"})
	}

	err := scaledObjectClient.Namespace("apps").Delete(context.TODO(), scaledObjectName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error("Error deleting scaledobjects", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("ScaledObject deleted", zap.String("scaledObjectName", scaledObjectName))
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Deleted ScaledObject %s", scaledObjectName)})
}

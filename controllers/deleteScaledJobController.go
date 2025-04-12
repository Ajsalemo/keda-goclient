package controllers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func DeleteScaledJob(c *fiber.Ctx) error {
	scaledJobClient := config.DynammicKubeConfig("scaledJob")
	scaledJobName := c.Params("scaledJobName")

	if scaledJobName == "" {
		zap.L().Error("ScaledJob name is required")
		return c.Status(400).JSON(fiber.Map{"error": "ScaledJob name is required"})
	}

	err := scaledJobClient.Namespace("apps").Delete(context.TODO(), scaledJobName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error("Error deleting scaledjobs", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("ScaledJob deleted", zap.String("scaledJobName", scaledJobName))
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Deleted ScaledJob %s", scaledJobName)})
}

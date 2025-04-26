package controllers

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func GetScaledJob(c *fiber.Ctx) error {
	scaledJobClient := config.DynammicKubeConfig("scaledJob")
	scaledJobName := c.Params("scaledJobName")

	if scaledJobName == "" {
		zap.L().Error("ScaledJob name is required")
		return c.Status(400).JSON(fiber.Map{"error": "ScaledJob name is required"})
	}

	scaledJob, err := scaledJobClient.Namespace("apps").Get(context.TODO(), scaledJobName, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "scaledjobs.keda.sh") && strings.Contains(err.Error(), scaledJobName) && strings.Contains(err.Error(), "not found") {
			zap.L().Error("ScaledJob not found", zap.Error(err))
			return c.Status(404).JSON(fiber.Map{"error": "ScaledJob not found"})
		}
		zap.L().Error("Error retrieving scaledjobs", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": scaledJob.Object})
}

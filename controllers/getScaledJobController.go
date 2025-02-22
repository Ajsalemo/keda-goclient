package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func GetScaledJob(c *fiber.Ctx) error {
	clientType := "scaledJob"
	scaledJobName := c.Params("scaledJobName")

	if scaledJobName == "" {
		zap.L().Error("ScaledJob name is required")
		return c.Status(400).JSON(fiber.Map{"error": "ScaledJob name is required"})
	}

	scaledJobClient := config.KubeConfig(clientType)
	scaledJob, err := scaledJobClient.Namespace("default").Get(context.TODO(), scaledJobName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error("Error retrieving scaledjobs", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("ScaledJob created", zap.Any("scaledJob", scaledJob))
	return c.JSON(fiber.Map{"message": scaledJob.Object})

}

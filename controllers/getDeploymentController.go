package controllers

import (
	"context"
	"strings"

	config "github.com/Ajsalemo/keda-goclient/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployment(c *fiber.Ctx) error {
	kubeconfig, err := config.KubeConfig()
	if err != nil {
		zap.L().Error("Error getting kubeconfig", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	deploymentName := c.Params("deploymentName")

	if deploymentName == "" {
		zap.L().Error("ScaledJob name is required")
		return c.Status(400).JSON(fiber.Map{"error": "ScaledJob name is required"})
	}

	deploymentClient := kubeconfig.AppsV1().Deployments("apps")
	deployment, err := deploymentClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "deployments.apps") && strings.Contains(err.Error(), deploymentName) && strings.Contains(err.Error(), "not found") {
			zap.L().Error("Deployment not found", zap.Error(err))
			return c.Status(404).JSON(fiber.Map{"error": "Deployment not found"})
		}

		zap.L().Error("Error retrieving deployment", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": deployment.Spec})
}

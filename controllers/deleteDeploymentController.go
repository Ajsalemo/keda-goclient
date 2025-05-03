package controllers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func DeleteDeployment(c *fiber.Ctx) error {
	deploymentKubeConfig, err := config.KubeConfig()
	if err != nil {
		zap.L().Error("Error getting kubeconfig", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	deploymentName := c.Params("deploymentName")
	deploymentClient := deploymentKubeConfig.AppsV1().Deployments("apps")

	if deploymentName == "" {
		zap.L().Error("Deployment name is required")
		return c.Status(400).JSON(fiber.Map{"error": "Deployment name is required"})
	}

	err = deploymentClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		zap.L().Error("Error deleting deployments", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("Deployment deleted", zap.String("deploymentName", deploymentName))
	return c.JSON(fiber.Map{"message": fmt.Sprintf("Deleted Deployment %s", deploymentName)})
}

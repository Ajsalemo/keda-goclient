package controllers

import (
	"context"

	config "github.com/Ajsalemo/keda-goclient/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(c *fiber.Ctx) error {
	deploymentKubeConfig, err := config.KubeConfig()
	if err != nil {
		zap.L().Error("Error getting kubeconfig", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	deploymentClient := deploymentKubeConfig.AppsV1().Deployments("apps")
	// Since metadata varies depending on the scaler, we use a map of interface to handle the metadata
	// This lets us deserialize the dynamic POST body into scalerMetadataInterfaceMap without having to rely on a typed Struct
	var deploymentStruct config.DeploymentStruct

	if err := c.BodyParser(&deploymentStruct); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Check if containers is provided
	if len(deploymentStruct.Containers) == 0 {
		zap.L().Error("No containers provided")
		return c.Status(400).JSON(fiber.Map{"error": "No containers provided"})
	}

	deploymentResource := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentStruct.Name,
			Namespace: "apps",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentStruct.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentStruct.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentStruct.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: deploymentStruct.Containers,
				},
			},
		},
	}

	deployment, err := deploymentClient.Create(context.TODO(), deploymentResource, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("Error creating deployments", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("Deployment created", zap.Any("deployment", "Created Deployment "+deployment.GetName()))
	return c.JSON(fiber.Map{"message": "Created Deployment " + deployment.GetName()})
}

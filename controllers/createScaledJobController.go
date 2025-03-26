package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func CreateScaledJob(c *fiber.Ctx) error {
	scaledJobClient := config.DynammicKubeConfig("scaledJob")
	// Since metadata varies depending on the scaler, we use a map of interface to handle the metadata
	// This lets us deserialize the dynamic POST body into scalerMetadataInterfaceMap without having to rely on a typed Struct
	var scaledJobStruct config.ScaledJobStruct

	if err := c.BodyParser(&scaledJobStruct); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Check if metadata is provided
	if len(scaledJobStruct.Metadata) == 0 {
		zap.L().Error("No metadata provided")
		return c.Status(400).JSON(fiber.Map{"error": "No metadata provided"})
	}
	// Check if containers is provided
	if len(scaledJobStruct.Containers) == 0 {
		zap.L().Error("No containers provided")
		return c.Status(400).JSON(fiber.Map{"error": "No containers provided"})
	}

	zap.L().Info("scaledJobStruct.Metadata", zap.Any("scaledJobStruct.Metadata", scaledJobStruct.Metadata))
	zap.L().Info("scaledJobStruct.Env", zap.Any("scaledJobStruct.Env", scaledJobStruct.Containers))

	deployment := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "keda.sh/v1alpha1",
			"kind":       "ScaledJob",
			"metadata": map[string]any{
				"name":      "test",
				"namespace": "apps",
			},
			"spec": map[string]any{
				"maxReplicaCount": 1,
				"minReplicaCount": 0,
				"jobTargetRef": map[string]any{
					"template": map[string]any{
						"spec": map[string]any{
							"containers": scaledJobStruct.Containers,
						},
					},
				},
				"triggers": []map[string]any{
					{
						"type":     "github-runner",
						"metadata": scaledJobStruct.Metadata,
						"authenticationRef": map[string]any{
							"name": "github-runner-auth",
						},
					},
				},
			},
		},
	}

	scaledObject, err := scaledJobClient.Namespace("apps").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("Error creating scaledjobs", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("ScaledJob created", zap.Any("scaledObject", scaledObject))
	return c.JSON(fiber.Map{"message": "Created deployment " + scaledObject.GetName()})
}

package controllers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	config "github.com/Ajsalemo/keda-goclient/config"
)

func CreateScaledObject(c *fiber.Ctx) error {
	scaledObjectClient := config.DynammicKubeConfig("scaledObject")
	// Since metadata varies depending on the scaler, we use a map of interface to handle the metadata
	// This lets us deserialize the dynamic POST body into scalerMetadataInterfaceMap without having to rely on a typed Struct
	var scaledObjectStruct config.ScaledObjectStruct

	if err := c.BodyParser(&scaledObjectStruct); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Check if triggers are provided
	if len(scaledObjectStruct.Triggers) == 0 {
		zap.L().Error("No triggers provided")
		return c.Status(400).JSON(fiber.Map{"error": "No triggers provided"})
	}

	zap.L().Info("scaledObjectStruct.Containers", zap.Any("scaledObjectStruct.Containers", scaledObjectStruct.Containers))
	zap.L().Info("scaledObjectStruct.Triggers", zap.Any("scaledObjectStruct.Triggers", scaledObjectStruct.Triggers))

	deployment := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "keda.sh/v1alpha1",
			"kind":       "ScaledObject",
			"metadata": map[string]any{
				"name":      fmt.Sprintf("%s-%s", scaledObjectStruct.Name, "keda-scaledobject"),
				"namespace": "apps",
			},
			"spec": map[string]any{
				"maxReplicaCount": 1,
				"minReplicaCount": 0,
				"scaleTargetRef": map[string]any{
					"name": scaledObjectStruct.Name,
				},
				"triggers": scaledObjectStruct.Triggers,
			},
		},
	}

	scaledObject, err := scaledObjectClient.Namespace("apps").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("Error creating scaledobjects", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("ScaledJob created", zap.Any("scaledObject", scaledObject))
	return c.JSON(fiber.Map{"message": "Created deployment " + scaledObject.GetName()})
}

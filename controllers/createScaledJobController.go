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
	clientType := "scaledJob"
	scaledJobClient := config.KubeConfig(clientType)

	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "keda.sh/v1alpha1",
			"kind":       "ScaledJob",
			"metadata": map[string]interface{}{
				"name":      "test",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"maxReplicaCount": 1,
				"minReplicaCount": 0,
				"jobTargetRef": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{
								{
									"name":            "github-runner",
									"image":           "self-hosted-github-action-runner:latest",
									"imagePullPolicy": "IfNotPresent",
									"env": []map[string]interface{}{
										{
											"name":  "GITHUB_PAT",
											"value": "test-runner",
										},
										{
											"name":  "REPO_OWNER",
											"value": "Ajsalemo",
										},
										{
											"name":  "REPO_NAME",
											"value": "self-hosted-github-action-runner",
										},
										{
											"name":  "REPO_URL",
											"value": "https://github.com/$(REPO_OWNER)/$(REPO_NAME)",
										},
										{
											"name":  "REGISTRATION_TOKEN_API_URL",
											"value": "https://api.github.com/repos/$(REPO_OWNER)/$(REPO_NAME)/actions/runners/registration-token",
										},
									},
								},
							},
						},
					},
				},
				"triggers": []map[string]interface{}{
					{
						"type": "github-runner",
						"metadata": map[string]interface{}{
							"ownerFromEnv":              "REPO_OWNER",
							"runnerScope":               "repo",
							"repoFromEnv":               "REPO_NAME",
							"targetWorkflowQueueLength": "1",
						},
						"authenticationRef": map[string]interface{}{
							"name": "github-runner-auth",
						},
					},
				},
			},
		},
	}

	scaledObject, err := scaledJobClient.Namespace("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		zap.L().Error("Error creating scaledjobs", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("ScaledJob created", zap.Any("scaledObject", scaledObject))
	return c.JSON(fiber.Map{"message": "Created deployment " + scaledObject.GetName()})
}

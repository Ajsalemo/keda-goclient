package controllers

import (
	"context"
	"fmt"
	"time"

	config "github.com/Ajsalemo/keda-goclient/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreateAuthSecret(c *fiber.Ctx) error {
	clientset, err := config.KubeConfig()
	triggerAuthenticationClient := config.DynammicKubeConfig("triggerAutentication")
	if err != nil {
		zap.L().Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var authenticationSecretStruct config.AuthenticationSecret

	if err := c.BodyParser(&authenticationSecretStruct); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// Check if the secret name was provided in the incoming request body
	if authenticationSecretStruct.Name == "" {
		zap.L().Error("No secret name provided")
		return c.Status(400).JSON(fiber.Map{"error": "No secret name provided"})
	}
	// Check if the secret parameter was provided in the incoming request body
	if authenticationSecretStruct.Parameter == "" {
		zap.L().Error("No secret parameter provided")
		return c.Status(400).JSON(fiber.Map{"error": "No secret parameter provided"})
	}
	// Check if the secret value was provided in the incoming request body
	if authenticationSecretStruct.Value == "" {
		zap.L().Error("No secret value provided")
		return c.Status(400).JSON(fiber.Map{"error": "No secret value provided"})
	}

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: authenticationSecretStruct.Name,
		},
		Type: "Opaque",
		Data: map[string][]byte{authenticationSecretStruct.Parameter: []byte(authenticationSecretStruct.Value)},
	}

	_, secretErr := clientset.CoreV1().Secrets("apps").Create(context.TODO(), secret, metav1.CreateOptions{})
	if secretErr != nil {
		zap.L().Error(secretErr.Error())
		return c.Status(500).JSON(fiber.Map{"error": secretErr.Error()})
	}
	// Poll every .5 seconds with a timeout of 60 seconds for secret creation
	// If for some reason the secret creation takes longer than 60 seconds, return a 500
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	start := time.Now()
	for {
		select {
		case <-ticker.C:
			secret, err := clientset.CoreV1().Secrets("apps").Get(context.TODO(), authenticationSecretStruct.Name, metav1.GetOptions{})
			if err != nil {
				zap.L().Error(err.Error())
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			elapsed := time.Since(start)
			zap.L().Info("Polling secret creation: " + elapsed.String())
			// If the elapsed time is greater than 60 seconds, return a 500
			if elapsed > 60*time.Second {
				zap.L().Info("Elapsed time: " + elapsed.String())
				zap.L().Warn("Secret creation took longer than 30 seconds, exiting")
				return c.Status(500).JSON(fiber.Map{"error": "Creation took longer than 30 seconds"})
			}
			zap.L().Info("Polling for secret creation progress..")
			// If the secret has been created, create the KEDA TriggerAuthentication
			if secret.GetName() != "" {
				zap.L().Info("Kind 'Secret' has been created with name " + secret.ObjectMeta.Name)
				zap.L().Info("Creating Kind of 'TriggerAuthentication' using the secret name " + secret.ObjectMeta.Name + " as the secretTargetRef")

				triggerAuthenticationName := fmt.Sprintf("%s-trigger-auth", secret.ObjectMeta.Name)
				triggerAuthenticationDeployment := &unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "keda.sh/v1alpha1",
						"kind":       "TriggerAuthentication",
						"metadata": map[string]any{
							"name":      triggerAuthenticationName,
							"namespace": "apps",
						},
						"spec": map[string]any{
							"secretTargetRef": []map[string]any{
								{
									"parameter": authenticationSecretStruct.Parameter,
									"name":      secret.ObjectMeta.Name,
									"key":       authenticationSecretStruct.Parameter,
								},
							},
						},
					},
				}

				triggerAuthentication, err := triggerAuthenticationClient.Namespace("apps").Create(context.TODO(), triggerAuthenticationDeployment, metav1.CreateOptions{})
				if err != nil {
					zap.L().Error("Error creating Kind of 'TriggerAuthentication' with name '"+triggerAuthenticationName+"' using the secret name '"+secret.ObjectMeta.Name+"' as the secretTargetRef", zap.Error(err))
					return c.Status(500).JSON(fiber.Map{"error": err.Error()})
				}

				zap.L().Info("Succesfully created Kind of 'TriggerAuthentication' using the secret name " + secret.ObjectMeta.Name + " as the secretTargetRef")
				return c.JSON(fiber.Map{"message": "Created TriggerAuthentication " + triggerAuthentication.GetName()})
			}
		}
	}
}

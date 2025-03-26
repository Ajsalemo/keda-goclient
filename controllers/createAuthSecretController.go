package controllers

import (
	"context"

	config "github.com/Ajsalemo/keda-goclient/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateAuthSecret(c *fiber.Ctx) error {
	clientset, err := config.KubeConfig()
	if err != nil {
		zap.L().Error(err.Error())
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-secret",
		},
		Type: "Opaque",
		Data: map[string][]byte{"pat": []byte("123")},
	}

	_, secretErr := clientset.CoreV1().Secrets("apps").Create(context.TODO(), secret, metav1.CreateOptions{})
	if secretErr != nil {
		zap.L().Error(secretErr.Error())
		return c.Status(500).JSON(fiber.Map{"error": secretErr.Error()})
	}

	zap.L().Info("Created secret " + secret.ObjectMeta.Name)

	return c.SendString("CreateAuthSecret")
}

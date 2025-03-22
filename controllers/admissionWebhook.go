package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func AdmissionWebhook(c *fiber.Ctx) error {
	zap.L().Info("AdmissionWebhook called")
	zap.L().Info("Decoding request body into AdmissionReview{}")
	admissionReviewRequest := &admissionv1.AdmissionReview{}
	if err := c.BodyParser(admissionReviewRequest); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	scaledJob := &unstructured.Unstructured{}
	if err := json.Unmarshal(admissionReviewRequest.Request.Object.Raw, scaledJob); err != nil {
		zap.L().Error("Error unmarshalling AdmissionReview request", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	admissionResponse := &admissionv1.AdmissionResponse{}
	admissionResponse.Allowed = true

	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admissionResponse
	admissionReviewResponse.SetGroupVersionKind(admissionReviewRequest.GroupVersionKind())
	admissionReviewResponse.Response.UID = admissionReviewRequest.Request.UID

	zap.L().Info("Logging AdmissionReview response", zap.Any("admissionReviewResponse", admissionReviewResponse))
	return c.JSON(admissionReviewResponse)
}

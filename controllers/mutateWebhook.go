package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func MutateWebhook(c *fiber.Ctx) error {
	c.Request().Header.VisitAll(func(key, value []byte) {
		zap.L().Info("req headerKey", zap.String("key", string(key)), zap.String("value", string(value)))
	})
	zap.L().Info("MutateWebhook called")
	zap.L().Info("Logging Adminssion request body", zap.Any("req body", c.Body()))
	zap.L().Info("Decoding request body into AdmissionReview{}")
	admissionReviewRequest := &admissionv1.AdmissionReview{}
	if err := c.BodyParser(admissionReviewRequest); err != nil {
		zap.L().Error("Error parsing request body", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("Logging AdmissionReview request", zap.Any("admissionReviewRequest", admissionReviewRequest))

	scaledJob := &unstructured.Unstructured{}
	if err := json.Unmarshal(admissionReviewRequest.Request.Object.Raw, scaledJob); err != nil {
		zap.L().Error("Error unmarshalling AdmissionReview request", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("Logging ScaledJob", zap.Any("scaledJob", scaledJob))

	admissionResponse := &admissionv1.AdmissionResponse{}
	var patch string
	patchType := v1.PatchTypeJSONPatch
	// if _, ok := scaledJob.Object["metadata.labels.type"]; !ok {
	// 	zap.L().Info("ScaledJob does not have a type label")
	patch = `[{"op":"add","path":"/metadata/labels","value":{"type":"job"}}]`
	// }

	admissionResponse.Allowed = true
	if patch != "" {
		admissionResponse.PatchType = &patchType
		admissionResponse.Patch = []byte(patch)
	}

	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admissionResponse
	admissionReviewResponse.SetGroupVersionKind(scaledJob.GetObjectKind().GroupVersionKind())
	admissionReviewResponse.Response.UID = admissionReviewRequest.Request.UID

	res, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		zap.L().Error("Error marshalling AdmissionReview response", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	zap.L().Info("Logging AdmissionReview response", zap.Any("res", res))

	return c.JSON(res)
}

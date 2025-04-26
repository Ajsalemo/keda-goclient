package main

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	controllers "github.com/Ajsalemo/keda-goclient/controllers"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	app := fiber.New()
	app.Post("/api/scaledjob/create", controllers.CreateScaledJob)
	app.Post("/api/scaledobject/create", controllers.CreateScaledObject)
	app.Post("/api/deployment/create", controllers.CreateDeployment)
	app.Get("/api/scaledjob/get/:scaledJobName", controllers.GetScaledJob)
	app.Get("/api/scaledobject/get/:scaledObjectName", controllers.GetScaledObject)
	app.Delete("/api/scaledjob/delete/:scaledJobName", controllers.DeleteScaledJob)
	app.Delete("/api/scaledobject/delete/:scaledObjectName", controllers.DeleteScaledObject)
	app.Post("/api/secret/create", controllers.CreateAuthSecret)
	app.Post("/mutate-pods", controllers.MutateWebhook)
	app.Post("/validate-pods", controllers.AdmissionWebhook)
	// Start the servers concurrently through goroutines
	go func() {
		zap.L().Info("Fiber listening on port 8080")
		err := app.Listen(":8080")
		if err != nil {
			zap.L().Fatal(err.Error())
		}
	}()

	zap.L().Info("Go Fiber starting HTTPS server on port 3443")
	err := app.ListenTLS(":3443", "/etc/admission-webhook/tls/tls.crt", "/etc/admission-webhook/tls/tls.key")
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

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
	app.Get("/api/scaledjob/get/:scaledJobName", controllers.GetScaledJob)

	zap.L().Info("Fiber listening on port 8080")
	err := app.Listen(":8080")

	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

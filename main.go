package main

import (
	"gofiber-crud/config"
	"gofiber-crud/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	config.ConnectDB()

	app := fiber.New()

	routes.SetupUserRoutes(app, config.DB)
	routes.SetupAuthRoutes(app, config.DB)

	port := os.Getenv("FIBER_PORT")
	if port != "" {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	} else {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}
}

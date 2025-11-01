package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"gofiber-crud/handlers"
	"gofiber-crud/services"

	"gofiber-crud/middlewares"
)

func SetupUserRoutes(app *fiber.App, db *gorm.DB) {
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	api := app.Group("/api/v1/users")

	api.Use(middlewares.AuthMiddleware)

	api.Post("/", userHandler.CreateUser)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)
	api.Get("/", userHandler.GetAllUsers)
	api.Get("/:id", userHandler.GetUserById)

}

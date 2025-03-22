package routes

import (
	"cs371-backend/internal/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/test-db", controllers.TestDBConnection)
	userController := controllers.NewUserController()
    // Use userController to set up user-related routes
	app.Get("/users", userController.GetAllUsers)
	app.Get("/users/:id", userController.GetUser)
	app.Post("/users", userController.CreateUser)
	app.Put("/users/:id", userController.UpdateUser)
	app.Delete("/users/:id", userController.DeleteUser)

	app.Post("/login", userController.LoginHandler)
}
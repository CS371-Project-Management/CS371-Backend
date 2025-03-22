package routes

import (
	"cs371-backend/internal/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api1 := app.Group("/api1")

	api1.Get("/test-db", controllers.TestDBConnection)

	userController := controllers.NewUserController()
	api1.Get("/users", userController.GetAllUsers)
	api1.Get("/users/:id", userController.GetUser)
	api1.Post("/users", userController.CreateUser)
	api1.Put("/users/:id", userController.UpdateUser)
	api1.Delete("/users/:id", userController.DeleteUser)

	api1.Post("/login", userController.LoginHandler)
	api1.Post("/logout", userController.Logout)

	api1.Post("/request-reset-password", userController.RequestResetPassword)
	api1.Post("/reset-password", userController.ResetPassword)

	courseController := controllers.NewCourseController()
	api1.Post("/course", courseController.CreateCourse)
}

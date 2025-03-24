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
	api1.Get("/courses/:classID", courseController.GetCoursesByClassID)
	api1.Post("/courses", courseController.CreateCourse)
	api1.Put("/courses/:id", courseController.UpdateCourse)

	quizController := controllers.NewQuizController()
	api1.Post("/quizzes", quizController.CreateChoiceQuiz)
	api1.Get("/quizzes/:course_id", quizController.GetAllQuizByCourseID)

	classController := controllers.NewClassController()
	api1.Post("/classes", classController.CreateClassHandler)
	api1.Get("/classes", classController.GetAllClassesHandler)
	api1.Get("/classes/:id", classController.GetClassHandler)
	api1.Put("/classes/:id", classController.UpdateClassHandler)
	api1.Delete("/classes/:id", classController.DeleteClassHandler)	

	api1.Get("/classes/:id/invite_code", classController.GetInviteCodeHandler)
}
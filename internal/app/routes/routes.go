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
	api1.Get("/courses/class/:classID", courseController.GetCoursesByClassID)
	api1.Post("/courses", courseController.CreateCourse)
	api1.Put("/courses/:course_id", courseController.UpdateCourse)
	api1.Delete("/courses/:course_id", courseController.DeleteCourseByID)

	quizController := controllers.NewQuizController()
	api1.Post("/quizzes", quizController.CreateChoiceQuiz)
	api1.Get("/quizzes/:course_id", quizController.GetAllQuizByCourseID)
	api1.Delete("/quizzes/:quiz_id", quizController.DeleteQuizByID)
	api1.Get("/users/:user_id/courses/:course_id/progress", quizController.GetCourseProgress)

	quizHistoryController := controllers.NewQuizHistoryController()
	api1.Post("/quizzes/:quiz_id/submissions", quizHistoryController.TakeQuiz)
	api1.Get("/quiz-histories/:history_id", quizHistoryController.GetQuizHistoryByID)
	api1.Get("/quizzes/:quiz_id/histories", quizHistoryController.GetAllQuizHistoryByQuizID)
}

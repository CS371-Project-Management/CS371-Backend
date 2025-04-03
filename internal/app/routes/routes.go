package routes

import (
	"cs371-backend/internal/app/controllers"
	"cs371-backend/internal/app/middlewares"
	"cs371-backend/internal/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api1 := app.Group("/api1")

	api1.Get("/test-db", controllers.TestDBConnection)

	userController := controllers.NewUserController(services.NewUserService())
	api1.Get("/users", userController.GetAllUsers)
	api1.Get("/users/:id", userController.GetUser)
	api1.Post("/users", userController.CreateUser)
	api1.Put("/users/:id", userController.UpdateUser)
	api1.Delete("/users/:id", userController.DeleteUser)

	api1.Post("/login", userController.LoginHandler)
	api1.Post("/logout", userController.Logout)

	api1.Post("/request-reset-password", userController.RequestResetPassword)
	api1.Post("/reset-password", userController.ResetPassword)

	classController := controllers.NewClassController(services.NewClassService())
	api1.Post("/classes", middlewares.AuthMiddleware, classController.CreateClassHandler)
	api1.Get("/classes", middlewares.AuthMiddleware, classController.GetAllClassesHandler)
	api1.Get("/classes/:id", middlewares.AuthMiddleware, classController.GetClassHandler)
	api1.Put("/classes/:id", middlewares.AuthMiddleware, classController.UpdateClassHandler)
	api1.Delete("/classes/:id", middlewares.AuthMiddleware, classController.DeleteClassHandler)

	api1.Get("/classes/:id/invite_code", middlewares.AuthMiddleware, classController.GetInviteCodeHandler)

	api1.Post("/classes/:id/join-public", middlewares.AuthMiddleware, classController.JoinPublicClassHandler)
	api1.Post("/classes/join-private", middlewares.AuthMiddleware, classController.JoinPrivateClassHandler)

	api1.Post("/classes/leaveClass/:id", middlewares.AuthMiddleware, classController.LeaveClassHandler) //****

	api1.Delete("/classes/:class_id/users/:user_id", middlewares.AuthMiddleware, classController.RemoveUserFromClassHandler)

	api1.Get("/classes/:id/users", classController.GetUsersByClassIDHandler)
	api1.Get("/classes/owned/:user_id", classController.GetOwnedClassesHandler)

	api1.Get("/classes/joined/:user_id", classController.GetClassUserJoinByUserIDHandler)

	courseController := controllers.NewCourseController(services.NewCourseService())
	api1.Get("/courses/class/:classID", middlewares.AuthMiddleware, courseController.GetCoursesByClassID)
	api1.Post("/courses", middlewares.AuthMiddleware, courseController.CreateCourse)
	api1.Put("/courses/:course_id", middlewares.AuthMiddleware, courseController.UpdateCourse)
	api1.Delete("/courses/:course_id", middlewares.AuthMiddleware, courseController.DeleteCourseByID)
	api1.Get("/courses/:course_id", middlewares.AuthMiddleware, courseController.GetCourseByIDHandler)

	quizController := controllers.NewQuizController(services.NewQuizService())
	api1.Post("/quizzes", middlewares.AuthMiddleware, quizController.CreateChoiceQuiz)
	api1.Get("/quizzes/:course_id", middlewares.AuthMiddleware, quizController.GetAllQuizByCourseID)
	api1.Delete("/quizzes/:quiz_id", middlewares.AuthMiddleware, quizController.DeleteQuizByID)
	api1.Get("/users/:user_id/courses/:course_id/progress", middlewares.AuthMiddleware, quizController.GetCourseProgress)

	quizHistoryController := controllers.NewQuizHistoryController(services.NewQuizHistoryService())
	api1.Post("/quizzes/submissions", middlewares.AuthMiddleware, quizHistoryController.TakeQuiz)
	api1.Get("/quiz-histories/:history_id", middlewares.AuthMiddleware, quizHistoryController.GetQuizHistoryByID)
	api1.Get("/quizzes/:quiz_id/histories", middlewares.AuthMiddleware, quizHistoryController.GetAllQuizHistoryByQuizID)
}
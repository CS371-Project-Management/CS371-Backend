package routes

import (
	"cs371-backend/internal/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/test-db", controllers.TestDBConnection)
}

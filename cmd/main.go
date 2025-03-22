package main

import (
	"cs371-backend/config"
	"cs371-backend/db"
	"cs371-backend/internal/app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

func init() {

	//check
	log.Println("DB_HOST:", config.GetEnv("DB_HOST", ""))
	log.Println("DB_PORT:", config.GetEnv("DB_PORT", ""))
	log.Println("DB_USER:", config.GetEnv("DB_USER", ""))
	log.Println("DB_PASSWORD:", config.GetEnv("DB_PASSWORD", ""))
	log.Println("DB_NAME:", config.GetEnv("DB_NAME", ""))
	//check

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// if err := db.RunMigrations(); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Seed
	//seeder := seeders.NewSeeder()
	//seeder.AddSeeder(seeders.SeedUsers)
	//
	//if err := seeder.RunAllSeeders(); err != nil {
	//	log.Fatalf("Failed to run seeders: %v", err)
	//}
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	routes.SetupRoutes(app)

	port := config.GetEnv("APP_PORT", "8080")

	log.Fatal(app.Listen(":" + port))
}

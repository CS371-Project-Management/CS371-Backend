package main

import (
	"cs371-backend/config"
	"cs371-backend/db"
	"cs371-backend/internal/app/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {

	//check
	log.Println("DB_HOST:", config.GetEnv("DB_HOST", ""))
	log.Println("DB_PORT:", config.GetEnv("DB_PORT", ""))
	log.Println("DB_USER:", config.GetEnv("DB_USER", ""))
	log.Println("DB_PASSWORD:", config.GetEnv("DB_PASSWORD", ""))
	log.Println("DB_NAME:", config.GetEnv("DB_NAME", ""))
	//check

	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using default values")
	}

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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET is not set in environment variables!")
	}

	log.Println("✅ JWT_SECRET loaded successfully!")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://0.0.0.0:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	routes.SetupRoutes(app)

	port := config.GetEnv("APP_PORT", "8080")

	log.Fatal(app.Listen(":" + port))
}

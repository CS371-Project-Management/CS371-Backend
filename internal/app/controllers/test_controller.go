package controllers

import (
	"cs371-backend/db"
	"github.com/gofiber/fiber/v2"
)

func TestDBConnection(c *fiber.Ctx) error {
	if err := db.DB.Ping(); err != nil {
		return c.Status(500).SendString("Database connection failed: " + err.Error())
	}
	return c.SendString("Database connected successfully!")
}

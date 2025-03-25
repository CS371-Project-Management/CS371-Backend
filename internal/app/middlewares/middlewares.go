package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func LoggingMiddleware() {
	fmt.Println("Logging middleware")
}

//ก็อปมา ใช้ แก้ token
func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// ใช้ secret จาก .env
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server misconfiguration: missing JWT_SECRET"})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	c.Locals("user_id", claims["user_id"])
	return c.Next()
}
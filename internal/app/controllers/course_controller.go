package controllers

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
	"github.com/gofiber/fiber/v2"
)

type CourseController struct {
	service *services.CourseService
}

func NewCourseController() *CourseController {
	return &CourseController{
		service: services.NewCourseService(),
	}
}

// CreateCourse สร้าง course ใหม่
func (c *CourseController) CreateCourse(ctx *fiber.Ctx) error {
	course := new(models.Course)
	if err := ctx.BodyParser(course); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.service.CreateCourse(course); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(course)
}

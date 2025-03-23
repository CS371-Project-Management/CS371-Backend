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

// GetCoursesByClassID ดึงข้อมูล course ทั้งหมดจาก classID
func (c *CourseController) GetCoursesByClassID(ctx *fiber.Ctx) error {
	classID := ctx.Params("classID")

	courses, err := c.service.GetCoursesByClassID(classID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(courses)
}

// CreateCourse สร้าง course ใหม่
func (c *CourseController) CreateCourse(ctx *fiber.Ctx) error {
	course := new(models.CreateCourseRequest)
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

func (c *CourseController) UpdateCourse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	course := new(models.UpdateCourseRequest)
	if err := ctx.BodyParser(course); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	course.ID = id

	if err := c.service.UpdateCourse(course); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(course)
}

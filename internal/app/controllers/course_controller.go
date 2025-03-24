package controllers

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type CourseController struct {
	service *services.CourseService
}

func NewCourseController() *CourseController {
	return &CourseController{
		service: services.NewCourseService(),
	}
}

func (c *CourseController) GetCoursesByClassID(ctx *fiber.Ctx) error {
	classID := ctx.Params("classID")

	courses, err := c.service.GetCoursesByClassID(classID)
	if err != nil {
		if strings.Contains(err.Error(), "error executing query") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		} else if strings.Contains(err.Error(), "no courses found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(courses)
}

func (c *CourseController) CreateCourse(ctx *fiber.Ctx) error {
	request := new(services.CreateCourseRequest)
	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	course, err := c.service.CreateCourse(request)
	if err != nil {
		if strings.Contains(err.Error(), "error executing query") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(course)
}

func (c *CourseController) UpdateCourse(ctx *fiber.Ctx) error {
	id := ctx.Params("course_id")

	// Parse request body
	course := new(models.UpdateCourseRequest)
	if err := ctx.BodyParser(course); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	course.ID = id

	// Call service to update course
	err := c.service.UpdateCourse(course)
	if err != nil {
		// Handle specific error for not found course
		if strings.Contains(err.Error(), "course not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// Handle other errors (database or server errors)
		if strings.Contains(err.Error(), "error executing query") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return updated course
	return ctx.Status(fiber.StatusOK).JSON(course)
}

func (c *CourseController) DeleteCourseByID(ctx *fiber.Ctx) error {
	courseID := ctx.Params("course_id")

	err := c.service.DeleteCourseByID(courseID)
	if err != nil {
		if strings.Contains(err.Error(), "Course not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Course deleted successfully",
	})
}

package controllers

import (
	"cs371-backend/internal/app/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type QuizController struct {
	quizService *services.QuizService
}

func NewQuizController() *QuizController {
	return &QuizController{
		quizService: services.NewQuizService(),
	}
}

func (c *QuizController) CreateChoiceQuiz(ctx *fiber.Ctx) error {
	request := new(services.CreateQuizRequest)

	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	quiz, err := c.quizService.CreateQuiz(request)
	if err != nil {
		if strings.Contains(err.Error(), "invalid quiz type") == true {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		} else if strings.Contains(err.Error(), "error executing query") == true {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Quiz created successfully",
		"quiz":    quiz,
	})

}

func (c *QuizController) GetAllQuizByCourseID(ctx *fiber.Ctx) error {
	courseID := ctx.Params("course_id")

	quizzes, err := c.quizService.GetAllQuizByCourseID(courseID)
	if err != nil {
		if strings.Contains(err.Error(), "error executing query") == true {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		} else if strings.Contains(err.Error(), "no quizzes found") == true {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"quizzes": quizzes,
	})
}

func (c *QuizController) DeleteQuizByID(ctx *fiber.Ctx) error {
	quizID := ctx.Params("quiz_id")

	err := c.quizService.DeleteQuizByID(quizID)
	if err != nil {
		if strings.Contains(err.Error(), "no quiz found with id") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Quiz not found",
			})
		}

		if strings.Contains(err.Error(), "error deleting quiz") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Quiz deleted successfully",
	})
}

func (c *QuizController) GetCourseProgress(ctx *fiber.Ctx) error {
	courseID := ctx.Params("course_id")
	userID := ctx.Params("user_id")

	result, err := c.quizService.CheckCourseCompletion(courseID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to check course progress",
			"details": err.Error(),
		})
	}

	return ctx.JSON(result)
}

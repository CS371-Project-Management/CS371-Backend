package controllers

import (
	"cs371-backend/internal/app/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type QuizHistoryController struct {
	Service services.QuizHistoryService
}

func NewQuizHistoryController(service services.QuizHistoryService) *QuizHistoryController {
	return &QuizHistoryController{
		Service: service,
	}
}

func (c *QuizHistoryController) TakeQuiz(ctx *fiber.Ctx) error {
	request := new(services.CreateQuizHistoryRequest)

	if err := ctx.BodyParser(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	quizHistory, err := c.Service.TakeQuiz(request)
	if err != nil {
		if strings.Contains(err.Error(), "invalid quiz type") == true {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		} else if strings.Contains(err.Error(), "error executing query") == true {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Quiz submitted successfully",
		"quizHistory": quizHistory,
	})
}

func (c *QuizHistoryController) GetAllQuizHistoryByQuizID(ctx *fiber.Ctx) error {
	quizID := ctx.Params("quiz_id")

	quizHistories, err := c.Service.GetAllQuizHistoryByQuizID(quizID)
	if err != nil {
		if strings.Contains(err.Error(), "error executing query") {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request data",
				"details": err.Error(),
			})
		} else if strings.Contains(err.Error(), "no quiz histories found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "No quiz histories found for this quiz",
				"details": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"quiz_histories": quizHistories,
	})
}

func (c *QuizHistoryController) GetQuizHistoryByID(ctx *fiber.Ctx) error {
	historyID := ctx.Params("history_id")

	history, err := c.Service.GetQuizHistoryByID(historyID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get quiz history",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"quiz_history": history,
	})
}

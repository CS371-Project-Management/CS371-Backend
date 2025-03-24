package controllers

import (
	"cs371-backend/internal/app/services"
	"github.com/gofiber/fiber/v2"
	"log"
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
		log.Println(request, err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	quiz, err := c.quizService.CreateQuiz(request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Quiz created successfully",
		"quiz":    quiz,
	})

}

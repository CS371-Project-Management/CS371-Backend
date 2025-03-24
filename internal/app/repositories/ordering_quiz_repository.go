package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"fmt"
)

type OrderingQuizRepository struct{}

func NewOrderingQuizRepository() *OrderingQuizRepository {
	return &OrderingQuizRepository{}
}

func (r *OrderingQuizRepository) CreateOrderingQuiz(ordering *models.OrderingQuiz) error {
	query := `
		INSERT INTO ordering_quizzes (quiz_id, question)
		VALUES (?, ?)
	`

	_, err := db.DB.Exec(query, ordering.QuizID, ordering.Question)
	if err != nil {
		return fmt.Errorf("CreateOrderingQuiz: error executing query: %w", err)
	}

	return nil
}

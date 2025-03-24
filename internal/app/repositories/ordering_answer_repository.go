package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type OrderingAnswerRepository struct{}

func NewOrderingAnswerRepository() *OrderingAnswerRepository {
	return &OrderingAnswerRepository{}
}

func (r *OrderingAnswerRepository) CreateOrderingAnswer(answer *models.OrderingAnswer) error {
	query := `
    INSERT INTO ordering_answers (id,quiz_id, answer, 'order')
    VALUES (?,?, ?, ?)
  `

	id, err := utils.GenerateUniqueID(db.DB, "ordering_answers", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	answer.ID = id

	_, err = db.DB.Exec(query, answer.ID, answer.QuizID, answer.Answer, answer.Order)
	if err != nil {
		return fmt.Errorf("CreateOrderingAnswer: error executing query: %w", err)
	}

	return nil
}

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

func (r *OrderingQuizRepository) GetOrderingByQuizID(id string) (*models.OrderingQuiz, error) {
	quiz := new(models.OrderingQuiz)

	query := "SELECT quiz_id, question FROM ordering_quizzes WHERE quiz_id = ?"
	err := db.DB.QueryRow(query, id).Scan(&quiz.QuizID, &quiz.Question)
	if err != nil {
		return nil, fmt.Errorf("GetOrderingQuizID: error scanning row: %w", err)
	}

	return quiz, nil
}

func (r *OrderingQuizRepository) DeleteOrderingQuizByID(quizID string) error {
	query := "DELETE FROM ordering_quizzes WHERE quiz_id = ?"
	_, err := db.DB.Exec(query, quizID)
	if err != nil {
		return fmt.Errorf("DeleteOrderingQuizByID: error executing query: %w", err)
	}
	return nil
}

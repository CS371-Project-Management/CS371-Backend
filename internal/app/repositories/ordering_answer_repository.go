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
	query := "INSERT INTO ordering_answers (id ,quiz_id, answer, `order`) VALUES (?,?, ?, ?)"

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

func (r *OrderingAnswerRepository) GetOrderingAnswerByQuizID(id string) ([]models.OrderingAnswer, error) {
	query := "SELECT id, quiz_id, answer, `order` FROM ordering_answers WHERE quiz_id = ?"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("GetOrderingAnswerQuizID: error executing query: %w", err)
	}

	var answers []models.OrderingAnswer
	for rows.Next() {
		var answer models.OrderingAnswer
		err = rows.Scan(&answer.ID, &answer.QuizID, &answer.Answer, &answer.Order)
		if err != nil {
			return nil, fmt.Errorf("GetOrderingAnswerQuizID: error scanning row: %w", err)
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *OrderingAnswerRepository) DeleteOrderingAnswerByID(answerID string) error {
	query := "DELETE FROM ordering_answers WHERE id = ?"
	_, err := db.DB.Exec(query, answerID)
	if err != nil {
		return fmt.Errorf("DeleteOrderingAnswerByID: error executing query: %w", err)
	}
	return nil
}

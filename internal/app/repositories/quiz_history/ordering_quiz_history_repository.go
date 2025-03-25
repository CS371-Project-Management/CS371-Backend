package quiz_history

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz_history"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type OrderingQuizHistoryRepository struct{}

func NewOrderingQuizHistoryRepository() *OrderingQuizHistoryRepository {
	return &OrderingQuizHistoryRepository{}
}

func (r *OrderingQuizHistoryRepository) CreateOrderingHistory(orderingHistory *quiz_history.OrderingHistory) error {
	id, err := utils.GenerateUniqueID(db.DB, "ordering_histories", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	orderingHistory.ID = id

	query := "INSERT INTO ordering_histories (id, quiz_history_id, `order`, answer, result) VALUES (?, ?, ?, ?, ?)"
	_, err = db.DB.Exec(query, orderingHistory.ID, orderingHistory.QuizHistoryID, orderingHistory.Order, orderingHistory.Answer, orderingHistory.Result)
	if err != nil {
		return fmt.Errorf("CreateOrderingHistory: error executing query: %w", err)
	}

	return nil
}

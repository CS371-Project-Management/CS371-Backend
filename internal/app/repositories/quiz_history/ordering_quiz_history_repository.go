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

func (r *OrderingQuizHistoryRepository) GetOrderingHistoryByQuizHistoryID(historyID string) ([]quiz_history.OrderingHistory, error) {
	query := "SELECT id, quiz_history_id, answer, `order`, result FROM ordering_histories WHERE quiz_history_id = ?"
	rows, err := db.DB.Query(query, historyID)
	if err != nil {
		return nil, fmt.Errorf("GetOrderingHistoryByQuizHistoryID: error executing query: %w", err)
	}
	defer rows.Close()

	var answers []quiz_history.OrderingHistory
	for rows.Next() {
		var answer quiz_history.OrderingHistory
		err = rows.Scan(&answer.ID, &answer.QuizHistoryID, &answer.Answer, &answer.Order, &answer.Result)
		if err != nil {
			return nil, fmt.Errorf("GetOrderingHistoryByQuizHistoryID: error scanning row: %w", err)
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

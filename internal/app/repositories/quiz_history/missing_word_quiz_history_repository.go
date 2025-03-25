package quiz_history

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz_history"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type MissingWordQuizHistoryRepository struct{}

func NewMissingWordQuizHistoryRepository() *MissingWordQuizHistoryRepository {
	return &MissingWordQuizHistoryRepository{}
}

func (r *MissingWordQuizHistoryRepository) CreateMissingWordHistory(missingWordHistory *quiz_history.MissingWordHistory) error {
	id, err := utils.GenerateUniqueID(db.DB, "missing_words_quiz_histories", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	missingWordHistory.ID = id

	query := "INSERT INTO missing_words_histories (id, quiz_history_id, answer, result) VALUES (?, ?, ?, ?)"
	_, err = db.DB.Exec(query, missingWordHistory.ID, missingWordHistory.QuizHistoryID, missingWordHistory.Answer, missingWordHistory.Result)
	if err != nil {
		return fmt.Errorf("CreateMissingWordHistory: error executing query: %w", err)
	}

	return nil
}

func (r *MissingWordQuizHistoryRepository) GetMissingWordHistoryByQuizHistoryID(historyID string) (*quiz_history.MissingWordHistory, error) {
	history := new(quiz_history.MissingWordHistory)

	query := "SELECT id, quiz_history_id, answer, result FROM missing_words_histories WHERE quiz_history_id = ?"
	err := db.DB.QueryRow(query, historyID).Scan(&history.ID, &history.QuizHistoryID, &history.Answer, &history.Result)
	if err != nil {
		return nil, fmt.Errorf("GetMissingWordHistoryByQuizHistoryID: error scanning row: %w", err)
	}

	return history, nil
}

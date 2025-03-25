package quiz_history

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz_history"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type ChoiceQuizHistoryRepository struct{}

func NewChoiceQuizHistoryRepository() *ChoiceQuizHistoryRepository {
	return &ChoiceQuizHistoryRepository{}
}

func (r *ChoiceQuizHistoryRepository) CreateChoiceHistory(choiceHistory *quiz_history.ChoiceHistory) error {
	id, err := utils.GenerateUniqueID(db.DB, "choice_quiz_histories", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	choiceHistory.ID = id

	query := "INSERT INTO choice_histories (id, quiz_history_id, answer, result) VALUES (?, ?, ?, ?)"
	_, err = db.DB.Exec(query, choiceHistory.ID, choiceHistory.QuizHistoryID, choiceHistory.Answer, choiceHistory.Result)
	if err != nil {
		return fmt.Errorf("CreateChoiceHistory: error executing query: %w", err)
	}

	return nil
}

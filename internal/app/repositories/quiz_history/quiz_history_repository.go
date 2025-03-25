package quiz_history

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz_history"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type QuizHistoryRepository struct{}

func NewQuizHistoryRepository() *QuizHistoryRepository {
	return &QuizHistoryRepository{}
}

func (r *QuizHistoryRepository) CreateQuizHistory(quizHistory *quiz_history.QuizHistory) error {

	id, err := utils.GenerateUniqueID(db.DB, "quiz_histories", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	quizHistory.ID = id

	query := `
		INSERT INTO quiz_histories (id, user_id, quiz_id, quiz_type, note, passed)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = db.DB.Exec(query, quizHistory.ID, quizHistory.UserID, quizHistory.QuizID, quizHistory.QuizType, quizHistory.Note, quizHistory.Passed)
	if err != nil {
		return fmt.Errorf("CreateQuizHistory: error executing query: %w", err)
	}

	return nil
}

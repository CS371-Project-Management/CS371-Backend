package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"fmt"
)

type MissingWordQuizRepository struct{}

func NewMissingWordQuizRepository() *MissingWordQuizRepository {
	return &MissingWordQuizRepository{}
}

func (r *MissingWordQuizRepository) CreateMissingWordQuiz(missingWord *models.MissingWordQuiz) error {
	query := `
    INSERT INTO missing_word_quizzes (quiz_id, question, answer)
    VALUES (?, ?, ?)
  `

	_, err := db.DB.Exec(query, missingWord.QuizID, missingWord.Question, missingWord.Answer)
	if err != nil {
		return fmt.Errorf("CreateMissingWordQuiz: error executing query: %w", err)
	}

	return nil
}

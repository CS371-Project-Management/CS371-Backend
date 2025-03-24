package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"fmt"
)

type ChoiceQuizRepository struct{}

func NewChoiceQuizRepository() *ChoiceQuizRepository {
	return &ChoiceQuizRepository{}
}

func (r *ChoiceQuizRepository) CreateChoiceQuiz(choice *models.ChoiceQuiz) error {
	query := `
        INSERT INTO choice_quizzes (quiz_id, question, type)
        VALUES (?, ?, ?)
    `

	_, err := db.DB.Exec(query, choice.QuizID, choice.Question, choice.Type)
	if err != nil {
		return fmt.Errorf("CreateChoiceQuiz: error executing query: %w", err)
	}

	return nil
}

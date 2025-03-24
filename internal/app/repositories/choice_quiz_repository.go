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

func (r *ChoiceQuizRepository) GetChoiceByQuizID(id string) (*models.ChoiceQuiz, error) {
	quiz := new(models.ChoiceQuiz)

	query := "SELECT quiz_id, question, type FROM choice_quizzes WHERE quiz_id = ?"
	err := db.DB.QueryRow(query, id).Scan(&quiz.QuizID, &quiz.Question, &quiz.Type)
	if err != nil {
		return nil, fmt.Errorf("GetChoiceQuizID: error scanning row: %w", err)
	}

	return quiz, nil
}

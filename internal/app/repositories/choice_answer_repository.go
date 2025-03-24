package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type ChoiceAnswerRepository struct{}

func NewChoiceAnswerRepository() *ChoiceAnswerRepository {
	return &ChoiceAnswerRepository{}
}

func (r *ChoiceAnswerRepository) CreateChoiceAnswer(answer *models.ChoiceAnswer) error {
	query := "INSERT INTO choice_answers (id, quiz_id, answer, result) VALUES (?, ?, ?, ?)"

	id, err := utils.GenerateUniqueID(db.DB, "quizzes", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	answer.ID = id

	_, err = db.DB.Exec(query, answer.ID, answer.QuizID, answer.Answer, answer.Result)
	if err != nil {
		return fmt.Errorf("CreateChoiceAnswers: error executing query: %w", err)
	}

	return nil
}

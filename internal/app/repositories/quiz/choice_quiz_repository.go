package quiz

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz"
	"fmt"
)

type ChoiceQuizRepository struct{}

func NewChoiceQuizRepository() *ChoiceQuizRepository {
	return &ChoiceQuizRepository{}
}

func (r *ChoiceQuizRepository) CreateChoiceQuiz(choice *quiz.ChoiceQuiz) error {
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

func (r *ChoiceQuizRepository) GetChoiceByQuizID(id string) (*quiz.ChoiceQuiz, error) {
	quiz := new(quiz.ChoiceQuiz)

	query := "SELECT quiz_id, question, type FROM choice_quizzes WHERE quiz_id = ?"
	err := db.DB.QueryRow(query, id).Scan(&quiz.QuizID, &quiz.Question, &quiz.Type)
	if err != nil {
		return nil, fmt.Errorf("GetChoiceQuizID: error scanning row: %w", err)
	}

	return quiz, nil
}

func (r *ChoiceQuizRepository) DeleteChoiceQuizByID(quizID string) error {
	query := "DELETE FROM choice_quizzes WHERE quiz_id = ?"
	_, err := db.DB.Exec(query, quizID)
	if err != nil {
		return fmt.Errorf("DeleteChoiceQuizByID: error executing query: %w", err)
	}
	return nil
}

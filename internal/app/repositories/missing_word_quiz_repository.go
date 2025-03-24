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
    INSERT INTO missing_words_quizzes (quiz_id, question, answer)
    VALUES (?, ?, ?)
  `

	_, err := db.DB.Exec(query, missingWord.QuizID, missingWord.Question, missingWord.Answer)
	if err != nil {
		return fmt.Errorf("CreateMissingWordQuiz: error executing query: %w", err)
	}

	return nil
}

func (r *MissingWordQuizRepository) GetMissingWordByQuizID(id string) (*models.MissingWordQuiz, error) {
	quiz := new(models.MissingWordQuiz)

	query := "SELECT quiz_id, question, answer FROM missing_words_quizzes WHERE quiz_id = ?"
	err := db.DB.QueryRow(query, id).Scan(&quiz.QuizID, &quiz.Question, &quiz.Answer)
	if err != nil {
		return nil, fmt.Errorf("GetMissingWordQuizID: error scanning row: %w", err)
	}

	return quiz, nil
}

func (r *MissingWordQuizRepository) DeleteMissingWordQuizByID(quizID string) error {
	query := "DELETE FROM missing_words_quizzes WHERE quiz_id = ?"
	_, err := db.DB.Exec(query, quizID)
	if err != nil {
		return fmt.Errorf("DeleteMissingWordQuizByID: error executing query: %w", err)
	}
	return nil
}

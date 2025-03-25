package quiz

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type ChoiceAnswerRepository struct{}

func NewChoiceAnswerRepository() *ChoiceAnswerRepository {
	return &ChoiceAnswerRepository{}
}

func (r *ChoiceAnswerRepository) CreateChoiceAnswer(answer *quiz.ChoiceAnswer) error {
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

func (r *ChoiceAnswerRepository) GetChoiceAnswerByQuizID(id string) ([]quiz.ChoiceAnswer, error) {
	query := "SELECT id, quiz_id, answer, result FROM choice_answers WHERE quiz_id = ?"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("GetChoiceAnswerQuizID: error executing query: %w", err)
	}

	var answers []quiz.ChoiceAnswer
	for rows.Next() {
		var answer quiz.ChoiceAnswer
		err = rows.Scan(&answer.ID, &answer.QuizID, &answer.Answer, &answer.Result)
		if err != nil {
			return nil, fmt.Errorf("GetChoiceAnswerQuizID: error scanning row: %w", err)
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *ChoiceAnswerRepository) DeleteChoiceAnswerByID(answerID string) error {
	query := "DELETE FROM choice_answers WHERE id = ?"
	_, err := db.DB.Exec(query, answerID)
	if err != nil {
		return fmt.Errorf("DeleteChoiceAnswerByID: error executing query: %w", err)
	}
	return nil
}

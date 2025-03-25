package quiz_history

import "cs371-backend/internal/app/models/quiz"

type QuizHistory struct {
	ID       string        `json:"id" db:"id"`
	UserID   string        `json:"user_id" db:"user_id"`
	QuizID   string        `json:"quiz_id" db:"quiz_id"`
	QuizType quiz.QuizType `json:"quiz_type" db:"quiz_type"`
	Note     string        `json:"note" db:"note"`
	Passed   bool          `json:"passed" db:"passed"`
}

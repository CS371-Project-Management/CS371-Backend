package models

type QuizHistory struct {
	ID       string   `json:"id" db:"id"`
	UserID   string   `json:"user_id" db:"user_id"`
	QuizID   string   `json:"quiz_id" db:"quiz_id"`
	QuizType QuizType `json:"quiz_type" db:"quiz_type"`
	Note     string   `json:"note" db:"note"`
	Passed   bool     `json:"passed" db:"passed"`
}

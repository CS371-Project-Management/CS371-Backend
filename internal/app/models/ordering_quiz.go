package models

type OrderingQuiz struct {
	QuizID   string `json:"quiz_id" db:"quiz_id"`
	Question string `json:"question" db:"question"`
}

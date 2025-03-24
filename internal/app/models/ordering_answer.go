package models

type OrderingAnswer struct {
	ID     string `json:"id" db:"id"`
	QuizID string `json:"quiz_id" db:"quiz_id"`
	Answer string `json:"answer" db:"answer"`
	Order  string `json:"order" db:"order"`
}

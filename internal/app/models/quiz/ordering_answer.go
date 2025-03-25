package quiz

type OrderingAnswer struct {
	ID     string `json:"id" db:"id"`
	QuizID string `json:"quiz_id" db:"quiz_id"`
	Answer string `json:"answer" db:"answer"`
	Order  int    `json:"order" db:"order"`
}

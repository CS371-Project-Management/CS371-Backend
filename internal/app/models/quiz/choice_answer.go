package quiz

type ChoiceAnswer struct {
	ID     string `json:"id" db:"id"`
	QuizID string `json:"quiz_id" db:"quiz_id"`
	Answer string `json:"answer" db:"answer"`
	Result bool   `json:"result" db:"result"`
}

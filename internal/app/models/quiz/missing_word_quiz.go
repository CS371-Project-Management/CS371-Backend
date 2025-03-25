package quiz

type MissingWordQuiz struct {
	QuizID   string `json:"quiz_id" db:"quiz_id"`
	Question string `json:"question" db:"question"`
	Answer   string `json:"answer" db:"answer"`
}

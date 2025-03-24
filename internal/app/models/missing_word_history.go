package models

type MissingWordHistory struct {
	ID            string `json:"id" db:"id"`
	QuizHistoryID string `json:"quiz_history_id" db:"quiz_history_id"`
	Answer        string `json:"answer" db:"answer"`
	Result        bool   `json:"result" db:"result"`
}

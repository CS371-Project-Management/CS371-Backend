package quiz_history

type OrderingHistory struct {
	ID            string `json:"id" db:"id"`
	QuizHistoryID string `json:"quiz_history_id" db:"quiz_history_id"`
	Order         string `json:"order" db:"order"`
	Answer        string `json:"answer" db:"answer"`
	Result        bool   `json:"result" db:"result"`
}

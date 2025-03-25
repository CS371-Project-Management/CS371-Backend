package services

type CreateQuizHistoryRequest struct {
	ID       string `json:"id,omitempty"`
	UserID   string `json:"user_id"`
	QuizID   string `json:"quiz_id"`
	QuizType string `json:"quiz_type"`
	Note     string `json:"note"`
	Passed   bool   `json:"passed"`
}

type CreateChoiceQuizHistoryData struct {
	Answer string `json:"answer"`
	Result bool   `json:"result"`
}

type CreateOrderQuizHistoryData struct {
	Answer string `json:"answer"`
	Order  int    `json:"order"`
	Result bool   `json:"result"`
}

type CreateMissingWordQuizHistoryData struct {
	Answer string `json:"answer"`
	Result bool   `json:"result"`
}

type QuizHistoryService interface{}

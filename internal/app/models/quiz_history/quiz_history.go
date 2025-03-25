package quiz_history

import (
	"cs371-backend/internal/app/models/quiz"
	"fmt"
)

type Status string

const (
	InProgress Status = "in_progress"
	InCorrect  Status = "in_correct"
	Correct    Status = "correct"
)

type QuizHistory struct {
	ID       string    `json:"id" db:"id"`
	QuizID   string    `json:"quiz_id" db:"quiz_id"`
	UserID   string    `json:"user_id" db:"user_id"`
	QuizType quiz.Type `json:"quiz_type" db:"quiz_type"`
	Note     string    `json:"note" db:"note"`
	Status   Status    `json:"status" db:"status"`
}

func ValidateQuizStatus(status string) (Status, error) {
	switch status {
	case string(InProgress):
		return InProgress, nil
	case string(InCorrect):
		return InCorrect, nil
	case string(Correct):
		return Correct, nil
	default:
		return "", fmt.Errorf("invalid quiz status: %s", status)
	}
}

func QuizStatusToString(status Status) string { return string(status) }

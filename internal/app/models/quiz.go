package models

import "fmt"

type QuizType string

const (
	QuizTypeChoice       QuizType = "choice"
	QuizTypeOrdering     QuizType = "ordering"
	QuizTypeMissingWords QuizType = "missing_words"
)

type Quiz struct {
	ID       string   `json:"id" db:"id"`
	CourseID string   `json:"course_id" db:"course_id"`
	Number   int      `json:"number" db:"number"`
	Point    int      `json:"point" db:"point"`
	QuizType QuizType `json:"quiz_type" db:"quiz_type"`
	Title    string   `json:"title" db:"title"`
	Lesson   string   `json:"lesson" db:"lesson"`
}

func ValidateQuizType(quizType string) (QuizType, error) {
	switch quizType {
	case string(QuizTypeChoice):
		return QuizTypeChoice, nil
	case string(QuizTypeOrdering):
		return QuizTypeOrdering, nil
	case string(QuizTypeMissingWords):
		return QuizTypeMissingWords, nil
	default:
		return "", fmt.Errorf("invalid quiz type: %s", quizType)
	}
}

func QuizTypeToString(quizType QuizType) string {
	return string(quizType)
}

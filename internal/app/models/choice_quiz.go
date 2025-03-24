package models

import "fmt"

type ChoiceType string

const (
	ChoiceSingle   ChoiceType = "single"
	ChoiceMultiple ChoiceType = "multiple"
)

type ChoiceQuiz struct {
	QuizID   string     `json:"quiz_id" db:"quiz_id"`
	Question string     `json:"question" db:"question"`
	Type     ChoiceType `json:"type" db:"type"`
}

func ValidateChoiceType(choiceType string) (ChoiceType, error) {
	switch choiceType {
	case string(ChoiceSingle):
		return ChoiceSingle, nil
	case string(ChoiceMultiple):
		return ChoiceMultiple, nil
	default:
		return "", fmt.Errorf("Invalid choice type: %s", choiceType)
	}
}

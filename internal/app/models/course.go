package models

// Course Create
type Course struct {
	ID              string `json:"id,omitempty" db:"id"`
	ClassID         string `json:"class_id" db:"class_id"`
	Title           string `json:"title" db:"title"`
	Description     string `json:"description" db:"description"`
	DifficultyLevel string `json:"difficulty_level" db:"difficulty_level"`
	Number          int    `json:"number" db:"number"`
}

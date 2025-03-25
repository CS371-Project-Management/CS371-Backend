package models

type user_class struct {
	UserID  uint `json:"user_id" db:"user_id"`
	ClassID uint `json:"class_id" db:"class_id"`
}

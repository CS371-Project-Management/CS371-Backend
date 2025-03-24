package models

type UserClass struct {
	UserID  uint `json:"user_id" db:"user_id"`
	ClassID uint `json:"class_id" db:"class_id"`
}

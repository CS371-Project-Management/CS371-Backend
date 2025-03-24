package models

type Class struct {
	ID            uint   `json:"id" db:"id"`
	InviteCode    string `json:"invite_code" db:"invite_code"`
	Title         string `json:"title" db:"title"`
	Description   string `json:"description" db:"description"`
	Accessibility string `json:"accessibility" db:"accessibility"`
}
package models

type Class struct {
    ID            string `json:"id" db:"id"`
    UserID        string `json:"user_id" db:"user_id"`
    InviteCode    string `json:"invite_code" db:"invite_code"`
    Title         string `json:"title" db:"title"`
    Description   string `json:"description" db:"description"`
    Accessibility bool   `json:"accessibility" db:"accessibility"`
}
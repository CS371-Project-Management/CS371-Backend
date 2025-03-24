package models

type UserCourse struct {
	ID       string `json:"id" db:"id"`
	UserID   string `json:"user_id" db:"user_id"`
	CourseID string `json:"course_id" db:"course_id"`
}

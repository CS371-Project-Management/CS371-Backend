package models

import "fmt"

type UserCourseStatus string

const (
	UserCourseStatusInProgress UserCourseStatus = "in_progress"
	UserCourseStatusCompleted  UserCourseStatus = "completed"
)

type UserCourse struct {
	ID       string           `json:"id" db:"id"`
	UserID   string           `json:"user_id" db:"user_id"`
	Status   UserCourseStatus `json:"status" db:"status"`
	CourseID string           `json:"course_id" db:"course_id"`
}

func ValidateUserCourseStatus(status string) (UserCourseStatus, error) {
	switch status {
	case string(UserCourseStatusInProgress):
		return UserCourseStatusInProgress, nil
	case string(UserCourseStatusCompleted):
		return UserCourseStatusCompleted, nil
	default:
		return "", fmt.Errorf("invalid user course status: %s", status)
	}

}

func UserCourseStatusToString(status UserCourseStatus) string { return string(status) }

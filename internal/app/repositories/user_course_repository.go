package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type UserCourseRepository struct{}

func NewUserCourseRepository() *UserCourseRepository {
	return &UserCourseRepository{}
}

func (r *UserCourseRepository) Create(userCourse *models.UserCourse) error {

	id, err := utils.GenerateUniqueID(db.DB, "user_courses", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}
	userCourse.ID = id

	query := "INSERT INTO user_courses (id, user_id, course_id) VALUES (?, ?, ?)"
	_, err = db.DB.Exec(query, userCourse.ID, userCourse.UserID, userCourse.CourseID)
	if err != nil {
		return fmt.Errorf("Create: error executing query: %w", err)
	}
	return nil
}

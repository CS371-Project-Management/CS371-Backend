package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"fmt"
	"log"
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

func (r *UserCourseRepository) EnrollAllCourseByUserID(userID , classID string) error{
	query := `
		INSERT INTO user_courses (user_id, course_id)
		SELECT ? AS user_id, 
    		c.id AS course_id
		FROM courses c
		WHERE 
    	c.class_id = ?
    	AND NOT EXISTS (
        	SELECT 1 
        	FROM user_courses uc 
        	WHERE uc.user_id = ?  AND uc.course_id = c.id
    );`
	_,err := db.DB.Exec(query, userID, classID, userID)
	log.Println("query run")
	if err != nil {
		return err
	}
	return nil
}
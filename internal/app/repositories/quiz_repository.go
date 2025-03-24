package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type QuizRepository struct{}

func NewQuizRepository() *QuizRepository { return &QuizRepository{} }

func (r *QuizRepository) Create(quiz *models.Quiz) error {

	id, err := utils.GenerateUniqueID(db.DB, "quizzes", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	quiz.ID = id

	query := "INSERT INTO quizzes (id,course_id,number,quiz_type,title,lesson) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.DB.Exec(query, quiz.ID, quiz.CourseID, quiz.Number, quiz.QuizType, quiz.Title, quiz.Lesson)
	if err != nil {
		return fmt.Errorf("Create: error executing query: %w", err)
	}
	return nil
}

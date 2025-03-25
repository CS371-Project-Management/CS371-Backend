package quiz_history

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz_history"
	"cs371-backend/internal/app/utils"
	"database/sql"
	"fmt"
)

type QuizHistoryRepository struct{}

func NewQuizHistoryRepository() *QuizHistoryRepository {
	return &QuizHistoryRepository{}
}

func (r *QuizHistoryRepository) CreateQuizHistory(quizHistory *quiz_history.QuizHistory) error {

	id, err := utils.GenerateUniqueID(db.DB, "quiz_histories", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	quizHistory.ID = id

	query := `
		INSERT INTO quiz_histories (id, user_id, quiz_id, quiz_type, note, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = db.DB.Exec(query, quizHistory.ID, quizHistory.UserID, quizHistory.QuizID, quizHistory.QuizType, quizHistory.Note, quizHistory.Status)
	if err != nil {
		return fmt.Errorf("CreateQuizHistory: error executing query: %w", err)
	}

	return nil
}

func (r *QuizHistoryRepository) GetAllQuizHistoryByQuizID(quizID string) ([]quiz_history.QuizHistory, error) {
	query := "SELECT id, user_id, quiz_id, quiz_type, note, status FROM quiz_histories WHERE quiz_id = ?"
	rows, err := db.DB.Query(query, quizID)
	if err != nil {
		return nil, fmt.Errorf("GetAllQuizHistoryByQuizID: error querying quiz_history: %w", err)
	}
	defer rows.Close()

	var quizHistories []quiz_history.QuizHistory
	for rows.Next() {
		var quizHistory quiz_history.QuizHistory
		err = rows.Scan(
			&quizHistory.ID,
			&quizHistory.UserID,
			&quizHistory.QuizID,
			&quizHistory.QuizType,
			&quizHistory.Note,
			&quizHistory.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllQuizHistoryByQuizID: error scanning row: %w", err)
		}
		quizHistories = append(quizHistories, quizHistory)
	}

	return quizHistories, nil
}
func (r *QuizHistoryRepository) GetQuizHistoryByID(historyID string) (*quiz_history.QuizHistory, error) {
	query := `
        SELECT id, user_id, quiz_id, quiz_type, note, status
        FROM quiz_histories 
        WHERE id = ?
    `

	history := new(quiz_history.QuizHistory)
	err := db.DB.QueryRow(query, historyID).Scan(
		&history.ID,
		&history.UserID,
		&history.QuizID,
		&history.QuizType,
		&history.Note,
		&history.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quiz history not found with id: %s", historyID)
		}
		return nil, fmt.Errorf("error fetching quiz history: %w", err)
	}

	return history, nil
}

func (r *QuizHistoryRepository) DeleteQuizHistoryByQuizID(quizID string, userID string) error {
	query := "DELETE FROM quiz_histories WHERE quiz_id = ? AND user_id = ?"
	_, err := db.DB.Exec(query, quizID, userID)
	if err != nil {
		return fmt.Errorf("error deleting quiz history: %w", err)
	}
	return nil
}

func (r *QuizHistoryRepository) CountAnsweredQuizzes(courseID, userID string) (int, error) {
	var count int
	query := `
        SELECT COUNT(DISTINCT q.id) 
        FROM quizzes q
        JOIN quiz_histories h ON q.id = h.quiz_id
        WHERE q.course_id = ? AND h.user_id = ?
    `
	err := db.DB.QueryRow(query, courseID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting answered quizzes: %w", err)
	}
	return count, nil
}

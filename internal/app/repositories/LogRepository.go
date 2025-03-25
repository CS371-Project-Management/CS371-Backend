package repositories

import (
    "cs371-backend/db"
)

type LogRepository struct{}

func NewLogRepository() *LogRepository {
    return &LogRepository{}
}

// DeleteUserLogs - ลบประวัติทั้งหมดใน quiz_histories, choice_histories, ordering_histories, missing_words_histories
func (r *LogRepository) DeleteUserLogs(userID string) error {
    // quiz_histories
    if err := r.deleteQuizHistories(userID); err != nil {
        return err
    }
    // choice_histories
    if err := r.deleteChoiceHistories(userID); err != nil {
        return err
    }
    // ordering_histories
    if err := r.deleteOrderingHistories(userID); err != nil {
        return err
    }
    // missing_words_histories
    if err := r.deleteMissingWordsHistories(userID); err != nil {
        return err
    }
    return nil
}

func (r *LogRepository) deleteQuizHistories(userID string) error {
    query := `DELETE FROM quiz_histories WHERE user_id = ?`
    _, err := db.DB.Exec(query, userID)
    return err
}

func (r *LogRepository) deleteChoiceHistories(userID string) error {
    query := `DELETE FROM choice_histories WHERE user_id = ?`
    _, err := db.DB.Exec(query, userID)
    return err
}

func (r *LogRepository) deleteOrderingHistories(userID string) error {
    query := `DELETE FROM ordering_histories WHERE user_id = ?`
    _, err := db.DB.Exec(query, userID)
    return err
}

func (r *LogRepository) deleteMissingWordsHistories(userID string) error {
    query := `DELETE FROM missing_words_histories WHERE user_id = ?`
    _, err := db.DB.Exec(query, userID)
    return err
}

package utils

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

func GenerateUniqueID(db *sql.DB, tableName string, fieldName string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("GenerateUniqueID: database connection is nil")
	}

	id := uuid.New().String()
	var exists bool

	for {
		checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", tableName, fieldName)

		err := db.QueryRow(checkQuery, id).Scan(&exists)
		if err != nil {
			return "", fmt.Errorf("GenerateUniqueID: error checking existing id in %s: %w", tableName, err)
		}

		if !exists {
			break
		}

		id = uuid.New().String()
	}

	return id, nil
}

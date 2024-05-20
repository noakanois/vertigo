package database

import (
	"os"
)

func ReadSQLFile(filePath string) (string, error) {
	sql, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(sql), nil
}


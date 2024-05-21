package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func GetDB(databasePath string) (*DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	return &DB{db}, nil
}

func (db *DB) Initialize() error {
	query, sqlErr := ReadSQLFile("data/sql/tables/shoes.sql")
	if sqlErr != nil {
		return fmt.Errorf("can't read file: %v", sqlErr)
	}
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating shoe table: %v", err)
	}
	return nil
}

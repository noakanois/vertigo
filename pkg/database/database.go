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
	query2, sqlErr2 := ReadSQLFile("data/sql/tables/shoentries.sql")
	if sqlErr2 != nil {
		return fmt.Errorf("can't read file: %v", sqlErr2)
	}
	query3, sqlErr3 := ReadSQLFile("data/sql/tables/restaurants.sql")
	if sqlErr3 != nil {
		return fmt.Errorf("can't read file: %v", sqlErr3)
	}
	query4, sqlErr4 := ReadSQLFile("data/sql/tables/foodentries.sql")
	if sqlErr4 != nil {
		return fmt.Errorf("can't read file: %v", sqlErr4)
	}
	query5, sqlErr5 := ReadSQLFile("data/sql/tables/pictures.sql")
	if sqlErr3 != nil {
		return fmt.Errorf("can't read file: %v", sqlErr5)
	}

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating shoe table: %v", err)
	}
	_, err2 := db.Exec(query2)
	if err2 != nil {
		return fmt.Errorf("error creating shoe table: %v", err2)
	}
	_, err3 := db.Exec(query3)
	if err3 != nil {
		return fmt.Errorf("error creating shoe table: %v", err3)
	}
	_, err4 := db.Exec(query4)
	if err4 != nil {
		return fmt.Errorf("error creating shoe table: %v", err3)
	}
	_, err5 := db.Exec(query5)
	if err5 != nil {
		return fmt.Errorf("error creating shoe table: %v", err3)
	}
	return nil
}

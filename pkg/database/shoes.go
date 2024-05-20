package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"vertigo/pkg/dataitems"
)

func (db *DB) InsertShoe(shoe dataitems.Shoe) error {
	query := `INSERT INTO shoes (brand, name, silhouette, image_url, nicknames) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, shoe.Brand, shoe.Name, shoe.Silhouette, shoe.ImageUrl, shoe.Nicknames)
	if err != nil {
		return fmt.Errorf("error inserting new shoe: %v", err)
	}
	return nil
}
func (db *DB) QueryShoesTemplate(query string, params ...interface{}) ([]dataitems.Shoe, error) {
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error querying shoes: %v", err)
	}
	defer rows.Close()

	var shoesList []dataitems.Shoe
	for rows.Next() {
		var shoe dataitems.Shoe
		err := rows.Scan(&shoe.ID, &shoe.Brand, &shoe.Name, &shoe.Silhouette, &shoe.ImageUrl, &shoe.Nicknames)
		if err != nil {
			return nil, fmt.Errorf("error scanning shoe: %v", err)
		}
		shoesList = append(shoesList, shoe)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading shoe rows: %v", err)
	}
	return shoesList, nil
}

func (db *DB) QueryShoeByName(name string) ([]dataitems.Shoe, error) {
	query := `SELECT id, brand, name, silhouette, image_url, nicknames FROM shoes WHERE name = ?`
	shoesList, err := db.QueryShoesTemplate(query, name)
	if err != nil {
		return nil, err
	}
	return shoesList, nil
}

func (db *DB) QueryShoes() ([]dataitems.Shoe, error) {
	query := `SELECT id, brand, name, silhouette, image_url, nicknames FROM shoes`
	shoesList, err := db.QueryShoesTemplate(query)
	if err != nil {
		return nil, err
	}
	return shoesList, nil
}

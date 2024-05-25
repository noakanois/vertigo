package database

import (
	"encoding/json"
	"fmt"
	"vertigo/pkg/stockx"

	_ "github.com/mattn/go-sqlite3"
)

func (db *DB) InsertShoe(pd stockx.ProductDetails) error {
    attributesJSON, err := json.Marshal(pd.Attributes)
    if err != nil {
        return fmt.Errorf("error marshalling attributes to JSON: %v", err)
    }
    query := `INSERT INTO shoes (Name, Subtitle, LastSale, ProductName, MainPicture, Attributes, Description) VALUES (?, ?, ?, ?, ?, ?, ?)`
    _, err = db.Exec(query, pd.Name, pd.Subtitle, pd.LastSale, pd.ProductName, pd.MainPicture, attributesJSON, pd.Description)
    if err != nil {
        return fmt.Errorf("error inserting new product details: %v", err)
    }
    return nil
}

func (db *DB) QueryShoesTemplate(query string, params ...interface{}) ([]stockx.ProductDetails, error) {
    rows, err := db.Query(query, params...)
    if err != nil {
        return nil, fmt.Errorf("error querying product details: %v", err)
    }
    defer rows.Close()

    var productsList []stockx.ProductDetails
    for rows.Next() {
        var pd stockx.ProductDetails
        var attributesJSON string
        err := rows.Scan(&pd.ID, &pd.Name, &pd.Subtitle, &pd.LastSale, &pd.ProductName, &pd.MainPicture, &attributesJSON, &pd.Description)
        if err != nil {
            return nil, fmt.Errorf("error scanning product details: %v", err)
        }
        json.Unmarshal([]byte(attributesJSON), &pd.Attributes)
        productsList = append(productsList, pd)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error reading product details rows: %v", err)
    }
    return productsList, nil
}

func (db *DB) QueryShoeByName(name string) ([]stockx.ProductDetails, error) {
    query := `SELECT ID, Name, Subtitle, LastSale, ProductName, MainPicture, Attributes, Description FROM shoes WHERE Name = ?`
    return db.QueryShoesTemplate(query, name)
}

func (db *DB) QueryShoes() ([]stockx.ProductDetails, error) {
    query := `SELECT ID, Name, Subtitle, LastSale, ProductName, MainPicture, Attributes, Description FROM shoes`
    return db.QueryShoesTemplate(query)
}
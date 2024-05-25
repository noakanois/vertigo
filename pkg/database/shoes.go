package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
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

func (db *DB) InsertPicture(localLocation, discordImageUrl, discordMessageId string, latitude, longitude float64, takenAt time.Time) (int64, error) {
	query := `INSERT INTO pictures (LocalLocation, DiscordImageLink, DiscordMessageId, Latitude, Longitude, TakenAt) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, localLocation, discordImageUrl, discordMessageId, latitude, longitude, takenAt)
	if err != nil {
		return 0, fmt.Errorf("error inserting picture: %v", err)
	}
    log.Printf("succesfully uploaded", discordImageUrl)
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}
	return id, nil
}

func (db *DB) UpdatePictureFilePathAndTimestamp(id int64, newFilePath string) error {
	query := `UPDATE pictures SET LocalLocation = ?, UpdatedAt = ? WHERE id = ?`
	_, err := db.Exec(query, newFilePath, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating picture file path and timestamp: %v", err)
	}
	return nil
}

func (db *DB) InsertShoentry(itemID int64, pictureID int64) (int64, error) {
	query := `INSERT INTO shoentries (ItemID, PictureID) VALUES (?, ?)`
	result, err := db.Exec(query, itemID, pictureID)
	if err != nil {
		return 0, fmt.Errorf("error inserting shoentry: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}
	return id, nil
}


type Shoe struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Subtitle    string    `json:"subtitle"`
	LastSale    string    `json:"last_sale"`
	ProductName string    `json:"product_name"`
	MainPicture string    `json:"main_picture"`
	Attributes  string    `json:"attributes"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type Shoentry struct {
	ID        int64     `json:"id"`
	ItemID    int64     `json:"item_id"`
	PictureID int64     `json:"picture_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}


func (db *DB) GetShoeByProductName(name string) (*Shoe, error) {
	query := `SELECT ID, Name, Subtitle, LastSale, ProductName, MainPicture, Attributes, Description, Timestamp FROM shoes WHERE ProductName = ?`
	row := db.QueryRow(query, name)

	var shoe Shoe
	err := row.Scan(&shoe.ID, &shoe.Name, &shoe.Subtitle, &shoe.LastSale, &shoe.ProductName, &shoe.MainPicture, &shoe.Attributes, &shoe.Description, &shoe.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No shoe found with the given name
		}
		return nil, fmt.Errorf("error retrieving shoe: %v", err)
	}

	return &shoe, nil
}


func (db *DB) GetShoentryByID(id int64) (*Shoentry, error) {
	query := `SELECT ID, ItemID, PictureID, UpdatedAt, CreatedAt FROM shoentries WHERE ID = ?`
	row := db.QueryRow(query, id)

	var shoentry Shoentry
	err := row.Scan(&shoentry.ID, &shoentry.ItemID, &shoentry.PictureID, &shoentry.UpdatedAt, &shoentry.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No shoentry found with the given ID
		}
		return nil, fmt.Errorf("error retrieving shoentry: %v", err)
	}

	return &shoentry, nil
}
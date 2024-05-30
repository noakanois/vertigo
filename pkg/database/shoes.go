package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

type ShoentryDetails struct {
	ShoentryID        int64     `json:"shoentry_id"`
	ItemID            int64     `json:"item_id"`
	ShoeID            int64     `json:"shoe_id"`
	ShoeName          string    `json:"shoe_name"`
	ShoeSubtitle      string    `json:"shoe_subtitle"`
	ShoeLastSale      string    `json:"shoe_last_sale"`
	ShoeProductName   string    `json:"shoe_product_name"`
	ShoeMainPicture   string    `json:"shoe_main_picture"`
	ShoeAttributes    string    `json:"shoe_attributes"`
	ShoeDescription   string    `json:"shoe_description"`
	ShoeTimestamp     time.Time `json:"shoe_timestamp"`
	PictureID         int64     `json:"picture_id"`
	PictureLocalPath  string    `json:"picture_local_path"`
	PictureDiscordURL string    `json:"picture_discord_url"`
	PictureMessageID  string    `json:"picture_message_id"`
	PictureLatitude   float64   `json:"picture_latitude"`
	PictureLongitude  float64   `json:"picture_longitude"`
	PictureTakenAt    time.Time `json:"picture_taken_at"`
	PictureUpdatedAt  time.Time `json:"picture_updated_at"`
	PictureCreatedAt  time.Time `json:"picture_created_at"`
	ShoentryUpdatedAt time.Time `json:"shoentry_updated_at"`
	ShoentryCreatedAt time.Time `json:"shoentry_created_at"`
}

func (db *DB) GetShoeByProductName(name string) (*Shoe, error) {
	query := `SELECT ID, Name, Subtitle, LastSale, ProductName, MainPicture, Attributes, Description, Timestamp FROM shoes WHERE ProductName = ?`
	row := db.QueryRow(query, name)

	var shoe Shoe
	err := row.Scan(&shoe.ID, &shoe.Name, &shoe.Subtitle, &shoe.LastSale, &shoe.ProductName, &shoe.MainPicture, &shoe.Attributes, &shoe.Description, &shoe.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving shoe: %v", err)
	}

	return &shoe, nil
}

func (db *DB) GetShoentryByID(id int64) (*ShoentryDetails, error) {
	query := `
		SELECT 
			shoentries.ID AS ShoentryID,
			shoentries.ItemID,
			shoes.ID AS ShoeID,
			shoes.Name AS ShoeName,
			shoes.Subtitle AS ShoeSubtitle,
			shoes.LastSale AS ShoeLastSale,
			shoes.ProductName AS ShoeProductName,
			shoes.MainPicture AS ShoeMainPicture,
			shoes.Attributes AS ShoeAttributes,
			shoes.Description AS ShoeDescription,
			shoes.Timestamp AS ShoeTimestamp,
			shoentries.PictureID,
			pictures.LocalLocation AS PictureLocalPath,
			pictures.DiscordImageLink AS PictureDiscordURL,
			pictures.DiscordMessageId AS PictureMessageID,
			pictures.Latitude AS PictureLatitude,
			pictures.Longitude AS PictureLongitude,
			pictures.TakenAt AS PictureTakenAt,
			pictures.UpdatedAt AS PictureUpdatedAt,
			pictures.CreatedAt AS PictureCreatedAt,
			shoentries.UpdatedAt AS ShoentryUpdatedAt,
			shoentries.CreatedAt AS ShoentryCreatedAt
		FROM 
			shoentries
		INNER JOIN 
			shoes ON shoentries.ItemID = shoes.ID
		INNER JOIN 
			pictures ON shoentries.PictureID = pictures.ID
		WHERE 
			shoentries.ID = ?
	`

	row := db.QueryRow(query, id)

	var details ShoentryDetails
	err := row.Scan(
		&details.ShoentryID,
		&details.ItemID,
		&details.ShoeID,
		&details.ShoeName,
		&details.ShoeSubtitle,
		&details.ShoeLastSale,
		&details.ShoeProductName,
		&details.ShoeMainPicture,
		&details.ShoeAttributes,
		&details.ShoeDescription,
		&details.ShoeTimestamp,
		&details.PictureID,
		&details.PictureLocalPath,
		&details.PictureDiscordURL,
		&details.PictureMessageID,
		&details.PictureLatitude,
		&details.PictureLongitude,
		&details.PictureTakenAt,
		&details.PictureUpdatedAt,
		&details.PictureCreatedAt,
		&details.ShoentryUpdatedAt,
		&details.ShoentryCreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No shoentry found with the given ID
		}
		return nil, fmt.Errorf("error retrieving shoentry: %v", err)
	}

	return &details, nil
}

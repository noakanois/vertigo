package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"vertigo/pkg/restaurant"

	_ "github.com/mattn/go-sqlite3"
)

func (db *DB) InsertRestaurant(rt restaurant.RestaurantDetails) (int64, error) {
	attributesJSON, err := json.Marshal(rt.Attributes)
	if err != nil {
		return 0, fmt.Errorf("error marshalling attributes to JSON: %v", err)
	}
	query := `INSERT INTO restaurants (Name, Attributes) VALUES (?, ?)`
	result, err := db.Exec(query, rt.Name, attributesJSON)
	if err != nil {
		return 0, fmt.Errorf("error inserting new product details: %v", err)
	}
	return result.LastInsertId()
}

func (db *DB) QueryRestaurantTemplate(query string, params ...interface{}) ([]restaurant.RestaurantDetails, error) {
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error querying product details: %v", err)
	}
	defer rows.Close()

	var productsList []restaurant.RestaurantDetails
	for rows.Next() {
		var rt restaurant.RestaurantDetails
		var attributesJSON string
		err := rows.Scan(&rt.ID, &rt.Name, &attributesJSON)
		if err != nil {
			return nil, fmt.Errorf("error scanning product details: %v", err)
		}
		json.Unmarshal([]byte(attributesJSON), &rt.Attributes)
		productsList = append(productsList, rt)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading product details rows: %v", err)
	}
	return productsList, nil
}

func (db *DB) QueryRestaurantByName(name string) ([]restaurant.RestaurantDetails, error) {
	query := `SELECT ID, Name, Attributes FROM restaurants WHERE Name = ?`
	return db.QueryRestaurantTemplate(query, name)
}

func (db *DB) QueryRestaurants() ([]restaurant.RestaurantDetails, error) {
	query := `SELECT ID, Name, Attributes FROM restaurants`
	return db.QueryRestaurantTemplate(query)
}

func (db *DB) InsertFoodentry(name string, itemID int64, pictureID int64) (int64, error) {
	query := `INSERT INTO foodentries (Name, ItemID, PictureID) VALUES (?, ?, ?)`
	result, err := db.Exec(query, name, itemID, pictureID)
	if err != nil {
		return 0, fmt.Errorf("error inserting foodentry: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}
	return id, nil
}

type Restaurant struct {
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

type Foodentry struct {
	ID        int64     `json:"id"`
	Name      string    `json:"foodname"`
	ItemID    int64     `json:"item_id"`
	PictureID int64     `json:"picture_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type FoodentryDetails struct {
	FoodentryID          int64     `json:"shoentry_id"`
	FoodentryName        string    `json:"foodentry_name"`
	ItemID               int64     `json:"item_id"`
	RestaurantID         int64     `json:"shoe_id"`
	RestaurantName       string    `json:"shoe_name"`
	RestaurantAttributes string    `json:"shoe_attributes"`
	RestaurantTimestamp  time.Time `json:"shoe_timestamp"`
	PictureID            int64     `json:"picture_id"`
	PictureLocalPath     string    `json:"picture_local_path"`
	PictureDiscordURL    string    `json:"picture_discord_url"`
	PictureMessageID     string    `json:"picture_message_id"`
	PictureLatitude      float64   `json:"picture_latitude"`
	PictureLongitude     float64   `json:"picture_longitude"`
	PictureTakenAt       time.Time `json:"picture_taken_at"`
	PictureUpdatedAt     time.Time `json:"picture_updated_at"`
	PictureCreatedAt     time.Time `json:"picture_created_at"`
	FoodentryUpdatedAt   time.Time `json:"shoentry_updated_at"`
	FoodentryCreatedAt   time.Time `json:"shoentry_created_at"`
}

func (db *DB) GetRestaurantByName(name string) (*Restaurant, error) {
	query := `SELECT ID, Name, Attributes, Timestamp FROM restaurants WHERE Name = ?`
	row := db.QueryRow(query, name)

	var restaurant Restaurant
	err := row.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Attributes, &restaurant.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving restaurant: %v", err)
	}

	return &restaurant, nil
}

func (db *DB) GetFoodEntryByID(id int64) (*FoodentryDetails, error) {
	query := `
		SELECT 
			foodentries.ID AS FoodentryID,
			foodentries.ItemID,
			foodentries.Name AS FoodentryName,
			restaurants.ID AS RestaurantID,
			restaurants.Name AS RestaurantName,
			restaurants.Attributes AS RestaurantAttributes,
			restaurants.Timestamp AS RestaurantTimestamp,
			foodentries.PictureID,
			pictures.LocalLocation AS PictureLocalPath,
			pictures.DiscordImageLink AS PictureDiscordURL,
			pictures.DiscordMessageId AS PictureMessageID,
			pictures.Latitude AS PictureLatitude,
			pictures.Longitude AS PictureLongitude,
			pictures.TakenAt AS PictureTakenAt,
			pictures.UpdatedAt AS PictureUpdatedAt,
			pictures.CreatedAt AS PictureCreatedAt,
			foodentries.UpdatedAt AS FoodentryUpdatedAt,
			foodentries.CreatedAt AS FoodentryCreatedAt
		FROM 
			foodentries
		INNER JOIN 
			restaurants ON foodentries.ItemID = restaurants.ID
		INNER JOIN 
			pictures ON foodentries.PictureID = pictures.ID
		WHERE 
			foodentries.ID = ?
	`

	row := db.QueryRow(query, id)

	var details FoodentryDetails
	err := row.Scan(
		&details.FoodentryID,
		&details.ItemID,
		&details.FoodentryName,
		&details.RestaurantID,
		&details.RestaurantName,
		&details.RestaurantAttributes,
		&details.RestaurantTimestamp,
		&details.PictureID,
		&details.PictureLocalPath,
		&details.PictureDiscordURL,
		&details.PictureMessageID,
		&details.PictureLatitude,
		&details.PictureLongitude,
		&details.PictureTakenAt,
		&details.PictureUpdatedAt,
		&details.PictureCreatedAt,
		&details.FoodentryUpdatedAt,
		&details.FoodentryCreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No shoentry found with the given ID
		}
		return nil, fmt.Errorf("error retrieving shoentry: %v", err)
	}

	return &details, nil
}

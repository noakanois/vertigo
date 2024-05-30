package database

import (
	"fmt"
	"time"
)

func (db *DB) InsertPicture(localLocation, discordImageUrl, discordMessageId string, latitude, longitude float64, takenAt time.Time) (int64, error) {
	query := `INSERT INTO pictures (LocalLocation, DiscordImageLink, DiscordMessageId, Latitude, Longitude, TakenAt) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, localLocation, discordImageUrl, discordMessageId, latitude, longitude, takenAt)
	if err != nil {
		return 0, fmt.Errorf("error inserting picture: %v", err)
	}
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

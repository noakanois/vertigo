package main

import (
	"log"
	"net/http"
	"vertigo/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db *database.DB

func initDB() {
	var err error
	db, err = database.GetDB("data/database/test.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

func main() {

	initDB()
	defer db.Close()

	r := gin.Default()

	r.Use(cors.Default())

	r.Static("/img_data", "./img_data")

	r.GET("/shoes", handleShoes)
	r.GET("/shoes/:productName", handleShoeDetails)
	r.GET("/shoentries/:id", handleShoentries)
	r.GET("/recent-shoentries", handleRecentShoentries)
	log.Println("Server is running on port 8080...")
	r.Run(":8080")
}

func handleShoes(c *gin.Context) {
	shoes, err := db.QueryShoes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shoes"})
		return
	}

	c.JSON(http.StatusOK, shoes)
}

func handleShoeDetails(c *gin.Context) {
	name := c.Param("productName")
	shoe, err := db.GetShoeByProductName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shoe details"})
		return
	}
	if shoe == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shoe not found"})
		return
	}

	c.JSON(http.StatusOK, shoe)
}

func handleShoentries(c *gin.Context) {
	shoeID := c.Param("id")
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
            shoes.ID = ?
    `

	rows, err := db.Query(query, shoeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shoentries"})
		return
	}
	defer rows.Close()

	var shoentries []database.ShoentryDetails
	for rows.Next() {
		var details database.ShoentryDetails
		err := rows.Scan(
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan shoentry details"})
			return
		}
		shoentries = append(shoentries, details)
	}

	c.JSON(http.StatusOK, shoentries)
}

func handleRecentShoentries(c *gin.Context) {
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
            pictures ON shoentries.PictureID = pictures.ID
        INNER JOIN 
            shoes ON shoentries.ItemID = shoes.ID
        ORDER BY shoentries.CreatedAt DESC
        LIMIT 10
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying recent shoentries: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent shoentries"})
		return
	}
	defer rows.Close()

	var shoentries []database.ShoentryDetails
	for rows.Next() {
		var shoentry database.ShoentryDetails
		err := rows.Scan(
			&shoentry.ShoentryID,
			&shoentry.ItemID,
			&shoentry.ShoeID,
			&shoentry.ShoeName,
			&shoentry.ShoeSubtitle,
			&shoentry.ShoeLastSale,
			&shoentry.ShoeProductName,
			&shoentry.ShoeMainPicture,
			&shoentry.ShoeAttributes,
			&shoentry.ShoeDescription,
			&shoentry.ShoeTimestamp,
			&shoentry.PictureID,
			&shoentry.PictureLocalPath,
			&shoentry.PictureDiscordURL,
			&shoentry.PictureMessageID,
			&shoentry.PictureLatitude,
			&shoentry.PictureLongitude,
			&shoentry.PictureTakenAt,
			&shoentry.PictureUpdatedAt,
			&shoentry.PictureCreatedAt,
			&shoentry.ShoentryUpdatedAt,
			&shoentry.ShoentryCreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning recent shoentry: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning recent shoentry"})
			return
		}
		shoentries = append(shoentries, shoentry)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error reading recent shoentry rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading recent shoentry rows"})
		return
	}
	c.JSON(http.StatusOK, shoentries)
}

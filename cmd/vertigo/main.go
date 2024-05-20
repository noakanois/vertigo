package main

import (
	"log"
	"vertigo/pkg/database"
	// "vertigo/pkg/shoes"
)

func main() {
	db, err := database.GetDB("data/database/test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// shoe := dataitems.Shoe{
	// 	Name:       "Chicago 3",
	// 	Brand:      "Nike",
	// 	Silhouette: "Air Jordan 1",
	// 	Nicknames:  "Chicago",
	// }

	// err = db.InsertShoe(shoe)
	// if err != nil {
	// 	log.Fatalf("Failed to insert shoe: %v", err)
	// }

	// shoes, err := db.QueryShoeByName("Fragment 1")
	// if err != nil {
	// 	log.Fatalf("Failed to query shoes: %v", err)
	// }

    shoes, err := db.QueryShoes()
	if err != nil {
		log.Fatalf("Failed to query shoes: %v", err)
	}

	for _, shoe := range shoes {
		log.Printf("%+v\n", shoe)
	}
}

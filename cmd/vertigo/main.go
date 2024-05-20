package main

import (
	"log"
    "fmt"
	"vertigo/pkg/database"
    "vertigo/pkg/commandline"
    "flag"
)

func main() {

    listItems := flag.String("list", "", "List all items of type, -list shoes")
    addItems := flag.String("add", "", "Add an item of type, -add shoes")
    flag.Parse()

	db, err := database.GetDB("data/database/test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
    
    if *listItems == "shoes" {
        shoes, err := db.QueryShoes()
        if err != nil {
            log.Fatalf("Failed to query shoes: %v", err)
        }

        for _, shoe := range shoes {
            fmt.Printf("%+v\n", shoe)
        }
    } else if *addItems == "shoes" {
        shoeParams := flag.Args()
        shoe, err := commandline.ParseShoeParams(shoeParams)
        if err != nil {
            log.Fatalf("Failed to parse shoe params: %v", err)
        }

        err = db.InsertShoe(shoe)
        if err != nil {
            log.Fatalf("Failed to insert shoe: %v", err)
        }
        fmt.Println("Shoe added successfully:", shoe)
    } else {
        fmt.Println("Use `-list shoes` to list all shoes, `-add shoe` to add a new shoe, or `-h` for more options")
    }
	
}

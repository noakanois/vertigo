package main

import (
	"flag"
	"fmt"
	"log"
	"vertigo/pkg/database"
	discordBot "vertigo/pkg/discordBot"
	"vertigo/pkg/stockx"
)

func main() {
	listItems := flag.String("list", "", "List all items of type, -list shoes")
	addItems := flag.String("add", "", "Add an item of type, -add https://stockx.com/nike-air-force-1-low-07-chinese-new-year-2024")
	discordNotificationEnabled := flag.Bool("discord", false, "Notify with discord if you are adding a shoe or shoe entry, -discord, default false")
	shoeEntry := flag.String("shoentry", "", "Onboard a shoe image of type, -shoentry filepath -shoe name")
	shoeName := flag.String("shoe", "", "The name of the shoe for the shoentry, -shoe name")
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
	} else if *addItems != "" {
		product, err := stockx.GetShoeInformation(*addItems)
		if err != nil {
			log.Fatalf("Can't get shoe information from stockx: %v", err)
		}
		stockx.GetVisualItem(product.ProductName, product.MainPicture)
		err = db.InsertShoe(product)
		if err != nil {
			log.Fatalf("Failed to insert shoe: %v", err)
		}
		fmt.Println("Shoe added successfully:", product)

		if *discordNotificationEnabled {
			err = discordBot.PostNewShoe(product)
			if err != nil {
				log.Printf("Discord couldn't be notified. %v", err)
			} else {
				fmt.Println("Discord successfully notified.")
			}
		}
	} else if *shoeEntry != "" && *shoeName != "" {
		shoe, err := db.GetShoeByProductName(*shoeName)
		if err != nil {
			log.Fatalf("Failed to get shoe by name: %v", err)
		}
		if shoe == nil {
			log.Fatalf("No shoe found with the name: %s", *shoeName)
		}

		// Onboard the image using the OnboardNewImage function
		pictureID, discordImageUrl, err := discordBot.OnboardNewImage(*shoeEntry)
		if err != nil {
			log.Printf("Failed to onboard new image: %v", err)
		}

		fmt.Println(discordImageUrl)
		fmt.Println(discordImageUrl + "\n\nn\n\n")

		shoentryID, err := db.InsertShoentry(shoe.ID, pictureID)
		if err != nil {
			log.Fatalf("Failed to insert shoentry: %v", err)
		}

		// Insert the new shoentry
		shoentry, err := db.GetShoentryByID(shoentryID)
		if err != nil {
			log.Fatalf("Failed to retrieve shoentry: %v", err)
		}

		fmt.Println("Shoe entry added successfully")

		if *discordNotificationEnabled {
			fmt.Println(discordImageUrl)
			fmt.Println(discordImageUrl)
			err = discordBot.PostNewShoeEntry(shoe.Name, shoentry.ID, discordImageUrl)
			if err != nil {
				log.Printf("Discord couldn't be notified. %v", err)
			} else {
				fmt.Println("Discord successfully notified.")
			}
		}
	} else {
		fmt.Println("Use `-list shoes` to list all shoes, `-add shoe` to add a new shoe, `-shoentry filepath -shoe name` to onboard a shoe image, or `-h` for more options")
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"vertigo/pkg/database"
	discordBot "vertigo/pkg/discordBot"
	rt "vertigo/pkg/restaurant"
	"vertigo/pkg/stockx"
)

func onboardNewRestaurantIfNeeded(db *database.DB, foodname string, foodpath string) (int64, error) {
	var rtDetails rt.RestaurantDetails
	if foodname != "" {
		rtDetails = rt.RestaurantDetails{
			Name: foodname,
		}
	} else {
		rtDetailsList, err := rt.FindRestaurants(foodpath)
		if err != nil {
			return 0, fmt.Errorf("Error requesting restaurant Details from OSM: %v", err)
		}
		if rtDetailsList == nil {
			return 0, fmt.Errorf("No Name provided and no retsaurant found.")
		}
		rtDetails = rtDetailsList[0] // For now, just take the first element in the list of possible restaurants
	}
	restaurant, err := db.GetRestaurantByName(rtDetails.Name)
	if err != nil {
		return 0, fmt.Errorf("Could not check if restaurant already exists: %v", err)
	}

	var id int64
	if restaurant == nil {
		id, err = db.InsertRestaurant(rtDetails)
	} else {
		id = restaurant.ID
	}

	if err != nil {
		return 0, fmt.Errorf("Failed to insert restaurant: %v", err)
	}

	return id, nil

}

func main() {
	listItems := flag.String("list", "", "List all items of type, -list shoes")
	addItems := flag.String("add", "", "Add an item of type, -add https://stockx.com/nike-air-force-1-low-07-chinese-new-year-2024")
	discordNotificationEnabled := flag.Bool("discord", false, "Notify with discord if you are adding a shoe or shoe entry, -discord, default false")

	shoeEntry := flag.String("shoentry", "", "Onboard a shoe image of type, -shoentry filepath -shoe name")
	shoeName := flag.String("shoe", "", "The name of the shoe for the shoentry, -shoe name")

	foodpath := flag.String("foodimage", "", "Upload food image. Will onboard new restaurant if required."+
		"If restaurant name is not found via OSM, you have to provide it yourself with -restaurant name"+
		"-foodImage path -foodName name")
	restaurantName := flag.String("restaurant", "", "Set Restaurant Name -restaurant name")
	foodName := flag.String("foodname", "", "Set Food Name -foodname name")
	flag.Parse()

	db, err := database.GetDB("data/database/test.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if *foodpath != "" && *foodName != "" {

		restaurantid, err := onboardNewRestaurantIfNeeded(db, *restaurantName, *foodpath)
		if err != nil {
			log.Fatalf("Could not add the restaurant: %v", err)
		}

		pictureID, discordImageUrl, err := discordBot.OnboardNewImage(*foodpath, "food")
		if err != nil {
			log.Printf("Failed to onboard new image: %v", err)
		}

		fmt.Println(discordImageUrl)
		fmt.Println(discordImageUrl + "\n\nn\n\n")

		foodentryID, err := db.InsertFoodentry(*foodName, restaurantid, pictureID)
		if err != nil {
			log.Fatalf("Failed to insert foodentry: %v", err)
		}

		foodentryDetails, err := db.GetFoodEntryByID(foodentryID)
		if err != nil {
			log.Fatalf("Failed to retrieve foodentry: %v", err)
		}

		fmt.Println("Shoe entry added successfully")

		if *discordNotificationEnabled {
			err = discordBot.PostNewFoodEntry(*foodentryDetails)
			if err != nil {
				log.Printf("Discord couldn't be notified. %v", err)
			} else {
				fmt.Println("Discord successfully notified.")
			}
		}

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

		pictureID, discordImageUrl, err := discordBot.OnboardNewImage(*shoeEntry, "shoe")
		if err != nil {
			log.Printf("Failed to onboard new image: %v", err)
		}

		fmt.Println(discordImageUrl)
		fmt.Println(discordImageUrl + "\n\nn\n\n")

		shoentryID, err := db.InsertShoentry(shoe.ID, pictureID)
		if err != nil {
			log.Fatalf("Failed to insert shoentry: %v", err)
		}

		shoentryDetails, err := db.GetShoentryByID(shoentryID)
		if err != nil {
			log.Fatalf("Failed to retrieve shoentry: %v", err)
		}

		fmt.Println("Shoe entry added successfully")

		if *discordNotificationEnabled {
			err = discordBot.PostNewShoeEntry(*shoentryDetails)
			if err != nil {
				log.Printf("Discord couldn't be notified. %v", err)
			} else {
				fmt.Println("Discord successfully notified.")
			}
		}
	}
}

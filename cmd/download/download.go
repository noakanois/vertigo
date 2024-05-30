package main

import (
	"log"
	"vertigo/pkg/discordBot"
)

func main() {
	// product, err := stockx.GetShoeInformation("https://stockx.com/nike-air-force-1-low-07-chinese-new-year-2024")
	// if err != nil {
	// 	log.Fatalf("Can't get shoe information from stockx: %v", err)
	// }
	// stockx.GetVisualItem(product.ProductName, product.MainPicture)
	_, _, err := discordBot.OnboardNewImage("img_data/shoentries/test.jpg")
	if err != nil {
		log.Fatalf("Can't onboard image: %v", err)
	}
}

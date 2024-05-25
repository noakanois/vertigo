package main

import (
	"vertigo/pkg/stockx"
)

func main() {
	product := stockx.GetShoeInformation("https://stockx.com/nike-air-force-1-low-07-chinese-new-year-2024")
	stockx.GetVisualItem("1", product.MainPicture)
}

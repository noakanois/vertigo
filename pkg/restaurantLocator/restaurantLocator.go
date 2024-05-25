package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"vertigo/pkg/imageMetadata"
)

type Element struct {
	Type string            `json:"type"`
	Tags map[string]string `json:"tags"`
}

type Response struct {
	Elements []Element `json:"elements"`
}

func findRestaurants(lat, lon float64) {
	query := fmt.Sprintf(
		"[out:json];node[amenity](around:10,%f,%f);out;", lat, lon)
	baseURL := "http://overpass-api.de/api/interpreter"

	v := url.Values{}
	v.Set("data", query)
	res, err := http.Get(baseURL + "?" + v.Encode())
	if err != nil {
		fmt.Println("Error in HTTP request:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, element := range response.Elements {
		fmt.Println("Restaurant Details:")
		for tag, value := range element.Tags {
			fmt.Printf("%s: %s\n", tag, value)
		}
		fmt.Println()
	}
}

func main() {
	metadata, _ := imageMetadata.GetImageMetaData("example.png")
	findRestaurants(metadata.Latitude, metadata.Longitude)
}

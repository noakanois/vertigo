package restaurant

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"vertigo/pkg/imageMetadata"
)

type restaurantJson struct {
	ID   int64          `json:"id"`
	Type string         `json:"type"`
	Lat  float64        `json:"lat"`
	Lon  float64        `json:"lon"`
	Tags RestaurantTags `json:"tags"`
}

type RestaurantTags struct {
	Name           string `json:"name,omitempty"`
	Cuisine        string `json:"cuisine,omitempty"`
	OpeningHours   string `json:"opening_hours,omitempty"`
	Website        string `json:"website,omitempty"`
	Phone          string `json:"phone,omitempty"`
	AddrFull       string `json:"addr:full,omitempty"`
	AddrCity       string `json:"addr:city,omitempty"`
	AddrStreet     string `json:"addr:street,omitempty"`
	AddrPostcode   string `json:"addr:postcode,omitempty"`
	Wheelchair     string `json:"wheelchair,omitempty"`
	Smoking        string `json:"smoking,omitempty"`
	OutdoorSeating string `json:"outdoor_seating,omitempty"`
	Delivery       string `json:"delivery,omitempty"`
}

type RestaurantDetails struct {
	ID         int               `json:"id"`
	Name       string            `json:"type"`
	Attributes map[string]string `json:"tags"`
}

func convertTagsToMap(tags RestaurantTags) map[string]string {
	result := make(map[string]string)
	result["Cuisine"] = tags.Cuisine
	result["OpeningHours"] = tags.OpeningHours
	result["Website"] = tags.Website
	result["Phone"] = tags.Phone
	result["AddrFull"] = tags.AddrFull
	result["AddrCity"] = tags.AddrCity
	result["AddrStreet"] = tags.AddrStreet
	result["AddrPostcode"] = tags.AddrPostcode
	result["Wheelchair"] = tags.Wheelchair
	result["Smoking"] = tags.Smoking
	result["OutdoorSeating"] = tags.OutdoorSeating
	result["Delivery"] = tags.Delivery
	return result
}

func FindRestaurants(foodpath string) ([]RestaurantDetails, error) {
	metadata, err := imageMetadata.GetImageMetaData(foodpath)
	if err != nil {
		return nil, fmt.Errorf("Error in getting Image Metadata: %v", err)
	}
	lat := metadata.Latitude
	lon := metadata.Longitude

	query := fmt.Sprintf(
		"[out:json];node[amenity](around:10,%f,%f);out;", lat, lon)
	baseURL := "http://overpass-api.de/api/interpreter"

	v := url.Values{}
	v.Set("data", query)
	res, err := http.Get(baseURL + "?" + v.Encode())
	if err != nil {
		return nil, fmt.Errorf("Error in HTTP request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	var response struct {
		Elements []restaurantJson `json:"elements"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	if len(response.Elements) == 0 {
		return nil, nil
	}

	restaurants := make([]RestaurantDetails, 0)
	for _, restaurantJson := range response.Elements {
		newRestaurant := RestaurantDetails{
			ID:         0,
			Name:       restaurantJson.Tags.Name,
			Attributes: convertTagsToMap(restaurantJson.Tags),
		}
		restaurants = append(restaurants, newRestaurant)
	}

	return restaurants, nil

}

// func main() {
// 	elements, err := FindRestaurants("example.jpg")
// 	if err != nil {
// 		log.Fatalf("Error finding restaurants: %v", err)
// 	}

// 	fmt.Println("Nearby Restaurants:")
// 	fmt.Println(elements)
// }

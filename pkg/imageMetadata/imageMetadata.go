package imageMetadata

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type ImageMetaData struct {
	Latitude     float64
	Longitude    float64
	CreationDate time.Time
}

func GetImageMetaData(imagePath string) (ImageMetaData, error) {
	latitude, longitude, err := getImageLocation(imagePath)
	if err != nil {
		return ImageMetaData{}, err
	}

	creationDate, err := getImageCreationDate(imagePath)
	if err != nil {
		return ImageMetaData{}, err
	}

	metaData := ImageMetaData{
		Latitude:     latitude,
		Longitude:    longitude,
		CreationDate: creationDate,
	}
	return metaData, nil
}

func getImageCreationDate(imagePath string) (creationDate time.Time, Error error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("error opening image: %v", err)
	}

	ExifMetadata, err := exif.Decode(file)
	if err != nil {
		return time.Time{}, fmt.Errorf("error decoding exif metadata: %v", err)
	}

	createdDate, err := ExifMetadata.DateTime()
	if err != nil {
		return time.Time{}, fmt.Errorf("error extracting date: %v", err)
	}

	return createdDate, nil
}

func getImageLocation(imagePath string) (latitude float64, longitude float64, Error error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("error opening image: %v", err)
	}

	ExifMetadata, err := exif.Decode(file)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("error decoding exif metadata: %v", err)
	}

	lat, long, err := ExifMetadata.LatLong()
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("error extracting location: %v", err)
	}

	return lat, long, nil
}

func main() {
	fl1, fl2, err := getImageLocation("example.png")
	fmt.Printf("error: %v ", err)
	fmt.Printf("Here %v %v:", fl1, fl2)
}

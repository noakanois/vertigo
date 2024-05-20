package imageMetadata

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func GetImageCreationDate(imagePath string) (creationDate time.Time, Error error) {
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

func GetImageLocation(imagePath string) (latitude float64, longitude float64, Error error) {
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

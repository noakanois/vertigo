package imageMetadata

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func GetImageCreationDate(imagePath string) (creationDate time.Time, Error error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return time.Time{}, err
	}

	ExifMetadata, err := exif.Decode(file)
	if err != nil {
		return time.Time{}, err
	}

	createdDate, err := ExifMetadata.DateTime()
	if err != nil {
		return time.Time{}, err
	}

	return createdDate, nil
}

func GetImageLocation(imagePath string) (latitude float64, longitude float64, Error error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0.0, 0.0, err
	}

	ExifMetadata, err := exif.Decode(file)
	if err != nil {
		return 0.0, 0.0, err
	}

	lat, long, err := ExifMetadata.LatLong()
	if err != nil {
		return 0.0, 0.0, err
	}

	return lat, long, nil
}

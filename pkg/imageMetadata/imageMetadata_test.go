package imageMetadata

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetImageCreationDate(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}

	root := filepath.Dir(filepath.Dir(wd))
	fileName := filepath.Join(root, "/testdata/imageMetadata/newYork.jpg")

	date, err := GetImageCreationDate(fileName)
	if err != nil {
		t.Fatalf("Could not read date from image: %v. %v", fileName, err)
	}

	cestLocation := time.FixedZone("CEST", 2*60*60)
	expectedDate := time.Date(2023, 5, 18, 16, 43, 37, 0, cestLocation)

	if !date.Equal(expectedDate) {
		t.Fatalf("Expected date: {%v}, got: {%v}", expectedDate, date)
	}
}

func TestGetImageLocation(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}

	root := filepath.Dir(filepath.Dir(wd))
	fileName := filepath.Join(root, "/testdata/imageMetadata/newYork.jpg")

	latitude, longitude, err := GetImageLocation(fileName)
	if err != nil {
		t.Fatalf("Could not read location from image: %v. %v", fileName, err)
	}

	expectedLatitude := 40.75397777777778
	expectedLongitude := -74.002425

	if latitude != expectedLatitude {
		t.Fatalf("Expected latitude: {%v}, got {%v}", expectedLatitude, latitude)
	}
	if longitude != expectedLongitude {
		t.Fatalf("Expected longitude: {%v}, got {%v}", expectedLongitude, longitude)
	}
}

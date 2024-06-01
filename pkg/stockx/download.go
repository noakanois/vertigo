package stockx

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vertigo/pkg/python"
	"vertigo/pkg/s3"
	"time"
	"github.com/nfnt/resize"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

const (
	imagePath      = "./img_data/shoes"
	numImages      = 36
	imageWidth     = 800
	whiteThreshold = 230
	bucketName     = "vertigo"
)


type DB struct {
	*sql.DB
}


func GetDB(databasePath string) (*DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	return &DB{db}, nil
}
func (db *DB) UpdateShoeImageURLs(productName, mainImgURL, spinningGifURL string) error {
	query := `UPDATE shoes SET MainPicture = ?, SpinningGifURL = ? WHERE ProductName = ?`
	_, err := db.Exec(query, mainImgURL, spinningGifURL, productName)
	if err != nil {
		return fmt.Errorf("error updating shoe image URLs: %v", err)
	}
	return nil
}


func saveImageURLsToDB(itemUUID, mainImgURL, spinningGifURL string) error {
	db, err :=  GetDB("data/database/test.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Assuming there's a function to update the main picture and spinning gif URLs in your database
	err = db.UpdateShoeImageURLs(itemUUID, mainImgURL, spinningGifURL)
	if err != nil {
		return fmt.Errorf("failed to update image URLs in database: %v", err)
	}

	return nil
}

func GetVisualItem(itemUUID, itemImgURL string) error {
	shoeFolderPath := filepath.Join(imagePath, itemUUID)
	firstImgPath := filepath.Join(shoeFolderPath, "main.png")
	if _, err := os.Stat(firstImgPath); os.IsNotExist(err) {
		if err := downloadFirstImg(itemUUID, itemImgURL, true); err != nil {
			return err
		}
		if err := download360Images(itemUUID, itemImgURL, true); err != nil {
			return err
		}
	}
	python.PythonGif(itemUUID, imagePath)
	if err := deleteImages(itemUUID); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	// Upload main.png to R2
	mainImgKey := fmt.Sprintf("%s/main.png", itemUUID)
	mainImgURL, err := s3.UploadToR2(bucketName, mainImgKey, firstImgPath)
	if err != nil {
		return fmt.Errorf("failed to upload main.png to R2: %v", err)
	}

	// Upload spinning.gif to R2
	spinningGifPath := filepath.Join(shoeFolderPath, "spinning.gif")
	spinningGifKey := fmt.Sprintf("%s/spinning.gif", itemUUID)
	spinningGifURL, err := s3.UploadToR2(bucketName, spinningGifKey, spinningGifPath)
	if err != nil {
		return fmt.Errorf("failed to upload spinning.gif to R2: %v", err)
	}

	// Save URLs to the database
	err = saveImageURLsToDB(itemUUID, mainImgURL, spinningGifURL)
	if err != nil {
		return fmt.Errorf("failed to save image URLs to database: %v", err)
	}

	fmt.Println(spinningGifURL, mainImgURL)
	return nil
}
func downloadFirstImg(itemUUID, imgURL string, redownload bool) error {
	shoeFolderPath := prepareShoeFolder(itemUUID)
	firstImgPath := filepath.Join(shoeFolderPath, "main.png")

	if !redownload && fileExists(firstImgPath) {
		return nil
	}

	if download360Image(imgURL, "01", firstImgPath) || downloadStandardImage(imgURL, firstImgPath) {
		if err := trimImage(firstImgPath); err != nil {
			return err
		}
	}
	return nil
}

func download360Image(baseURL, index, savePath string) bool {
	imgURL := convertURLTo360URL(baseURL, index)
	return downloadPicture(imgURL, savePath) == nil
}

func downloadStandardImage(baseURL, savePath string) bool {
	imgURL := fmt.Sprintf("%s?w=%d&bg=FFFFFF", strings.Split(baseURL, "?")[0], imageWidth)
	return downloadPicture(imgURL, savePath) == nil
}

func download360Images(itemUUID, baseURL string, redownload bool) error {
	shoeFolderPath := prepareShoeFolder(itemUUID)
	if !redownload && fileExists(filepath.Join(shoeFolderPath, "spinning.gif")) {
		return nil
	}

	for i := 1; i <= numImages; i++ {
		index := fmt.Sprintf("%02d", i)
		imgSavePath := filepath.Join(shoeFolderPath, index+".jpg")
		if !download360Image(baseURL, index, imgSavePath) {
			return fmt.Errorf("failed to download image %d for %s", i, itemUUID)
		}
	}
	return nil
}


func deleteImages(uuid string) error {
	shoeFolderPath := filepath.Join(imagePath, uuid)
	for i := 1; i <= numImages; i++ {
		index := fmt.Sprintf("%02d", i)
		imgPath := filepath.Join(shoeFolderPath, index+".jpg")
		if err := os.Remove(imgPath); err != nil {
			log.Printf("Failed to remove image %s for %s: %v", index, uuid, err)
		} else {
			log.Printf("Removed image %s for %s", index, uuid)
		}
	}
	return nil
}

func downloadPicture(imgURL, savePath string) error {
	resp, err := http.Get(imgURL)
	if err != nil {
		return fmt.Errorf("failed to download picture from %s: %w", imgURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download picture from %s: status code %d", imgURL, resp.StatusCode)
	}

	out, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to save picture to %s: %w", savePath, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write picture to %s: %w", savePath, err)
	}

	log.Printf("Successfully saved %s to %s", imgURL, savePath)
	return nil
}

func convertURLTo360URL(imageURL, index string) string {
	const imageURLIndex = 4
	parts := strings.Split(imageURL, "/")
	urlKey360 := strings.NewReplacer(
		"-Product.png", "",
		"-Product.jpg", "",
		"-Product_V2.jpg", "",
		"-Product_V2.png", "",
		"_V2", "",
		".png", "",
		".jpg", "",
	).Replace(parts[imageURLIndex])

	return fmt.Sprintf("https://images.stockx.com/360/%s/Images/%s/Lv2/img%s.jpg?w=%d",
		urlKey360, urlKey360, index, imageWidth)
}

func prepareShoeFolder(uuid string) string {
	shoeFolderPath := filepath.Join(imagePath, uuid)
	os.MkdirAll(shoeFolderPath, os.ModePerm)
	return shoeFolderPath
}

func isRowWhite(row []uint8, threshold uint8) bool {
	for _, pixel := range row {
		if pixel < threshold && pixel != 0 {
			return false
		}
	}
	return true
}

func trimImage(path string) error {
	img, err := loadImage(path)
	if err != nil {
		return fmt.Errorf("failed to load image for trimming: %w", err)
	}

	croppedImg := cropImage(img)
	if err := saveImage(path, croppedImg); err != nil {
		return err
	}
	return nil
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func cropImage(img image.Image) image.Image {
	grayImg := resize.Resize(uint(imageWidth), 0, img, resize.Lanczos3)
	bounds := grayImg.Bounds()
	pixels := make([][]uint8, bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]uint8, bounds.Dx())
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := grayImg.At(x, y).RGBA()
			row[x] = uint8((r + g + b) / 3 >> 8)
		}
		pixels[y] = row
	}

	topCrop := 0
	for _, row := range pixels {
		if isRowWhite(row, whiteThreshold) {
			topCrop++
		} else {
			break
		}
	}

	bottomCrop := 0
	for i := len(pixels) - 1; i >= 0; i-- {
		if isRowWhite(pixels[i], whiteThreshold) {
			bottomCrop++
		} else {
			break
		}
	}

	return img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, topCrop, bounds.Dx(), bounds.Dy()-bottomCrop))
}
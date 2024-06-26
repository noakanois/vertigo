package discordBot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vertigo/pkg/database"
	"vertigo/pkg/imageMetadata"
	"vertigo/pkg/stockx"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	session               *discordgo.Session
	botToken              string
	channelIDShoeUpdates  string
	channelIDUploadImages string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error initializing discord bot: %v", err)
	}

	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("Please set the DISCORD_BOT_TOKEN value in .env")
	}

	channelIDShoeUpdates = os.Getenv("DISCORD_NOTIFICATION_CHANNEL")
	if channelIDShoeUpdates == "" {
		log.Fatalf("Please set the DISCORD_NOTIFICATION_CHANNEL value in .env")
	}

	channelIDUploadImages = os.Getenv("DISCORD_IMAGE_CHANNEL")
	if channelIDShoeUpdates == "" {
		log.Fatalf("Please set the DISCORD_IMAGE_CHANNEL value in .env")
	}

	session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func downloadImage(ImageUrl string) (ImageFile *os.File, Error error) {
	resp, err := http.Get(ImageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := os.CreateTemp("", "image-*.jpg")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}

func uploadImage(imageUrl string) (discordImageUrl string, Error error) {
	file, err := downloadImage(imageUrl)
	if err != nil {
		log.Printf("error downloading image: %v", err)
		return "", err
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("error seeking file: %v", err)
		return "", err
	}

	imageFile := &discordgo.File{
		Name:   "image.jpg",
		Reader: file,
	}

	msg, err := session.ChannelFileSend(channelIDUploadImages, imageFile.Name, imageFile.Reader)
	if err != nil {
		fmt.Println("error uploading file,", err)
		return
	}

	if len(msg.Attachments) > 0 {
		fmt.Println("Image uploaded successfully, URL:", msg.Attachments[0].URL)
		return msg.Attachments[0].URL, nil
	} else {
		return "", fmt.Errorf("image uploaded, but no attachments found")
	}
}

func uploadLocalImage(filePath string) (discordImageUrl string, discordMessageId string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return "", "", err
	}
	defer file.Close()

	imageFile := &discordgo.File{
		Name:   file.Name(),
		Reader: file,
	}

	msg, err := session.ChannelFileSend(channelIDUploadImages, imageFile.Name, imageFile.Reader)
	if err != nil {
		log.Printf("error uploading file: %v", err)
		return "", "", err
	}

	if len(msg.Attachments) > 0 {
		fmt.Println("Image uploaded successfully, URL:", msg.Attachments[0].URL)
		return msg.Attachments[0].URL, msg.ID, nil
	} else {
		return "", "", fmt.Errorf("image uploaded, but no attachments found")
	}
}

func OnboardNewImage(filePath string, ImageType string) (int64, string, error) {
	if strings.ToLower(ImageType) != "shoe" && strings.ToLower(ImageType) != "food" {
		return 0, "", fmt.Errorf("cannot open the session: Invalid Image Type passed")
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		return 0, "", fmt.Errorf("cannot open the session: %v", err)
	}
	defer session.Close()

	discordImageUrl, discordMessageId, err := uploadLocalImage(filePath)
	if err != nil {
		return 0, "", fmt.Errorf("error uploading image to Discord: %v", err)
	}

	meta, _ := imageMetadata.GetImageMetaData(filePath)

	db, err := database.GetDB("data/database/test.db")
	if err != nil {
		return 0, "", fmt.Errorf("error connecting to the database: %v", err)
	}

	id, err := db.InsertPicture(filePath, discordImageUrl, discordMessageId, meta.Latitude, meta.Longitude, meta.CreationDate)
	if err != nil {
		return 0, "", fmt.Errorf("error inserting image data into the database: %v", err)
	}

	var newDir string
	switch ImageType {
	case "shoe":
		newDir = "img_data/shoentries"
	case "food":
		newDir = "img_data/food"
	}
	newFilePath := filepath.Join(newDir, fmt.Sprintf("%d.jpg", id))
	err = copyFile(filePath, newFilePath)
	if err != nil {
		return 0, "", fmt.Errorf("error copying image file: %v", err)
	}

	err = db.UpdatePictureFilePathAndTimestamp(id, newFilePath)
	if err != nil {
		return 0, "", fmt.Errorf("error updating image file path and timestamp in the database: %v", err)
	}

	return id, discordImageUrl, nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func PostNewFoodEntry(food database.FoodentryDetails) error {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		return fmt.Errorf("cannot open the session: %v", err)
	}

	if err != nil {
		return fmt.Errorf("cannot upload the image to Discord: %v", err)
	}

	var tags map[string]string
	json.Unmarshal([]byte(food.RestaurantAttributes), &tags)

	attributeString := ""

	if tags["Cuisine"] != "" {
		attributeString += fmt.Sprintf("%s\n", tags["Cuisine"])
	}

	if tags["AddrFull"] != "" {
		attributeString += fmt.Sprintf("Cuisine: %s\n", tags["Cuisine"])
	} else if tags["AddrCity"] != "" && tags["AddrPostcode"] != "" && tags["AddrStreet"] != "" {
		attributeString += fmt.Sprintf("%s, %s, %s\n", tags["AddrStreet"], tags["AddrCity"], tags["AddrPostcode"])
	}

	if tags["Website"] != "" {
		attributeString += fmt.Sprintf("%s\n", tags["Website"])
	}

	description := fmt.Sprintf(
		"A new food entry has been added for **%s**!\n\n%s\n%s\nTaken at %v",
		food.RestaurantName, food.FoodentryName, attributeString, food.PictureTakenAt,
	)
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("New Food Entry: %s", food.FoodentryName),
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: food.PictureDiscordURL,
		},
		Color: 0x4c00b0,
	}

	_, err = session.ChannelMessageSendEmbed(channelIDShoeUpdates, embed)
	if err != nil {
		return fmt.Errorf("cannot send the embedded message: %v", err)
	}

	err = session.Close()
	if err != nil {
		return fmt.Errorf("cannot close the session: %v", err)
	}

	return nil
}

func PostNewShoe(shoe stockx.ProductDetails) error {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		return fmt.Errorf("cannot open the session: %v", err)
	}
	path := "img_data/shoes/" + shoe.ProductName + "/gif/" + shoe.ProductName + ".gif"
	discordImageUrl, _, err := uploadLocalImage(path)
	if err != nil {
		return fmt.Errorf("cannot upload the image to Discord: %v", err)
	}

	var attributes []string
	for key, value := range shoe.Attributes {
		attributes = append(attributes, fmt.Sprintf("%s: %s", key, value))
	}

	attributesString := strings.Join(attributes, "\n")

	description := fmt.Sprintf(
		"%s\n%s\n\nLast Sale: %s\nAttributes: %v\n\nDescription:\n%s",
		shoe.Name,
		shoe.Subtitle,
		shoe.LastSale,
		attributesString,
		shoe.Description,
	)
	embed := &discordgo.MessageEmbed{
		Title:       shoe.Name,
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: discordImageUrl,
		},
		Color: 0x4c00b0,
	}

	_, err = session.ChannelMessageSendEmbed(channelIDShoeUpdates, embed)
	if err != nil {
		return fmt.Errorf("cannot send the embedded message: %v", err)
	}

	err = session.Close()
	if err != nil {
		return fmt.Errorf("cannot close the session: %v", err)
	}

	return nil
}

func PostNewShoeEntry(shoentry database.ShoentryDetails) error {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		return fmt.Errorf("cannot open the session: %v", err)
	}
	defer session.Close()

	description := fmt.Sprintf(
		"A new shoe entry has been added for **%s %s**!\nTaken at %s",
		shoentry.ShoeName,
		shoentry.ShoeSubtitle,
		shoentry.PictureTakenAt,
	)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("New Shoe Entry: %s", shoentry.ShoeName),
		Description: description,
		Image: &discordgo.MessageEmbedImage{
			URL: shoentry.PictureDiscordURL,
		},
		Color: 0x4c00b0,
	}

	_, err = session.ChannelMessageSendEmbed(channelIDShoeUpdates, embed)
	if err != nil {
		return fmt.Errorf("cannot send the embedded message: %v", err)
	}

	return nil
}

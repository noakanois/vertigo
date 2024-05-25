package discordBot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func uploadLocalImage(filePath string) (discordImageUrl string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return "", err
	}
	defer file.Close()

	imageFile := &discordgo.File{
		Name:   file.Name(),
		Reader: file,
	}

	msg, err := session.ChannelFileSend(channelIDUploadImages, imageFile.Name, imageFile.Reader)
	if err != nil {
		log.Printf("error uploading file: %v", err)
		return "", err
	}

	if len(msg.Attachments) > 0 {
		fmt.Println("Image uploaded successfully, URL:", msg.Attachments[0].URL)
		return msg.Attachments[0].URL, nil
	} else {
		return "", fmt.Errorf("image uploaded, but no attachments found")
	}
}

func PostNewShoe(shoe stockx.ProductDetails) error {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := session.Open()
	if err != nil {
		return fmt.Errorf("cannot open the session: %v", err)
	}
	path := "img_data/" + shoe.ProductName + "/gif/" + shoe.ProductName + ".gif"
	discordImageUrl, err := uploadLocalImage(path)
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

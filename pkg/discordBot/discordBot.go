package discordBot

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"vertigo/pkg/dataitems"

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
		return "", fmt.Errorf("Image uploaded, but no attachments found.")
	}
}

func PostNewShoe(shoe dataitems.Shoe) error {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := session.Open()
	if err != nil {
		return fmt.Errorf("Cannot open the session: %v", err)
	}

	discordImageUrl, err := uploadImage(shoe.ImageUrl)
	if err != nil {
		return fmt.Errorf("Cannot upload the image to Discord: %v", err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       shoe.Name,
		Description: fmt.Sprintf("Brand: %s\nSilhouette: %s\nTags: %s", shoe.Brand, shoe.Silhouette, shoe.Tags),
		Image: &discordgo.MessageEmbedImage{
			URL: discordImageUrl,
		},
		Color: 0x4c00b0,
	}

	_, err = session.ChannelMessageSendEmbed(channelIDShoeUpdates, embed)
	if err != nil {
		return fmt.Errorf("Cannot send the embedded message: %v", err)
	}

	err = session.Close()
	if err != nil {
		return fmt.Errorf("Cannot close the session: %v", err)
	}

	return nil
}

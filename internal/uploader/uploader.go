package uploader

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web %v", err)
	}
	return tok, nil
}

func uploadVideo(service *youtube.Service, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}
	defer file.Close()

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       "Test Video Title",
			Description: "This is a test video uploaded via the YouTube API",
			CategoryId:  "22",
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "public",
		},
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, video)
	response, err := call.Media(file).Do()
	if err != nil {
		log.Fatalf("Error making API call to upload the video: %v", err)
	}

	fmt.Printf("Upload successful! Video ID: %s\n", response.Id)
}

func Uploader(secretsPath string, videoPath string) {
	ctx := context.Background()

	// Load client secrets from a local file.
	b, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Configures OAuth2.0
	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Creates a new HTTP client, and logs in for authorization.
	token, err := tokenFromWeb(config)
	if err != nil {
		log.Fatalf("Unable to get token from web: %v", err)
	}

	client := config.Client(ctx, token)

	// Create the YouTube service
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Call the uploadVideo function with the created service and provided video path
	uploadVideo(service, videoPath)
}

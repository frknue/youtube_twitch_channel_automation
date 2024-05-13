package automater

import (
	"encoding/json"
	"log"
	"os"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/db"
	"github.com/frknue/youtube_twitch_channel_automation/internal/downloader"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
	"github.com/frknue/youtube_twitch_channel_automation/internal/uploader"
	"github.com/frknue/youtube_twitch_channel_automation/internal/video"
	"github.com/google/uuid"
)

type Video struct {
	RunID            string
	GameID           string
	Config           *config.Config
	ClipsData        []scraper.Clip
	VideoTitle       string
	VideoDescription string
	VideoEpisode     int
	VideoTags        []string
	VideoCategory    string
}

func Automate() {
	log.Println("Starting the application...")
	// Create lock
	if err := db.CreateLock(); err != nil {
		log.Fatalf("Failed to create lock: %v", err)
	}
	defer db.RemoveLock() // Ensure lock is always removed

	configPath := projectpath.Root + "/configs/config.yaml"
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	runID := uuid.New().String()
	gameID := config.Scraper.TwitchTracker.GameID

	cliPath := projectpath.Root + "/bin/TwitchDownloaderCLI"
	outputDir := projectpath.Root + config.Downloader.OutputPath + runID
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	clipsData, err := scraper.Scrape(outputDir, runID)
	if err != nil {
		log.Fatalf("Scraper failed with error: %v", err)
	}

	if len(clipsData) == 0 {
		log.Fatalf("No clips found.")
	}

	if err := downloader.Downloader(runID, clipsData, cliPath, outputDir); err != nil {
		log.Fatalf("Downloader failed with error: %v", err)
	}

	outputFile := outputDir + "/" + runID + ".mp4"
	if err := video.VideoCreator(clipsData, outputFile); err != nil {
		log.Fatalf("Video concatenation failed: %v", err)
	}

	// Create YouTube video title
	videoTitle, videoEpisode, err := video.CreateVideoTitle(gameID)
	// Create YouTube video description
	videoDescription := video.CreateVideoDescription(clipsData)
	// Create YouTube video tags
	videoTags := video.CreateVideoTags()
	// Create YouTube video category
	videoCategory := video.CreateVideoCategory()

	// Save the run data to a JSON file
	run := Video{RunID: runID, GameID: gameID, Config: config, ClipsData: clipsData, VideoTitle: videoTitle, VideoDescription: videoDescription, VideoEpisode: videoEpisode, VideoTags: videoTags, VideoCategory: videoCategory}
	jsonData, err := json.Marshal(run)
	if err != nil {
		log.Fatalf("Failed to marshal clips data: %v", err)
	}

	jsonFilePath := outputDir + "/run.json"
	if err := os.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to save JSON file: %v", err)
	}

	log.Printf("Saved %d clips to %s.\n", len(clipsData), jsonFilePath)

	//Upload the video to YouTube
	if err := uploader.Upload(outputFile, videoTitle, videoDescription, videoTags, videoCategory); err != nil {
		log.Fatalf("Failed to upload video: %v", err)
	}

	// Save clip IDs to the database
	for _, clip := range clipsData {
		if err := db.SaveClipID(clip.ClipID); err != nil {
			log.Fatalf("Database failed with error: %v", err)
		}
	}

	// Save run data to the database
	if err := db.SaveVideoData(run); err != nil {
		log.Fatalf("Failed to save video data: %v", err)
	}
}

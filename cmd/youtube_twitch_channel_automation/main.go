package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/db"
	"github.com/frknue/youtube_twitch_channel_automation/internal/downloader"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
	"github.com/frknue/youtube_twitch_channel_automation/internal/video"
	"github.com/google/uuid"
)

type Run struct {
	RunID     string
	Config    *config.Config
	ClipsData []scraper.Clip
	Bio       string
}

func main() {
	log.Println("Starting the application...")
	// Create lock
	if err := db.CreateLock(); err != nil {
		log.Fatalf("Failed to create lock: %v", err)
	}
	defer db.RemoveLock() // Ensure lock is always removed

	runID := uuid.New().String()
	configPath := projectpath.Root + "/configs/config.yaml"
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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

	bio, err := video.CreateYoutubeBioText(clipsData)
	if err != nil {
		log.Fatalf("Failed to create YouTube bio: %v", err)
	}

	for _, clip := range clipsData {
		if err := db.SaveClipID(clip.ClipID); err != nil {
			log.Fatalf("Database failed with error: %v", err)
		}
	}

	run := Run{RunID: runID, Config: config, ClipsData: clipsData, Bio: bio}
	jsonData, err := json.Marshal(run)
	if err != nil {
		log.Fatalf("Failed to marshal clips data: %v", err)
	}

	jsonFilePath := outputDir + "/run.json"
	if err := os.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to save JSON file: %v", err)
	}
	log.Printf("Saved %d clips to %s.\n", len(clipsData), jsonFilePath)
}

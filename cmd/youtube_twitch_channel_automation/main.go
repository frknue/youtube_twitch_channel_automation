package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
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
}

func main() {
	log.Println("Starting the application...")
	// Create a unique run ID using the UUID library
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

	log.Printf("Run ID: %s\n", runID)
	clipsData, err := scraper.Scrape(outputDir)
	if err != nil {
		log.Fatalf("Scraper failed with error: %v", err)
	}

	if len(clipsData) == 0 {
		log.Println("No clips found.")
		return
	}

	// Add the run ID to the Run struct and save it to the output directory as a JSON file
	run := Run{
		RunID:     runID,
		Config:    config,
		ClipsData: clipsData,
	}

	// Save run data to the output directory as a JSON file
	jsonData, err := json.Marshal(run)
	if err != nil {
		log.Fatalf("Failed to marshal clips data: %v", err)
	}
	jsonFilePath := outputDir + "/run.json"
	if err := os.WriteFile(jsonFilePath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to save JSON file: %v", err)
	}

	log.Printf("Saved %d clips to %s.\n", len(clipsData), jsonFilePath)
	log.Println("Starting the download process...")
	err = downloader.Downloader(runID, clipsData, cliPath, outputDir)
	if err != nil {
		log.Fatalf("Downloader failed with error: %v", err)
	}

	log.Println("Download process completed successfully.")

	outputFile := outputDir + "/" + runID + ".mp4"
	// err = video.ConcatenateVideos(videoFiles, outputFile)
	err = video.VideoCreator(clipsData, outputFile)

	if err != nil {
		log.Fatalf("Video concatenation failed: %v", err)
	}
	log.Println("Video concatenation completed successfully.")
}

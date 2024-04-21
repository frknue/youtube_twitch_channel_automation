package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/downloader"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
	"github.com/google/uuid"
)

func main() {
	log.Println("Starting the application...")
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
	clipsData, err := scraper.Scrape()
	if err != nil {
		log.Fatalf("Scraper failed with error: %v", err)
	}

	if len(clipsData) == 0 {
		log.Println("No clips found.")
		return
	}

	// Save clips data to the output directory as a JSON file
	jsonData, err := json.Marshal(clipsData)
	if err != nil {
		log.Fatalf("Failed to marshal clips data: %v", err)
	}
	jsonFilePath := outputDir + "/clipsData.json"
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
}

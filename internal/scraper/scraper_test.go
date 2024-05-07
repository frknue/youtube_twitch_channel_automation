package scraper

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/google/uuid"
)

type Run struct {
	RunID     string
	Config    *config.Config
	ClipsData []Clip
}

func TestScraper(t *testing.T) {
	runID := uuid.New().String()

	configPath := projectpath.Root + "/configs/config.yaml"

	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	outputDir := projectpath.Root + config.Downloader.OutputPath + runID

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	clipsData, err := Scrape(outputDir, runID)

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
}

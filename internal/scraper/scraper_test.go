package scraper

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/google/uuid"
)

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

	clipsData, err := Scrape(outputDir)

	if err != nil {
		log.Fatalf("Scraper failed with error: %v", err)
	}

	if len(clipsData) == 0 {
		log.Println("No clips found.")
		return
	}
	fmt.Println(clipsData)
}

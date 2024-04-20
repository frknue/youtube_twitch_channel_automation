package main

import (
	"log"

	"github.com/google/uuid"

	"github.com/frknue/youtube_twitch_channel_automation/internal/downloader"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
)

func main() {
	log.Println("Starting the application...")
	runID := uuid.New().String()
	log.Printf("Run ID: %s\n", runID)
	clipsData, err := scraper.Scrape()
    if err != nil {
        log.Fatalf("Scraper failed with error: %v", err)
    }

	if len(clipsData) == 0 {
		log.Println("No clips found.")
		return
	}

	log.Printf("Found %d clips.\n", len(clipsData))
	log.Println("Starting the download process...")
	err = downloader.Downloader(runID, clipsData)
	if err != nil {
		log.Fatalf("Downloader failed with error: %v", err)
	}
	log.Println("Download process completed successfully.")
}

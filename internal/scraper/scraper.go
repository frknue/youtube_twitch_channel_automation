package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/chromedp/chromedp"
	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
)

type Clip struct {
	RunID           string
	Channel         string
	Title           string
	URL             string
	Thumbnail       string
	Duration        string
	Views           string
	Created         string
	ClipID          string
	FileName        string
	DurationSeconds float64
	FilePath        string
}

func getHTML(url string) (string, error) {
	// Define the maximum number of retries
	const maxRetries = 5
	var html string
	var err error

	// Create context with non-headless options and more realistic parameters
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
		chromedp.Headless,
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
		// Add security flags
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	}

	for i := 0; i < maxRetries; i++ {
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		// Create a new context from the allocator
		ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancel()

		// Try to navigate to the URL and get the HTML content
		if err = chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("#clips-period", chromedp.ByQuery),
			chromedp.OuterHTML("html", &html),
		); err != nil {
			log.Printf("Attempt %d failed: %v", i+1, err)
			cancelCtx() // Ensure we cancel the context to clean up resources
			continue
		}

		// Check if the HTML contains necessary data
		if strings.Contains(html, "clips") {
			cancelCtx() // Ensure we cancel the context after a successful fetch
			return html, nil
		} else {
			log.Printf("Attempt %d fetched HTML but didn't contain necessary data.", i+1)
			cancelCtx() // Ensure we cancel the context to clean up resources
			continue
		}
	}

	// Return the last error after exhausting all retries
	return "", fmt.Errorf("failed to fetch valid HTML after %d attempts: %v", maxRetries, err)
}

func getClipsData(html string, maxDurationInSeconds float64, outputDir string, runID string) []Clip {
	var clips []Clip
	var totalDuration float64

	// Create a new document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		// Handle error
		fmt.Println("Error creating document:", err)
		return clips
	}

	// Get the element with the id "clips" and find each "clip-entity" within it
	doc.Find("#clips .clip-entity").Each(func(i int, s *goquery.Selection) {
		if totalDuration >= maxDurationInSeconds {
			fmt.Println("Total duration limit reached.")
			return
		}
		// Extract details for each clip
		title := s.Find(".clip-title").Text()
		fileName := ExtractFileName(title)
		url, exists := s.Find(".clip-tp").Attr("data-litebox")
		if !exists {
			url = "URL not found"
		}

		clipID := ExtractClipID(url)

		thumbnail, thumbExists := s.Find(".clip-thumbnail").Attr("src")
		if !thumbExists {
			thumbnail = "Thumbnail not found"
		}
		channel := s.Find(".clip-channel-name").Text()
		duration := s.Find(".clip-duration").Text()
		views := s.Find(".clip-views").Text()
		created := s.Find(".clip-created").Text()
		durationSeconds := ParseDuration(duration)

		totalDuration += durationSeconds

		filePath := outputDir + "/" + clipID + ".mp4"

		// Append the details to the clips slice
		clips = append(clips, Clip{
			RunID:           runID,
			Channel:         channel,
			Title:           title,
			URL:             url,
			Thumbnail:       thumbnail,
			Duration:        duration,
			Views:           views,
			Created:         created,
			ClipID:          clipID,
			FileName:        fileName + ".mp4",
			DurationSeconds: durationSeconds,
			FilePath:        filePath,
		})
	})
	log.Printf("Total duration of clips: %.2f seconds\n", totalDuration)

	return clips
}

func Scrape(outputDir string, runID string) ([]Clip, error) {
	configPath := projectpath.Root + "/configs/config.yaml"
	config, err := config.LoadConfig(configPath)

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return nil, err
	}

	url := config.Scraper.TwitchTracker.URL +
		"/" +
		config.Scraper.TwitchTracker.GameID +
		"/clips#" +
		config.Scraper.TwitchTracker.TimeStart +
		"-" + config.Scraper.TwitchTracker.TimeEnd

	html, err := getHTML(url)
	log.Println("Successfully fetched HTML data.")

	if err != nil {
		return nil, err
	}

	log.Println("Getting clips data...")
	clips := getClipsData(html, config.Downloader.MaxDurationInSeconds, outputDir, runID)
	log.Println("Successfully extracted clips data.")

	return clips, nil
}

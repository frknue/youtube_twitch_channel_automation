package video

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/frknue/youtube_twitch_channel_automation/internal/db"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
)

// CreateYoutubeBioText generates a YouTube bio text from clips data in JSON format.
func CreateVideoDescription(clipsData []scraper.Clip) (string, error) {
	// Greeting message.
	greeting := "Welcome to our YouTube channel!"

	// Create a set to collect unique Twitch channel URLs.
	channelURLs := make(map[string]bool)
	clipURLs := make(map[string]string) // A map to hold unique clip URLs and their titles

	for _, clip := range clipsData {
		channelURL := clip.ChannelURL
		channelURLs[channelURL] = true
		clipURLs[clip.URL] = clip.Title // Assuming clip.URL is the unique URL for each clip and clip.Title for the clip title
	}

	// Compile all unique Twitch channel URLs into a list.
	var urls []string
	for url := range channelURLs {
		urls = append(urls, url)
	}
	// Convert list of URLs to a single string for Twitch channels.
	urlsList := strings.Join(urls, "\n")

	// Compile clip URLs into a readable list format.
	var clipLinks []string
	for url, title := range clipURLs {
		clipLinks = append(clipLinks, fmt.Sprintf("%s - %s", title, url))
	}
	// Convert list of clip URLs to a single string.
	clipsList := strings.Join(clipLinks, "\n")

	// Call to subscribe and like.
	callToAction := "Don't forget to subscribe and hit that like button for more awesome content!"

	// Combine all parts into the final bio.
	description := fmt.Sprintf("%s\n\nChannels to check out:\n%s\n\nClips used in this video:\n%s\n\n%s", greeting, urlsList, clipsList, callToAction)

	// Convert the bio text to a JSON object.
	descriptionData := map[string]string{"description": description}
	descriptionJSON, err := json.Marshal(descriptionData)
	if err != nil {
		return "", err
	}

	return string(descriptionJSON), nil
}

// GetGameTitleByID reads the 'games_list.list' file and returns the game title by the gameID.
func GetGameTitleByID(gameID string) (string, error) {
	gamesListPath := projectpath.Root + "/configs/games_list.list"

	file, err := os.Open(gamesListPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		id := strings.TrimSpace(parts[0])
		title := strings.Trim(parts[1], "\"")

		if id == gameID {
			return title, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("game ID not found")
}

// CreateVideoTitle generates a video title based on the latest episode number.
func CreateVideoTitle(gameID string) (string, int, error) {
	// Get the game title by game ID
	gameTitle, err := GetGameTitleByID(gameID)

	e, err := db.GetLatestEpisodeByGameID(gameID)
	if err != nil {
		return gameTitle, 1, err
	}

	// Increment the episode number by 1
	episode := e + 1

	// Format the video title with the game title and episode number.
	videoTitle := fmt.Sprintf("%s MOST VIEWED Twitch Clips #%d", gameTitle, episode)

	return videoTitle, episode, nil
}

// Creating video Tags (static for now, but can be dynamic in the future)
func CreateVideoTags() []string {
	return []string{"twitch", "twitch clips", "twitch highlights", "gaming", "twitch moments"}
}

// Creating video category (static for now, but can be dynamic in the future)
func CreateVideoCategory() string {
	return "22"
}

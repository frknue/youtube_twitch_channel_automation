package video

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/frknue/youtube_twitch_channel_automation/internal/db"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
)

// CreateVideoDescription generates a video description text from clips data.
func CreateVideoDescription(clipsData []scraper.Clip) string {
	// Greeting message.
	greeting := "Welcome to our YouTube channel!"

	// Create a set to collect unique Twitch channel URLs.
	channelURLs := make(map[string]bool)
	var clipURLs []string // A list to hold unique clip URLs

	for _, clip := range clipsData {
		channelURL := clip.ChannelURL
		channelURLs[channelURL] = true
		clipURLs = append(clipURLs, clip.ClipURL) // Using ClipURL which is the correct field
	}

	// Compile all unique Twitch channel URLs into a list.
	var urls []string
	for url := range channelURLs {
		urls = append(urls, url)
	}
	// Convert list of URLs to a single string for Twitch channels.
	urlsList := strings.Join(urls, "\n")

	// Compile clip URLs into a readable list format.
	clipLinks := append([]string{"Clips used in this video:"}, clipURLs...)
	// Convert list of clip URLs to a single string.
	clipsList := strings.Join(clipLinks, "\n")

	// Call to subscribe and like.
	callToAction := "Don't forget to subscribe and hit that like button for more awesome content!"

	// Combine all parts into the final description.
	description := fmt.Sprintf("%s\n\nChannels to check out:\n%s\n\n%s\n\n%s", greeting, urlsList, clipsList, callToAction)

	return description
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
		return fmt.Sprintf("%s MOST VIEWED Twitch Clips #%d", gameTitle, e+1), e + 1, err
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
	return "20"
}

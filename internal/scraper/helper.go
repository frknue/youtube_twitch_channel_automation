package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"
)

// Helper function to extract clip ID from the URL
func ExtractClipID(clipURL string) string {
	// Parse the URL
	parsedURL, err := url.Parse(clipURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return ""
	}

	// Extract and decode query parameters
	queryParams, _ := url.ParseQuery(parsedURL.RawQuery)

	// Get the 'clip' parameter
	clipParam := queryParams.Get("clip")
	if clipParam == "" {
		return ""
	}

	// Handle possible nested URLs in 'clip' parameter by extracting the last valid segment
	clipParts := strings.Split(clipParam, "&clip=")
	lastPart := clipParts[len(clipParts)-1]
	// Trim any leading slashes
	finalClipID := strings.TrimLeft(lastPart, "/")

	// Return the entire final part after trimming, without splitting by '-'
	return finalClipID
}

func ExtractFileName(title string) string {
    // Trim any leading or trailing whitespaces
    title = strings.TrimSpace(title)

    // Replace invalid characters with underscore
    return strings.Map(func(r rune) rune {
        if r == '/' || r == '\\' || r == ':' || unicode.IsSpace(r) {
            return '_'
        }
        return r
    }, title)
}

func ParseDuration(durationStr string) float64 {
	duration, err := time.ParseDuration(strings.TrimSuffix(durationStr, "s") + "s")
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return 0
	}
	return duration.Seconds()
}
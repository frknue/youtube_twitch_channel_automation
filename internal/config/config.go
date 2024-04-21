package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds all the configuration for the application.
type Config struct {
	Scraper     ScraperConfig     `yaml:"scraper"`
	Downloader  DownloaderConfig  `yaml:"downloader"`
}

// ScraperConfig holds the configuration for the scraper part of the application.
type ScraperConfig struct {
	TwitchTracker TwitchTrackerConfig `yaml:"twitch_tracker"`
}

// TwitchTrackerConfig holds the specific configuration for tracking Twitch data.
type TwitchTrackerConfig struct {
	URL       string `yaml:"url"`
	GameID    string `yaml:"game_id"`
	TimeStart string `yaml:"time_start"`
	TimeEnd   string `yaml:"time_end"`
}

// DownloaderConfig holds the configuration for the downloader part of the application.
type DownloaderConfig struct {
	OutputPath  string `yaml:"output_path"`
	MaxDurationInSeconds float64 `yaml:"max_duration_in_seconds"`
}

// LoadConfig reads configuration from the specified file path and unmarshals it into Config struct.
func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

package downloader

import (
	"fmt"

	"os"
	"os/exec"

	"github.com/frknue/youtube_twitch_channel_automation/internal/config"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
)

func Downloader(runID string, clipData []scraper.Clip) error {
	configPath := projectpath.Root + "/configs/config.yaml"
	config, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cliPath := projectpath.Root + "/bin/TwitchDownloaderCLI"
	outputDir := projectpath.Root + config.Downloader.OutputPath + runID

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		fmt.Println("Output directory created successfully.")
	}

	// Loop through each clip in the provided slice
	for _, clip := range clipData {
		outputPath := fmt.Sprintf("%s/%s.mp4", outputDir, clip.ClipID)

		// Prepare the command with the specific clip ID and output path
		cmd := exec.Command(cliPath, "clipdownload", "--id", clip.ClipID, "--output", outputPath)

		// Execute the command and capture the output
		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error downloading clip %s: %v\n", clip.ClipID, err)
			fmt.Printf("Command output: %s\n", string(cmdOutput))
			continue // proceed with next clip despite the error
		}

		// Optionally print the successful download message
		fmt.Printf("Downloading clip %s completed successfully.\n", clip.ClipID)
		fmt.Printf("Command output: %s\n", string(cmdOutput))
	}

	return nil
}
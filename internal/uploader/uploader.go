package uploader

import (
	"fmt"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"os/exec"
	"path/filepath"
)

func Upload(filePath string, videoTitle string, videoDescription string, videoTags []string, videoCategory string) error {
	// Define the path to the Python executable inside the virtual environment
	pythonPath := filepath.Join(projectpath.Root, "scripts", "venv", "bin", "python")

	// Define the path to the Python script
	scriptPath := filepath.Join(projectpath.Root, "scripts", "upload_video.py")

	// join the video tags separated by commas
	tags := ""
	for _, tag := range videoTags {
		tags += tag + ","
	}

	// Prepare the arguments for the command
	args := []string{
		scriptPath, // First argument is always the script file name
		"--file", filePath,
		"--title", videoTitle,
		"--description", videoDescription,
		"--keywords", tags,
		"--category", videoCategory,
		"--privacyStatus", "private",
	}

	// Create the command with the Python interpreter and the arguments
	cmd := exec.Command(pythonPath, args...)

	// Running the command
	fmt.Println("Uploading video to YouTube...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error uploading video:", err)
		fmt.Println("Output:", string(output))
		return err
	}

	// Print the output from the Python script
	fmt.Println("Output from Python script:", string(output))
	return nil
}

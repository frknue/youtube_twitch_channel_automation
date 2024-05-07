package video

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/frknue/youtube_twitch_channel_automation/internal/scraper"
)

// ReencodeVideo re-encodes video files to a format optimized for YouTube, adds the channel name as a watermark, and overwrites the original.
func ReencodeVideo(inputFile, channelName string) error {
	tempOutputFile := inputFile + "_temp.mp4" // Temporary output file
	fmt.Printf("Re-encoding video %s to %s\n", inputFile, tempOutputFile)
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:v", "libx264", "-preset", "fast", "-crf", "22",
		"-vf", fmt.Sprintf("scale=1920:1080,drawtext=text='%s':x=w-tw-10:y=h-th-10:fontcolor=white@0.8:fontsize=24:box=1:boxcolor=black@0.5:boxborderw=5", channelName),
		"-c:a", "aac", "-b:a", "192k",
		tempOutputFile,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg re-encoding error: %s, output: %s", err, output)
	}

	if err := os.Rename(tempOutputFile, inputFile); err != nil {
		return fmt.Errorf("failed to overwrite original video file: %v", err)
	}
	return nil
}

// ConcatenateVideos concatenates a slice of Clip structs into a single file.
func ConcatenateVideos(clipsData []scraper.Clip, outputFile string) error {
	inputs := strings.Builder{}
	filterComplex := strings.Builder{}

	for i, clip := range clipsData {
		if err := ReencodeVideo(clip.FilePath, clip.Channel); err != nil {
			return err
		}
		inputs.WriteString(fmt.Sprintf("-i '%s' ", clip.FilePath))
		filterComplex.WriteString(fmt.Sprintf("[%d:v:0][%d:a:0]", i, i))
	}

	filterComplex.WriteString(fmt.Sprintf("concat=n=%d:v=1:a=1[v][a]", len(clipsData)))

	cmdString := fmt.Sprintf("ffmpeg %s -filter_complex '%s' -map '[v]' -map '[a]' -c:v libx264 -preset fast -crf 22 -c:a aac -b:a 192k '%s'", inputs.String(), filterComplex.String(), outputFile)
	cmd := exec.Command("bash", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat error: %s, output: %s", err, output)
	}

	fmt.Printf("Concatenation output:\n%s\n", output)
	return nil
}

func VideoCreator(clipsData []scraper.Clip, outputFile string) error {
	fmt.Println("Creating video...")
	if err := ConcatenateVideos(clipsData, outputFile); err != nil {
		return err
	}

	fmt.Println("Video created successfully.")
	return nil
}

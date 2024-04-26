package video

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func deleteTempFiles(files []string) {
	// Delete video files
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			fmt.Printf("Failed to delete temp file %s: %v\n", file, err)
		}
	}
	// Delete filelist.txt
	if err := os.Remove("filelist.txt"); err != nil {
		fmt.Printf("Failed to delete filelist.txt: %v\n", err)
	}
}

// ConcatenateVideos takes a slice of file paths and an output file path, and concatenates the videos into a single file.
func ReencodeVideo(inputFile, outputFile string) error {
	fmt.Printf("Re-encoding video %s to %s\n", inputFile, outputFile)
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:v", "libx264", "-crf", "23", "-preset", "fast",
		"-c:a", "aac", "-b:a", "192k",
		"-strict", "experimental",
		outputFile,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg re-encoding error: %s, output: %s", err, output)
	}
	return nil
}

func ConcatenateVideos(videoFiles []string, outputFile string) error {
	tempFileList := "filelist.txt"
	content := strings.Builder{}
	reencodedFiles := make([]string, len(videoFiles))

	for i, file := range videoFiles {
		outputFile := fmt.Sprintf("temp%d.mp4", i)
		if err := ReencodeVideo(file, outputFile); err != nil {
			return err
		}
		reencodedFiles[i] = outputFile
		content.WriteString(fmt.Sprintf("file '%s'\n", outputFile))
	}

	if err := os.WriteFile(tempFileList, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write file list: %v", err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", tempFileList,
		"-c", "copy",
		outputFile,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg concat error: %s, output: %s", err, output)
	}

	deleteTempFiles(reencodedFiles)

	fmt.Printf("Concatenation output:\n%s\n", output)
	return nil
}

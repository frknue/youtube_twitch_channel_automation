package uploader

import (
	"fmt"
	"github.com/frknue/youtube_twitch_channel_automation/internal/projectpath"
	"testing"
)

func TestUploader(t *testing.T) {
	secretsPath := projectpath.Root + "/client_secrets.json"
	videoPath := projectpath.Root + "/sample.mp4"
	fmt.Println(secretsPath)
	fmt.Println(videoPath)
	Uploader(secretsPath, videoPath)
}

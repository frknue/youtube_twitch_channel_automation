package video

import (
	"fmt"
	"testing"
)

func TestCreateVideoTitle(t *testing.T) {
	gameID := "538054672"

	videoTitle, videoEpisode, err := CreateVideoTitle(gameID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Video title:", videoTitle)
	fmt.Println("Video episode:", videoEpisode)
}

func TestGetGameTitleByID(t *testing.T) {
	gameID := "538054672"
	gameTitle, err := GetGameTitleByID(gameID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Game title:", gameTitle)
}

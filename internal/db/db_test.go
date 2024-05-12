package db

import (
	"fmt"
	"testing"
)

func TestSaveClipID(t *testing.T) {
	clipID := "test"
	err := SaveClipID(clipID)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCheckClipID(t *testing.T) {
	clipID := "test"
	clipIDExists := CheckClipID(clipID)
	if clipIDExists {
		fmt.Println("Clip ID exists")
	} else {
		fmt.Println("Clip ID does not exist")
	}
}

func TestCreateLock(t *testing.T) {
	err := CreateLock()
	if err != nil {
		fmt.Println(err)
	}
}

func TestRemoveLock(t *testing.T) {
	err := RemoveLock()
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetLatestVideoByGameID(t *testing.T) {
	gameID := "538054672"
	video, err := GetLatestVideoByGameID(gameID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(video)
}

func TestGetLatestEpisodeByGameID(t *testing.T) {
	gameID := "538054672"
	episode, err := GetLatestEpisodeByGameID(gameID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Latest episode:", episode)
}

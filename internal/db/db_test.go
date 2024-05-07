package db

import (
	"fmt"
	"testing"
)

func TestSaveClipID(t *testing.T) {
	clipID := "test"
	err := SaveClipID(clipID)
	fmt.Println("TestDatabase HELLOOOO")
	if err != nil {
		fmt.Println(err)
	}
}

func TestCheckClipID(t *testing.T) {
	clipID := "129e1"
	hasClip, err := CheckClipID(clipID)
	fmt.Println("hasClip: ", hasClip)
	if err != nil {
		fmt.Println(err)
	}
}

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
	clipID := "test"
	hasClip := CheckClipID(clipID)
	fmt.Println("hasClip: ", hasClip)
}

func TestPrintClipIDs(t *testing.T) {
	err := PrintClipIDs()
	if err != nil {
		fmt.Println(err)
	}
}

func TestCleanUpDB(t *testing.T) {
	err := CleanUpDB()
	if err != nil {
		fmt.Println(err)
	}
}

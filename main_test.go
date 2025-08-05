package main

import (
	"fmt"
	"os"
	"testing"
)

/// Test function to check if able to read individual image files inside the spritesheet
/// created by the json file
func TestReadImages(t *testing.T) {
	// Read json file to read location of image in the spritesheet
	content, err := os.ReadFile("spritesheet.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}
	fmt.Println("JSON content:", string(content))

}

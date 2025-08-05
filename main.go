package main

import (
	"path/filepath"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

type SerializedRect struct {
	X,Y,W,H float64
}
type SerializedPos struct {
	X,Y float64
}
type SerializedDim struct {
	W,H float64
}

type SerializedFrame struct {
	Rotated bool
	Trimmed bool
	Frame SerializedRect
	SpriteSourceSize SerializedRect
	SourceSize SerializedDim
	Pivot SerializedPos
}

type SerializedSpritesheet struct {
	ImageName string
	// Frames map[string]SerializedFrame
	// Meta map[string]interface{}
}

func main() {

	// User options to test read sprite or create a new spritesheet
	var choice int
	fmt.Println("Welcome to the Sprite Packer!")
	fmt.Println("1. Create new spritesheet")
	fmt.Println("2. Read existing spritesheet")
	fmt.Print("Enter your choice (1 or 2): ")

	_, err := fmt.Scanf("%d\n", &choice)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)	
	}

	switch choice {
		case 1:
			// Create the spritesheet
			fmt.Print("Enter the name for the spritesheet (without extension): ")
			var sprite_name string
			_, err := fmt.Scanln(&sprite_name)
			if err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
	
			fmt.Printf("Creating spritesheet with name: %s \n", sprite_name)
			CreateSpritesheet(sprite_name)
		case 2:
			//TODO: Implement reading the JSON file to get sprite data
			fmt.Println("Reading JSON file to get sprite data...")
	}
}

func CreateSpritesheet(sprite_name string) {
	// Create a new spritesheet from images in the specified directory
	flagVar := flag.String("path", "images", "Path to the image directory")
	flag.Parse()

	if _, err := os.Stat(*flagVar); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist.\n", *flagVar)
		return
	}
	
	// Load all images from the specified directory
	dir, err := os.ReadDir(*flagVar)
	if err != nil {panic(err)}

	// TODO: Create the Sprite to pack into
	offset_x := 0
	offset_y := 0
	sprite := image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	fmt.Println("Found images:", dir)
	files, err := filepath.Glob(*flagVar + "/*.png")
	if err != nil {panic(err)}

	// Iterate through the files and read each image
	for path := range files {
		img, err := ReadImage(files[path])
		if err != nil {
			fmt.Printf("Error reading:  %v:", err)
			continue
		}

		// Draw the image onto the sprite at the current offset
		draw.Draw(
			sprite, 
			image.Rect(
				offset_x, 
				offset_y, 
				offset_x + img.Bounds().Dx(),
				offset_y + img.Bounds().Dy(),
			), 
			img,
			image.Point{X: 0, Y: 0},
			draw.Over,
	)

		// Update the offsets for the next image
		offset_x += img.Bounds().Dx()
		if offset_x >= sprite.Bounds().Dx() {
			offset_x = 0
			offset_y += img.Bounds().Dy()
		}
	}

	// Save the sprite to a file
	outFile, err := os.Create(sprite_name + ".png")
	if err != nil {
		fmt.Printf("Error creating sprite file: %v\n", err)
		return
	}
	defer outFile.Close()

	// Encode the sprite image to PNG format
	if err := png.Encode(outFile, sprite); err != nil {
		fmt.Printf("Error encoding sprite image: %v\n", err)
	}
}

func LoadJson() (SerializedSpritesheet, error) {
	// TODO: Implement JSON serialization of the sprite data
	fmt.Println("Creating JSON data for the spritesheet...")
	if _, err := os.Stat("spritesheet.json"); err == nil {
		fmt.Println("spritesheet.json already exists, loading data...")
		json_data, err := ReadJson("spritesheet.json")
		return json_data, err
	}

	// Create the new json file
	file, err := os.Create("spritesheet.json")
	if err != nil {
		fmt.Printf("Error creating JSON file: %v\n", err)
		return SerializedSpritesheet{}, err
	}

	fmt.Println("JSON file created successfully:", file.Name())
	defer file.Close()

	// Read the json file to get the sprite data
	json_data, err := ReadJson("spritesheet.json")
	return json_data, err
}

func ReadJson(path string) (SerializedSpritesheet, error) {
	/// Read the JSON file to get the sprite data
	file, err := os.ReadFile("spritesheet.json")
	if err != nil {
		fmt.Printf("Error reading existing JSON file: %v\n", err)
		return SerializedSpritesheet{}, err
	}

	var payload SerializedSpritesheet
	err = json.Unmarshal(file, &payload)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON data: %v\n", err)
		return SerializedSpritesheet{}, err
	}
	return payload, nil
}


func ReadImage(path string) (image.Image, error) {
	// Read the image
	fmt.Println("Reading image:", path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening image file %s: %v\n", path, err)
		return nil, err
	}
	defer file.Close()

	// Read the image data
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image file %s: %v\n", path, err)
		return nil, err
	}
	return img, nil
}

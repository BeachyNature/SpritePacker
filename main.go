package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
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
}

type SerializedSpritesheet struct {
	Name string
	Frames map[string] SerializedFrame
}

func main() {

	// User options to test read sprite or create a new spritesheet
	var choice int

	for {
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
				return

			case 2:
				// Read the Json file then read the spritesheet
				fmt.Println("Reading JSON file to get sprite data...")
				sheet := ReadJson("spritesheet.json")
				fmt.Printf("Sheet: %v", sheet["example"])
				return

			default:
				fmt.Println("Invalid choice, please run the program again and select 1 or 2.")
		}
	}
}

// Read in the json file
func ReadSpriteSheet(path string) error {
	// TODO: Implement reading the spritesheet image and extracting frames based on JSON data
	return nil
}

// Create a new spritesheet from images in the specified directory
func CreateSpritesheet(sprite_name string) {
	flagVar := flag.String("path", "images", "Path to the image directory")
	flag.Parse()

	if _, err := os.Stat(*flagVar); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist.\n", *flagVar)
		return
	}
	
	// Load all images from the specified directory
	dir, err := os.ReadDir(*flagVar)
	if err != nil {panic(err)}

	offset_x := 0
	offset_y := 0
	sprite := image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	fmt.Println("Found images:", dir)
	files, err := filepath.Glob(*flagVar + "/*.png")
	if err != nil {panic(err)}

	// Iterate through the files and read each image
	sheet := SerializedSpritesheet{
		Name: sprite_name,
		Frames: make(map[string]SerializedFrame),
	}

	for path := range files {
		img, err := ReadImage(files[path])
		if err != nil {
			fmt.Printf("Error reading:  %v:", err)
			continue
		}
		fmt.Printf("Image %d: %s\n", path, files[path])

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

		// TODO: Add logic to handle wrapping to next row
		offset_x += img.Bounds().Dx()
		if offset_x >= sprite.Bounds().Dx() {
			offset_x = 0
			offset_y += img.Bounds().Dy()
		}
		
		// Create dicitonary to store frame data
		var file_name = filepath.Base(files[path])
		frame_loc := SerializedFrame {
			Rotated: false,
			Trimmed: false,
			Frame: SerializedRect {
				X: float64(offset_x),
				Y: float64(offset_y),
				W: float64(img.Bounds().Dx()),
				H: float64(img.Bounds().Dy()),
			},
		}
		sheet.Frames[file_name] = frame_loc
	}

	// Store the sheet into a json with framed values
	sprite_dict := make(map[string] SerializedSpritesheet)
	sprite_dict[sprite_name] = sheet

	// Write the JSON file with the sprite data
	err = WriteJson(sprite_dict)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
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

// Write the JSON file to store the sprite data
func WriteJson(data map[string] SerializedSpritesheet) error {
	if _, err := os.Stat("spritesheet.json"); err == nil {
		fmt.Println("spritesheet.json already exists, loading data...")
		// TODO: Add data to existing file
		return nil
	}

	// Create the new json file
	file, err := os.Create("spritesheet.json")
	if err != nil {
		fmt.Printf("Error creating JSON file: %v\n", err)
		return err
	}
	
	// Encode the data to JSON format
	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")
	enc.Encode(data)

	defer file.Close()
	return nil
}

// Read the JSON file to get the sprite data
func ReadJson(path string) map[string]SerializedSpritesheet{
	file, err := os.ReadFile("spritesheet.json")
	if err != nil {
		fmt.Printf("Error reading existing JSON file: %v\n", err)
		return nil
	}

	// Read json into Spritesheet map
	var payload map[string]SerializedSpritesheet
	err = json.Unmarshal(file, &payload)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON data: %v\n", err)
		return nil
	}
	return payload
}

// Open and read an image file
func ReadImage(path string) (image.Image, error) {
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

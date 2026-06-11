package main

import (
	"embed"
	"fmt"

	p "streamdeckVR/internal/profiles"
)

//go:embed all:frontend/dist
var assets embed.FS

type Profile struct {
	UUID  string
	name  string
	pages int
}

func main() {
	var path string
	path = "C:\\Users\\mfleu\\AppData\\Roaming\\Elgato\\StreamDeck\\bProfilesV3"
	files, err := p.ProfileScan(path)

	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Iterate and print basic info
	for _, entry := range files {
		fmt.Printf(
			"Name: %s, IsDir: %v, Type: %v\n",
			entry.Name(),
			entry.IsDir(),
			entry.Type(),
		)
	}
	// Create an instance of the app structure
	// app := NewApp()

	// // Create application with options
	// err := wails.Run(&options.App{
	// 	Title:  "Change",
	// 	Width:  1024,
	// 	Height: 768,
	// 	AssetServer: &assetserver.Options{
	// 		Assets: assets,
	// 	},
	// 	BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
	// 	OnStartup:        app.startup,
	// 	Bind: []interface{}{
	// 		app,
	// 	},
	// })

	// if err != nil {
	// 	println("Error:", err.Error())
	// }
}

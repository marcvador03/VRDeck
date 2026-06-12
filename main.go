package main

import (
	"embed"
	"fmt"
	"os"

	l "streamdeckVR/internal/logger"
	p "streamdeckVR/internal/profiles"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var path string
	path = "C:\\Users\\mfleu\\AppData\\Roaming\\Elgato\\StreamDeck\\bProfilesV3"
	log, err := l.NewLogger("./log.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Logger Error] %v\n", err)
		os.Exit(1)
	}
	ProfileList, err := p.NewProfileList(path, log)
	if err != nil {
		log.Error(err)
	}
	fmt.Println(ProfileList.Profiles)

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

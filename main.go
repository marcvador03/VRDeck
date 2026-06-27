package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"streamdeckVR/internal/logger"
	"streamdeckVR/internal/profiles"
	"streamdeckVR/internal/scanner"

	"go.uber.org/zap/zapcore"
)

//go:embed all:frontend/dist
var assets embed.FS

func inspect_data(profileList profiles.ProfileList) {

	for i, profile := range profileList.Profiles {
		fmt.Printf("Profile %d: %s (UUID: %s)\n", i, profile.Name, profile.UUID)
		for j, page := range profile.Pages {
			fmt.Printf("  Page %d: %s (UUID: %s)\n", j, page.Name, page.UUID)
			for k, button := range page.Buttons {
				fmt.Printf(
					"    Action %d: Row=%d, Col=%d, Title=%s, ActionID=%s\n",
					k, button.Row, button.Col, button.Title, button.ActionID,
				)
			}
		}
	}
}

func main() {
	path := filepath.Join(os.Getenv("APPDATA"), "Elgato", "StreamDeck", "bProfilesV3")
	logger.InitLogger("log.txt", zapcore.DebugLevel)
	log := logger.GetDefaultLogger()
	defer log.Sync()
	ProfileList, err := profiles.NewProfileList(path)
	if err != nil {
		log.Error("Stopping, error encountered")
		os.Exit(1)
	}
	inspect_data(*ProfileList)
	scanner.StartProfilesScan(path)

	//wails default code
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

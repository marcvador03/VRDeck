package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"streamdeckVR/internal/logger"
	p "streamdeckVR/internal/profiles"

	c "streamdeckVR/internal/connection"
	p "streamdeckVR/internal/profiles"

	"go.uber.org/zap/zapcore"

	"github.com/gorilla/websocket"
)

//go:embed all:frontend/dist
var assets embed.FS

func inspect_data(profileList p.ProfileList) {

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	path := filepath.Join(os.Getenv("APPDATA"), "Elgato", "StreamDeck", "bProfilesV3")
	logger.InitLogger("app.log", zapcore.DebugLevel)
	log := logger.GetDefaultLogger()
	defer log.Sync()
	ProfileList, err := p.NewProfileList(path)
	if err != nil {
		log.Error("Stopping, error encountered")
		os.Exit(1)
	}
	inspect_data(*ProfileList)

	c.InitiateStreamDeckConnection(log)

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

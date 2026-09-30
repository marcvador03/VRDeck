package main

import (
	"embed"
	"os"
	"path/filepath"

	"VRDeck/internal/logger"
	"VRDeck/internal/profiles"
	"VRDeck/internal/ws"

	"go.uber.org/zap/zapcore"
)

//go:embed all:frontend/dist

var assets embed.FS
var Version = "dev"

var ElgatoProgramPath = filepath.Join(os.Getenv("ProgramFiles"), "Elgato", "StreamDeck", "PageIcon")
var ElgatoProfilePath = filepath.Join(os.Getenv("APPDATA"), "Elgato", "StreamDeck", "ProfilesV3")

func main() {
	//path := filepath.Join(os.Getenv("APPDATA"), "Elgato", "StreamDeck", "ProfilesV3")
	logger.InitLogger("log.txt", zapcore.DebugLevel)
	log := logger.GetDefaultLogger()
	defer log.Sync()
	wsServer := ws.NewMSFSWebSocket()
	wsServer.CreateWebSocket()
	profilelist := profiles.NewProfileList(ElgatoProfilePath, wsServer)
	profilelist.CreateProfileList()
	//profilelist.InspectData()
	profilelist.StartProfilesScan()

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

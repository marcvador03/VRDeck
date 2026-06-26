package connection

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"streamdeckVR/internal/logger"
	"time"

	"github.com/coder/websocket"
	"go.uber.org/zap"
)

type params struct {
	port          int
	info          string
	pluginUUID    string
	registerEvent string
}

type DeviceType int

const (
	DeviceTypeStreamDeck        DeviceType = 0  // Stream Deck, Stream Deck Scissor Keys
	DeviceTypeStreamDeckMini    DeviceType = 1  // Stream Deck Mini
	DeviceTypeStreamDeckXL      DeviceType = 2  // Stream Deck XL
	DeviceTypeStreamDeckMobile  DeviceType = 3  // Stream Deck Mobile
	DeviceTypeCorsairGKeys      DeviceType = 4  // Corsair GKeys
	DeviceTypeStreamDeckPedal   DeviceType = 5  // Stream Deck Pedal
	DeviceTypeCorsairVoyager    DeviceType = 6  // Corsair Voyager
	DeviceTypeStreamDeckPlus    DeviceType = 7  // Stream Deck +
	DeviceTypeSCUFController    DeviceType = 8  // SCUF Controller
	DeviceTypeStreamDeckNeo     DeviceType = 9  // Stream Deck Neo
	DeviceTypeStreamDeckStudio  DeviceType = 10 // Stream Deck Studio
	DeviceTypeVirtualStreamDeck DeviceType = 11 // Virtual Stream Deck
	DeviceTypeGalleon100SD      DeviceType = 12 // Galleon 100 SD
	DeviceTypeStreamDeckPlusXL  DeviceType = 13 // Stream Deck + XL
)

type RegisterEvent struct {
	event string // -registerEvent parameter
	uuid  string // -pluginUUID parameter
}

type RegistrationInfo struct {
	Application struct {
		Font            string `json:"font"`
		Language        string `json:"language"` // e.g., "de", "en", "es", "fr", "ja", "ko", "zh_CN", "zh_TW"
		Platform        string `json:"platform"` // e.g., "mac", "windows"
		PlatformVersion string `json:"platformVersion"`
		Version         string `json:"version"`
	} `json:"application"`

	Colors struct {
		ButtonMouseOverBackgroundColor string `json:"buttonMouseOverBackgroundColor"`
		ButtonPressedBackgroundColor   string `json:"buttonPressedBackgroundColor"`
		ButtonPressedBorderColor       string `json:"buttonPressedBorderColor"`
		ButtonPressedTextColor         string `json:"buttonPressedTextColor"`
		HighlightColor                 string `json:"highlightColor"`
	} `json:"colors"`

	DevicePixelRatio float64 `json:"devicePixelRatio"`

	Devices []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Size struct {
			Columns int `json:"columns"`
			Rows    int `json:"rows"`
		} `json:"size"`
		Type DeviceType `json:"type"`
	} `json:"devices"`

	Plugin struct {
		UUID    string `json:"uuid"`
		Version string `json:"version"`
	} `json:"plugin"`
}

func getConnectionParams() (params, error) {
	log := logger.GetDefaultLogger()
	var p params

	flag.IntVar(&p.port, "port", 0, "WebSocket port")
	flag.StringVar(&p.info, "info", "", "Registration info (JSON)")
	flag.StringVar(&p.pluginUUID, "pluginUUID", "", "Plugin UUID")
	flag.StringVar(&p.registerEvent, "registerEvent", "", "Register event name")
	flag.Parse()

	if p.port == 0 || p.pluginUUID == "" || p.registerEvent == "" || p.info == "" {
		log.Error("Invalid command line arguments passed",
			zap.Int("Port", p.port),
			zap.String("pluginUUID", p.pluginUUID),
			zap.String("registerEvent", p.registerEvent),
			zap.String("info", p.info))
		return p, fmt.Errorf("error")
	}
	var regInfo RegistrationInfo
	if err := json.Unmarshal([]byte(p.info), &regInfo); err != nil {
		log.Error(("Failed to parse RegistrationInfo"),
			zap.Error(err))
	}
	log.Info("Received RegistrationInfo from command-line",
		zap.Any("registration_info", regInfo),
	)
	return p, nil
}

func connecttoStreamDeck(p params) error {
	log := logger.GetDefaultLogger()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	addr := "ws://localhost:" + strconv.Itoa(p.port)
	fmt.Printf(addr)
	c, _, err := websocket.Dial(ctx, addr, nil)
	defer c.CloseNow()
	if err != nil {
		log.Error("Error during websocket connection",
			zap.Int("Port", p.port),
			zap.String("pluginUUID", p.pluginUUID),
			zap.String("registerEvent", p.registerEvent),
			zap.Error(err))
		return err
	}
	regMsg := RegisterEvent{
		event: p.registerEvent,
		uuid:  p.pluginUUID,
	}
	jsonData, err := json.Marshal(regMsg)
	if err != nil {
		log.Error("Failed to marshal registration message",
			zap.Error(err))
	}

	if err := c.Write(ctx, websocket.MessageText, jsonData); err != nil {
		log.Error("Failed to register plugin: %v",
			zap.Error(err))
	} else {
		log.Info("Successfully Connected to StreamDeck VR",
			zap.Int("port", p.port))
	}
	c.Close(websocket.StatusNormalClosure, "")
	return nil
}

func InitiateStreamDeckConnection() error {

	p, err := getConnectionParams()
	if err != nil {
		return err
	}
	if err := connecttoStreamDeck(p); err != nil {
		return err
	}
	return nil
}

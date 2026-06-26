package connection

import (
	"context"
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

func getConnectionParams() (params, error) {
	log := logger.GetDefaultLogger()
	var p params

	flag.IntVar(&p.port, "port", 0, "WebSocket port")
	flag.StringVar(&p.info, "info", "", "Registration info (JSON)")
	flag.StringVar(&p.pluginUUID, "pluginUUID", "", "Plugin UUID")
	flag.StringVar(&p.registerEvent, "registerEvent", "", "Register event name")
	flag.Parse()

	if p.port == 0 || p.pluginUUID == "" || p.registerEvent == "" {
		log.Error("Invalid command line arguments passed",
			zap.Int("Port", p.port),
			zap.String("pluginUUID", p.pluginUUID),
			zap.String("registerEvent", p.registerEvent))
		return p, fmt.Errorf("error")
	}
	return p, nil
}

func connecttoStreamDeck(p params) error {
	log := logger.GetDefaultLogger()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	addr := "ws://localhost:" + strconv.Itoa(p.port)
	fmt.Printf(addr)
	c, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		log.Error("Error during websocket connection",
			zap.Int("Port", p.port),
			zap.String("pluginUUID", p.pluginUUID),
			zap.String("registerEvent", p.registerEvent),
			zap.Error(err))
		return err
	}
	defer c.CloseNow()

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

package connection

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	l "streamdeckVR/internal/logger"
	"time"

	"github.com/coder/websocket"
)

type params struct {
	port          int
	info          string
	pluginUUID    string
	registerEvent string
}

func getConnectionParams() (params, error) {
	var p params

	flag.IntVar(&p.port, "port", 0, "WebSocket port")
	flag.StringVar(&p.info, "info", "", "Registration info (JSON)")
	flag.StringVar(&p.pluginUUID, "pluginUUID", "", "Plugin UUID")
	flag.StringVar(&p.registerEvent, "registerEvent", "", "Register event name")
	flag.Parse()

	if p.port == 0 || p.pluginUUID == "" || p.registerEvent == "" {
		return params{}, fmt.Errorf("Invalid command line arguments passed")
	}
	return p, nil
}

func connecttoStreamDeck(p params) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	addr := "ws//localhost:" + strconv.Itoa(p.port)
	c, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return fmt.Errorf("Connection Error: %w", err)
	}
	defer c.CloseNow()

	c.Close(websocket.StatusNormalClosure, "")
	return nil
}

func InitiateStreamDeckConnection(log *l.Logger) error {
	p, err := getConnectionParams()
	if err != nil {
		log.Error(err)
		return err
	}
	if err := connecttoStreamDeck(p); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

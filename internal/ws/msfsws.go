package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"streamdeckVR/internal/logger"
	"sync"

	"github.com/coder/websocket"
	"go.uber.org/zap"
)

type MSFSWebSocket struct {
	clients   []*websocket.Conn
	clientsMu sync.Mutex
}

func NewMSFSWebSocket() *MSFSWebSocket {
	return &MSFSWebSocket{}
}

func (ws *MSFSWebSocket) BroadcastJSON(message any) error {
	log := logger.GetDefaultLogger()
	ws.clientsMu.Lock()
	defer ws.clientsMu.Unlock()

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	for _, client := range ws.clients {
		err := client.Write(context.Background(), websocket.MessageText, jsonData)
		if err != nil {
			log.Error("WebSocket write error",
				zap.Error(err))
			ws.removeClient(client)
		}
	}
	return nil
}

func (ws *MSFSWebSocket) removeClient(client *websocket.Conn) {
	log := logger.GetDefaultLogger()
	ws.clientsMu.Lock()
	for i, c := range ws.clients {
		if c == client {
			ws.clients = append(ws.clients[:i], ws.clients[i+1:]...)
			break
		}
	}
	ws.clientsMu.Unlock()
	log.Info("Client disconnected")
}

func (ws *MSFSWebSocket) CreateWebSocket() {
	log := logger.GetDefaultLogger()

	http.HandleFunc("/streamdeckvr", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Error(("Error while accepting websocket connections"),
				zap.Error(err))
			return
		}
		defer conn.Close(websocket.StatusInternalError, "server closed")
		ws.clientsMu.Lock()
		ws.clients = append(ws.clients, conn)
		ws.clientsMu.Unlock()
		log.Info("New client connected")

		// Keep the connection alive (no reading, just wait for disconnect)
		<-r.Context().Done()
		ws.removeClient(conn)
	})

	go func() {
		log.Info("WebSocket server running on ws://localhost:8080/streamdeckvr")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Error(("Server error:"),
				zap.Error(err))
		}
	}()
}

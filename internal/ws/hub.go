package ws

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/chinu-anand/crawlerx/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	JobIDs map[string]bool
}

type Hub struct {
	clients map[*Client]bool
	lock    sync.Mutex
}

var hub = Hub{
	clients: make(map[*Client]bool),
}

func (h *Hub) AddClient(client *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[client] = true
}

func (h *Hub) RemoveClient(client *Client) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.clients, client)
	client.Conn.Close()
}

func (h *Hub) Broadcast(id string, message []byte) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for client := range h.clients {
		if len(client.JobIDs) > 0 && !client.JobIDs[id] {
			continue
		}
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("❌ Error writing message:", err)
			h.RemoveClient(client)
		}
	}
}

func (h *Hub) BroadcastJobUpdate(id string, url string, status models.JobStatus, error string) {
	message := map[string]string{
		"id":     id,
		"status": string(status),
		"url":    url,
		"error":  error,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("❌ Error marshaling job update:", err)
		return
	}

	h.Broadcast(id, jsonData)
}

func GetHub() *Hub {
	return &hub
}

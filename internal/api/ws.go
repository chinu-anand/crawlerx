package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chinu-anand/crawlerx/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Websocket Upgrade Failed",
		})
		return
	}

	client := ws.Client{Conn: conn}
	hub := ws.GetHub()
	hub.AddClient(&client)

	fmt.Println("✅ Client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("❌ Client disconnected")
			hub.RemoveClient(&client)
			break
		}

		var incoming struct {
			JobID string `json:"job_id"`
		}

		if err := json.Unmarshal(msg, &incoming); err == nil && incoming.JobID != "" {
			if client.JobIDs == nil {
				client.JobIDs = make(map[string]bool)
			}
			client.JobIDs[incoming.JobID] = true
			fmt.Println("✅ Client subscribed to job:", incoming.JobID)
		}
	}
}

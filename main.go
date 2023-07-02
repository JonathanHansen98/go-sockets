package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/JonathanHansen98/go-sockets/v2/channel_manager"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var cm = channel_manager.NewChannelManager()

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

type socketData string

type EventPayload struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type Event struct {
	Event   string `json:"event"`
	Payload EventPayload
}

func reader(conn *websocket.Conn) {
	id := uuid.New().String()
	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			return
		}

		var m Event
		if err := json.Unmarshal(p, &m); err != nil {
			log.Println("Error parsing JSON:", err)
			return
		}

		switch m.Event {
		case "join":
			c, err := cm.GetChannel(m.Payload.Channel)

			if err != nil {
				fmt.Println(err)
				cm.AddChannel(m.Payload.Channel).AddClient(id, conn)
				break
			}

			c.AddClient(id, conn)

		case "send":
			c, err := cm.GetChannel(m.Payload.Channel)

			if err != nil {
				fmt.Println(err)
				break
			}

			c.Broadcast(messageType, m.Payload.Message)

		case "leave":
			c, err := cm.GetChannel(m.Payload.Channel)

			if err != nil {
				fmt.Println(err)
				break
			}

			if err := c.DisconnectClient(id); err != nil {
				log.Println(err)
			}

		default:
			fmt.Printf("Recieved invalid event type: %s\n", m.Event)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	log.Println("New websocket connection")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
}

func main() {
	setupRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

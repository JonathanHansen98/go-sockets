package channel_manager

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type Connections map[string]*websocket.Conn

type Channel struct {
	Name        string
	Connections Connections
}

func NewChannel(name string) *Channel {
	fmt.Printf("Creating channel: %s\n", name)

	p := Channel{Name: name, Connections: make(Connections)}

	return &p
}

func (c *Channel) DisconnectClient(name string) error {
	fmt.Printf("Disconnecting client: %s from channel: %s\n", name, c.Name)

	conn, ok := c.Connections[name]

	if !ok {
		m := fmt.Sprintf("Error disconnecting client: Does not exist.")

		return errors.New(m)
	}

	conn.Close()

	return nil
}

func (c *Channel) DisconnectClients() {
	fmt.Printf("Disconnecting all clients channel: %s\n", c.Name)

	for k := range c.Connections {
		c.DisconnectClient(k)
	}
}

func (c *Channel) AddClient(name string, conn *websocket.Conn) {
	fmt.Printf("Adding client: %s to channel: %s\n", name, c.Name)

	c.Connections[name] = conn
}

func (c *Channel) Broadcast(messageType int, message string) {
	for _, v := range c.Connections {
		v.WriteMessage(messageType, []byte(message))
	}
}

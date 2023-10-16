package websocket

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := string(p)
		_message := strings.Split(message, " ")
		action := _message[0]
		if action == "start" {
			log.Printf("Action Received: %s\n", action)
			c.Pool.Action <- Action{Action: message, Client: c}
		} else {
			log.Printf("Message Received: %+v\n", message)
			c.Pool.Broadcast <- Message{User: "system", Msg: message}
		}
	}
}

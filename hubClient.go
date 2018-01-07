package doWebStocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type hubClient struct {
	id   [16]byte
	send chan *HubMessage
	hub  *Hub
	conn *websocket.Conn
}

func (c *hubClient) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	if hubFuncs == nil {
		return
	}
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.process(message)
	}
}

func (c *hubClient) write() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				break
			}
			b, err := json.Marshal(msg)
			if err != nil {
				w.Write(b)
			}

			for len(c.send) > 0 {
				b2, err := json.Marshal(<-c.send)
				if err != nil {
					w.Write(b2)
				}
			}
		}
	}
}

func (c *hubClient) process(message []byte) {

	msg := parseMessage(message)
	msgFunc, ok := c.hub.hubFunc[msg.Name]

	if !ok {
		return
	}
	msgFunc(msg.Content)

	c.send <- msg
}

func parseMessage(message []byte) *HubMessage {
	return &HubMessage{}
}

//HubMessage : WebStocketHub return
type HubMessage struct {
	Name    string
	Content string
}

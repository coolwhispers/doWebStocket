package wsfunc

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type hubClient struct {
	id   [16]byte
	Send chan *HubMessage
	hub  *Hub
	conn *websocket.Conn
}

func (c *hubClient) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

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
		case msg, ok := <-c.Send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}

			c.sendMsg(msg)

			for len(c.Send) > 0 {
				c.sendMsg(<-c.Send)
			}
		}
	}
}

func (c *hubClient) sendMsg(msg *HubMessage) {

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("object convert fail:", msg)
		return
	}

	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	w.Write(b)
}

func (c *hubClient) process(message []byte) {

	msg := parseMessage(message)
	msgFunc, ok := c.hub.hubFunc[msg.Name]

	if !ok {
		log.Println(msg.Name, "not found")
		return
	}

	r, err := msgFunc(msg.JSON)

	if err == nil {
		log.Println(err)
		return
	}

	c.Send <- &r
}

func parseMessage(message []byte) *HubMessage {
	return &HubMessage{}
}

//HubMessage : WebStocketHub return
type HubMessage struct {
	Name string
	JSON string
}

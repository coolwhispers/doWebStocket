package wshub

import (
	"log"

	"github.com/gorilla/websocket"
)

//ClientID : Client ID
type ClientID [16]byte

type client struct {
	ID   ClientID
	Send chan []byte
	hub  *Hub
	conn *websocket.Conn
}

func (c *client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		err = c.hub.HubFunc(c.ID, msg)

		if err == nil {
			log.Println(err)
		}
	}
}

func (c *client) write() {
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

func (c *client) sendMsg(msg []byte) {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	w.Write(msg)
}

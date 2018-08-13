package wshub2

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//ClientID is a type use for Client ID type
type ClientID [16]byte

type client struct {
	ID   ClientID
	send chan []byte
	conn *websocket.Conn
	hub  *Hub
}

func (c *client) SendMsg(msg interface{}) {
	b, err := json.Marshal(msg)

	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	c.send <- b
}

func (c *client) write() {
	defer func() {
		c.hub.unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.send:

			c.sendMsg(ok, msg)

			if len(c.send) > 0 {
				msg, ok := <-c.send
				c.sendMsg(ok, msg)
			}
		}
	}
}

func (c *client) sendMsg(ok bool, msg []byte) {
	if !ok {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *client) read() {
	defer func() {
		c.hub.unregister <- c
	}()

	for {
		mt, msg, err := c.conn.ReadMessage()

		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		switch mt {
		default:
			c.hub.MsgFunc(c.ID, msg)
		}
	}
}

//Hub for websocket
type Hub struct {
	Clients    map[ClientID]*client
	register   chan *client
	unregister chan *client
	MsgFunc    func(clientID ClientID, msg []byte)
}

var upgrader websocket.Upgrader

//Handler is for register websocket client
func (h *Hub) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	c := &client{
		conn: conn,
		send: make(chan []byte),
		hub:  h,
	}

	go c.read()
	go c.write()

	h.register <- c
}

func newID(r *http.Request) ClientID {
	var id [16]byte
	id = uuid.New()
	h := sha256.New()
	h.Sum(id[:])

	return id
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			client.ID = func() ClientID {
				var id [16]byte
				id = uuid.New()
				h := sha256.New()
				h.Sum(id[:])

				return id
			}()

			h.Clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.Clients[client.ID]; ok {
				client.conn.Close()
				delete(h.Clients, client.ID)
				close(client.send)
			}
		}
	}
}

//New a Websocket Hub
func New() *Hub {
	h := &Hub{
		Clients:    make(map[ClientID]*client),
		register:   make(chan *client),
		unregister: make(chan *client),
	}

	go h.run()

	h.MsgFunc = func(clientID ClientID, msg []byte) {
		str := string(msg[:])
		log.Printf("client: %x, msg: %v", clientID, str)
		h.Clients[clientID].SendMsg(str)
	}

	return h
}

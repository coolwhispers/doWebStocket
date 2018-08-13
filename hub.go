package wshub

import (
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
	run  bool
}

func (c *client) SendMsg(msg interface{}) {
	b, err := json.Marshal(msg)

	if err != nil {
		if c.hub.OnError != nil {
			c.hub.OnError(c.ID, err)
		}
		return
	}

	c.send <- b
}

func (c *client) write() {
	defer func() {
		c.hub.unregister <- c
	}()

	for c.run {
		msg, ok := <-c.send

		c.sendMsg(ok, msg)

		if len(c.send) > 0 {
			msg, ok := <-c.send
			c.sendMsg(ok, msg)
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

	for c.run {
		mt, msg, err := c.conn.ReadMessage()

		if err != nil {
			if c.hub.OnError != nil {
				c.hub.OnError(c.ID, err)
			}
			break
		}

		switch mt {
		default:
			c.hub.OnMessage(c.ID, msg)
		}
	}
}

//Hub for websocket
type Hub struct {
	Clients    map[ClientID]*client
	register   chan *client
	unregister chan *client
	OnMessage  func(clientID ClientID, msg []byte)
	OnOpen     func(clientID ClientID)
	OnClose    func(clientID ClientID)
	OnError    func(clientID ClientID, err error)
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
		run:  true,
	}

	go c.read()
	go c.write()

	h.register <- c
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:

			client.ID = func(h *Hub) ClientID {
				var id [16]byte
				id = uuid.New()

				for _, ok := h.Clients[id]; ok; {
					id = uuid.New()
				}

				return id
			}(h)

			h.Clients[client.ID] = client

			h.OnOpen(client.ID)
		case client := <-h.unregister:
			client.run = false
			if _, ok := h.Clients[client.ID]; ok {
				h.OnClose(client.ID)
				client.conn.Close()
				close(client.send)
				delete(h.Clients, client.ID)
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

	h.OnMessage = func(cid ClientID, msg []byte) {
		str := string(msg[:])
		log.Printf("client %x send msg: %v", cid, str)
		h.Clients[cid].SendMsg(str)
	}

	return h
}

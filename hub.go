package wshub

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//Hub :
type Hub struct {
	Clients    map[ClientID]*client
	register   chan *client
	unregister chan *client
	HubFunc    func(clientID ClientID, msg []byte) (err error)
	upgrader   websocket.Upgrader
}

//Handler : HTTP Handler func
func (h *Hub) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	c := &client{ID: h.createID(r), hub: h, conn: conn, Send: make(chan []byte, 256)}

	go c.read()
	go c.write()

	h.register <- c
}

func newUUIDtoCID(r *http.Request) ClientID {
	var id [16]byte
	id = uuid.New()
	return id
}
func (h *Hub) createID(r *http.Request) [16]byte {
	clientID := newUUIDtoCID(r)
	_, ok := h.Clients[clientID]
	for ok {
		clientID = newUUIDtoCID(r)
		_, ok = h.Clients[clientID]
	}
	return clientID
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
		}
	}
}

//JsHandler :
func (h *Hub) JsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

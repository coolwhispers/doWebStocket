package wsfunc

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

//Hub :
type Hub struct {
	path       string
	Clients    map[[16]byte]*hubClient
	register   chan *hubClient
	unregister chan *hubClient
	hubFunc    map[string]func(JSON string) (msg HubMessage, err error)
	upgrader   websocket.Upgrader
}

//New : Create Hub Route
func New(hubPath string) *Hub {
	h := &Hub{
		path:       hubPath,
		Clients:    make(map[[16]byte]*hubClient),
		register:   make(chan *hubClient),
		unregister: make(chan *hubClient),
		hubFunc:    make(map[string]func(message string) (msg HubMessage, err error)),
		upgrader:   websocket.Upgrader{},
	}

	go h.run()

	return h
}

//AddFunc : add webstocket func
func (h *Hub) AddFunc(funcName string, msgFunc func(JSON string) (msg HubMessage, err error)) {
	funcName = strings.ToLower(funcName)
	if _, ok := h.hubFunc[funcName]; ok {
		panic("hubFunc is exist")
	}
	h.hubFunc[funcName] = msgFunc
}

//Handler : HTTP Handler func
func (h *Hub) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	c := &hubClient{id: h.createID(), hub: h, conn: conn, Send: make(chan *HubMessage, 256)}

	go c.read()
	go c.write()

	h.register <- c
}

func (h *Hub) createID() [16]byte {
	clientID := uuid.New()
	_, ok := h.Clients[clientID]
	for ok {
		clientID = uuid.New()
		_, ok = h.Clients[clientID]
	}
	return clientID
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.Clients[client.id]; ok {
				delete(h.Clients, client.id)
				close(client.Send)
			}
		}
	}
}

//JsHandler :
func (h *Hub) JsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

package wshub

import (
	"log"

	"github.com/gorilla/websocket"
)

//New : Create Hub Route
func New() *Hub {
	h := &Hub{
		Clients:    make(map[ClientID]*client),
		register:   make(chan *client),
		unregister: make(chan *client),
		HubFunc: func(cid ClientID, msg []byte) {
			log.Printf("cid: %v, msg: %v", cid, msg)
		},
		upgrader: websocket.Upgrader{},
	}

	go h.run()

	return h
}

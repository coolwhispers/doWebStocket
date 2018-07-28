package wshub

import "github.com/gorilla/websocket"

//New : Create Hub Route
func New() *Hub {
	h := &Hub{
		Clients:    make(map[ClientID]*client),
		register:   make(chan *client),
		unregister: make(chan *client),
		HubFunc:    nil,
		upgrader:   websocket.Upgrader{},
	}

	go h.run()

	return h
}

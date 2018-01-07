package wshub

import "github.com/gorilla/websocket"

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

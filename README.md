# WebSocket Hub

```go
func main() {
    hub := wshub.New()
    hub.HubFunc = func(cid wshub.ClientID, msg []byte) {

        // do something...ex: parseByteToJsonObject(msg)

    }
    // register router
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/ws", hub.Handler)

    // start server listening
    http.ListenAndServe("0.0.0.0:80", router)
}
```

```go
hub.OnMessage = func(cid wshub.ClientID, msg []byte) {

    // parse client msg like...
    obj := parseByteToJsonObject(msg)

}
```

```go
hub.OnMessage = func(cid wshub.ClientID, msg []byte) {

    // send a message to client
    hub.Clients[cid].Send("Hello!")

}
```

```go
hub.OnOpen = func(cid wshub.ClientID) {

    log.Printf("client %x was connected.", cid)

}
```

```go
hub.OnClose = func(cid wshub.ClientID) {

    log.Printf("client %x was disconnected.", cid)

}
```

```go
hub.OnError = func(cid wshub.ClientID, err error) {

    log.Printf("client %x has error, %v", cid, err)

}
```
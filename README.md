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
hub.HubFunc = func(cid wshub.ClientID, msg []byte) {

    // parse client msg like...
    obj := parseByteToJsonObject(msg)

}
```

```go
hub.HubFunc = func(cid wshub.ClientID, msg []byte) {

    // send a message to client
    hub.Clients[cid].Send("Hello!")

}
```
# wshub-go

```go
func main() {
    hub := wshub.New()
    hub.HubFunc = func(cid wshub.ClientID, msg []byte) {

        // do something...ex: parseByte(msg)

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
    parseByte(msg)

}
```

```go
hub.HubFunc = func(cid wshub.ClientID, msg []byte) {

    // send a message to client
    hub.Clients[cid].Send("Hello!")

}
```
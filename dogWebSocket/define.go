package dogWebSocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var CodeForcedLogout = "1000"

var Receive = func(uint64, string) {}

var GetUid = func(http.ResponseWriter, *http.Request) uint64 {
	return 0
}

var Register = func(uint64) bool {
	return true
}

var UnRegister = func(uint64) bool {
	return true
}

type client struct {
	uid  uint64
	conn *websocket.Conn
	send chan []byte
}

var clients = make(map[uint64]*client)
var registerChan = make(chan *client)
var unregisterChan = make(chan *client)

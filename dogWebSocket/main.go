package dogWebSocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/goxt/dog2/env"
	"net/http"
	"strconv"
)

// 启动ws主进程
func Run() {
	go func() {

		port := strconv.Itoa(env.Config.App.WsPort)

		fmt.Println("Now listening websocket on localhost:" + port)

		go start()

		http.HandleFunc("/ws", wsUpgrade)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic(err)
		}
	}()
}

// HTTP升级成WebSocket
func wsUpgrade(res http.ResponseWriter, req *http.Request) {

	uid := GetUid(res, req)
	if uid <= 0 {
		return
	}

	var up = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := up.Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	client := &client{
		uid:  uid,
		conn: conn,
		send: make(chan []byte),
	}

	registerChan <- client

	go client.read()
	go client.write()
}

// 连接和断开
func start() {
	for {
		select {
		case conn := <-registerChan:
			resister(conn)

		case conn := <-unregisterChan:
			unResister(conn)
		}
	}
}

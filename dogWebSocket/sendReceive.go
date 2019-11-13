package dogWebSocket

import (
	"encoding/json"
	"github.com/goxt/dog2/util"
	"websocket"
)

func Send(id uint64, d interface{}) bool {
	var client = clients[id]
	if client == nil {
		return false
	}

	jm, _ := json.Marshal(&d)
	client.send <- jm
	return true
}

func (c *client) read() {

	defer func() {
		unregisterChan <- c
	}()

	for {

		_, b, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		go func() {
			defer func() {
				if e := recover(); e != nil {
					switch v := e.(type) {
					case util.BizException:
						util.LogError("[warn] ws.receive - ", v)
					case util.SysException:
						util.LogException("[error] ws.receive - " + v.Msg)
					case string:
						util.LogException("[error] ws.receive - " + v)
					default:
						util.LogException("[error] ws.receive - " + (v.(error).Error()))
					}
				}
			}()
			Receive(c.uid, string(b))
		}()
	}
}

func (c *client) write() {

	defer func() {
		unregisterChan <- c
	}()

	for {
		select {
		case b, ok := <-c.send:
			if !ok {
				return
			}
			_ = c.conn.WriteMessage(websocket.TextMessage, b)
		}
	}
}

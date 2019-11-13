package dogWebSocket

func resister(conn *client) {

	defer func() {
		recover()
	}()

	if clients[conn.uid] != nil {
		Send(conn.uid, `{"code":` + CodeForcedLogout + `,"msg":"您的账号已在其他地方登录"}`)
		unResister(clients[conn.uid])
	}

	ok := Register(conn.uid)
	if !ok {
		return
	}

	clients[conn.uid] = conn
}

func unResister(conn *client) {

	defer func() {
		recover()
	}()

	ok := UnRegister(conn.uid)
	if !ok {
		return
	}

	if conn == nil {
		return
	}

	if clients[conn.uid] == nil {
		return
	}

	if conn.send != nil {
		close(conn.send)
	}

	if conn.conn != nil {
		_ = conn.conn.Close()
	}

	delete(clients, conn.uid)
	conn = nil
}

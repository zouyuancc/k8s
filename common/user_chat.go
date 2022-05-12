package common

import (
	"net"
)

type CUser struct {
	Name string
	Addr string

	C    chan string
	conn net.Conn

	server *Server
}

func NewCUser(conn net.Conn, server *Server) *CUser {
	user := &CUser{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		server: server,
		conn:   conn,
	}
	go user.listenmsg()
	return user
}

func (user *CUser) listenmsg() {
	for {
		msg := <-user.C
		go user.conn.Write([]byte(msg))
	}
}

func (user *CUser) Online() {
	user.server.CMaplock.Lock()
	user.server.COnlineMap[user.Name] = user
	user.server.CMaplock.Unlock()

	user.SendMsg([]byte("您已上线\n"))
}

func (user *CUser) Offline() {
	user.server.CMaplock.Lock()
	delete(user.server.COnlineMap, user.Name)
	user.server.CMaplock.Unlock()
}

func (user *CUser) SendMsg(msg []byte) {
	user.conn.Write(msg)
}

func (user *CUser) Domessage(msg string) {
	if msg == "OnlineUser" {
		//查询用户有哪些

		user.server.Maplock.Lock()
		for _, u := range user.server.COnlineMap {
			OnlineMsg := "[" + u.Addr + "]" + u.Name + ":" + "在线..." + "\n"
			user.SendMsg([]byte(OnlineMsg))
		}
		user.server.Maplock.Unlock()
	} else {
		//
	}
}

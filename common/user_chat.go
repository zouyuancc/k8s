package common

import (
	"net"
	"strings"
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
	if msg == "who" {
		//查询用户有哪些

		user.server.Maplock.Lock()
		for _, u := range user.server.COnlineMap {
			OnlineMsg := "[" + u.Addr + "]" + u.Name + ":" + "在线..." + "\n"
			user.SendMsg([]byte(OnlineMsg))
		}
		user.server.Maplock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMsg([]byte("当前用户名已经被使用\n"))
		} else {
			user.server.Maplock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.COnlineMap[newName] = user
			user.server.Maplock.Unlock()

			user.Name = newName
			user.SendMsg([]byte("您已经更新用户名" + user.Name + "\n"))
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式： to|张三|消息内容

		//1 获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg([]byte("消息格式不正确，请使用\"to|张三|你好\"格式."))
			return
		}

		//2 根据用户名获取对方User对象
		remoteUser, ok := user.server.COnlineMap[remoteName]
		if !ok {
			user.SendMsg([]byte("该用户名不存在"))
			return
		}

		//获取消息内容并发送
		msgsend := strings.Split(msg, "|")[2]
		if msgsend == "" {
			user.SendMsg([]byte("无消息内容，请重发\n"))
			return
		}
		remoteUser.SendMsg([]byte("message from " + user.Name + ":" + msgsend + "\n"))
	} else {
		//
	}
}
